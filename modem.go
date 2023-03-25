package main

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	ModemManagerInterface  = "org.freedesktop.ModemManager1"
	ModemManagerObjectPath = "/org/freedesktop/ModemManager1"

	ModemInterface           = ModemManagerInterface + ".Modem"
	ModemLocationInterface   = ModemInterface + ".Location"
	ModemLocationSetup       = ModemLocationInterface + ".Setup"
	ModemLocationGetLocation = ModemLocationInterface + ".GetLocation"

	MmModemLocationSourceGpsRaw uint32 = 2

	GetManagedObjects = "org.freedesktop.DBus.ObjectManager.GetManagedObjects"
)

// Simplified ModemManager interface (based on github.com/maltegrosse/go-modemmanager).
// Not using that library due to a bug in parsing location.
type Modem interface {
	SetupLocation() error
	GetLocation() (Location, error)
}

type Location struct {
	Latitude  float64
	Longitude float64
}

func NewModem() (Modem, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	var modemPaths []dbus.ObjectPath
	managedObjects := make(map[dbus.ObjectPath]interface{})

	modemManager := conn.Object(ModemManagerInterface, ModemManagerObjectPath)
	err = modemManager.Call(GetManagedObjects, 0).Store(&managedObjects)
	if err != nil {
		return nil, err
	}

	// Get the paths of all modems.
	for path := range managedObjects {
		modemPaths = append(modemPaths, path)
	}

	if len(modemPaths) != 1 {
		return nil, fmt.Errorf("unexpected number of modems found: %v", len(modemPaths))
	}

	modemObj := conn.Object(ModemManagerInterface, modemPaths[0])

	return &modem{modemObj}, nil
}

type modem struct {
	modemObj dbus.BusObject
}

func (m *modem) call(method string, args ...interface{}) error {
	return m.modemObj.Call(method, 0, args...).Err
}

func (m *modem) callWithReturn(ret interface{}, method string, args ...interface{}) error {
	return m.modemObj.Call(method, 0, args...).Store(ret)
}

func (m *modem) SetupLocation() error {
	return m.call(ModemLocationSetup, MmModemLocationSourceGpsRaw, false)
}

func (m *modem) GetLocation() (loc Location, err error) {
	var res map[uint32]dbus.Variant
	err = m.callWithReturn(&res, ModemLocationGetLocation)
	if err != nil {
		return
	}

	return m.createLocation(res)
}

func (m *modem) createLocation(res map[uint32]dbus.Variant) (loc Location, err error) {
	for key, element := range res {
		// Only interested in GpsRaw-type entries
		if key == MmModemLocationSourceGpsRaw {
			tmpMap, ok := element.Value().(map[string](dbus.Variant))

			if ok {
				for k, v := range tmpMap {
					switch k {
					case "latitude":
						tmpVal, ok := v.Value().(float64)
						if ok {
							loc.Latitude = tmpVal
						}
					case "longitude":
						tmpVal, ok := v.Value().(float64)
						if ok {
							loc.Longitude = tmpVal
						}
					}
				}
			}
		}
	}
	return
}

package main

import "fmt"

type Location struct {
	Latitude  float64
	Longitude float64
}

type LocationProvider interface {
	GetLocation() (Location, error)
}

type LocationProviderConstructor func() (LocationProvider, error)

var PROVIDERS_PRIORITY = []LocationProviderConstructor{
	NewModemLocationProvider,
	NewStarlinkLocationProvider,
}

// Finds the first available LocationProvider or returns an error if none are available.
func SelectLocationProvider() (LocationProvider, error) {
	for _, providerConstructor := range PROVIDERS_PRIORITY {
		provider, err := providerConstructor()
		if err == nil {
			return provider, nil
		}
	}
	return nil, fmt.Errorf("no suitable LocationProvider found")
}

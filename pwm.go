package main

import (
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

type Spi struct {
	port spi.PortCloser
	conn spi.Conn
}

func InitSpi() (*Spi, error) {
	host.Init()
	port, err := spireg.Open("SPI0.0")
	if err != nil {
		return nil, err
	}

	conn, err := port.Connect(physic.MegaHertz, spi.Mode1, 8)
	if err != nil {
		return nil, err
	}

	return &Spi{
		port,
		conn,
	}, nil
}

func (s *Spi) WritePwm(left byte, right byte) error {
	write := make([]byte, 2)
	read := make([]byte, len(write))
	write[0] = left
	write[1] = right

	return s.conn.Tx(write, read)
}

func (s *Spi) Close() {
	s.port.Close()
}

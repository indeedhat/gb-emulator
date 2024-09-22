package emu

import "log"

type MBC3 struct {
	data    []byte
	variant uint8
}

func NewMBC3(data []byte, variant uint8) MBC3 {
	return MBC3{data: data, variant: variant}
}

func (m MBC3) Read(address uint16) byte {
	return m.data[address]
}

func (m MBC3) Write(address uint16, value byte) {
	log.Print("cart.write not implemented")
}

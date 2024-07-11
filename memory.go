package main

type MemoryBus struct {
	Cart *Cartridge
}

func (b *MemoryBus) Read(address uint16) uint8 {
	switch true {
	case address < 0x8000:
		// TODO: this is probably wrong
		return b.Cart.Read(address)
	default:
		panic("not implemented")
	}
}

func (b *MemoryBus) Write(address uint16, value uint8) {
	switch true {
	case address < 0x8000:
		// TODO: this is probably wrong
		b.Cart.Write(address, value)
	default:
		panic("not implemented")
	}
}

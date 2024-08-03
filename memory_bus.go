package main

import (
	"log"
)

// 0x0000 - 0x3FFF:   16 KiB ROM bank 00              From cartridge, usually a fixed bank
// 0x4000 - 0x7FFF:   16 KiB ROM Bank 01–NN           From cartridge, switchable bank via mapper (if any)
// 0x8000 - 0x9FFF:   8 KiB Video RAM (VRAM)          In CGB mode, switchable bank 0/1
// 0xA000 - 0xBFFF:   8 KiB External RAM              From cartridge, switchable bank if any
// 0xC000 - 0xCFFF:   4 KiB Work RAM (WRAM)
// 0xD000 - 0xDFFF:   4 KiB Work RAM (WRAM)           In CGB mode, switchable bank 1–7
// 0xE000 - 0xFDFF:   Echo RAM (mirror of C000–DDFF)  Nintendo says use of this area is prohibited.
// 0xFE00 - 0xFE9F:   Object attribute memory (OAM)
// 0xFEA0 - 0xFEFF:   Not Usable                      Nintendo says use of this area is prohibited.
// 0xFF00 - 0xFF7F:   I/O Registers
// 0xFF80 - 0xFFFE:   High RAM (HRAM)
// 0xFFFF - 0xFFFF:   Interrupt Enable register (IE)

type MemoryBus struct {
	hram *RamBank
	wram *RamBank

	ctx *Context
}

func NewMemoryBus(ctx *Context) {
	ctx.membus = &MemoryBus{
		hram: &RamBank{
			offset: 0xFF80,
			data:   make([]byte, 0x80),
		},
		wram: &RamBank{
			offset: 0xC000,
			data:   make([]byte, 0x2000),
		},
		ctx: ctx,
	}
}

func (b *MemoryBus) Read(address uint16) uint8 {
	switch true {
	case address < 0x8000:
		return b.ctx.cart.Read(address)
	case address < 0xA000:
		return b.ctx.ppu.vram.Read(address)
	case address < 0xC000:
		// cart ram
		return b.ctx.cart.Read(address)
	case address < 0xE000:
		// working ram
		return b.wram.Read(address)
	case address < 0xFE00:
		// Echo ram is unusable
		return 0
	case address < 0xFEA0:
		if b.ctx.dma.Active {
			return 0xFF
		}

		value := b.ctx.ppu.oam.Read(address)
		if address == 0xFE40 {
			log.Fatalf("r %d,%d", address, value)
		}
		return value
	case address < 0xFF00:
		// reserved and unusable
		return 0
	case address < 0xFF80:
		// IO registers
		return b.ctx.io.Read(address)
	case address < 0xFFFF:
		// high ram/zero page
		return b.hram.Read(address)
	case address == 0xFFFF:
		// CPU ENABLE REIGSTER
		return b.ctx.cpu.interruptRegister
	default:
		log.Printf("unsupported mem.read 0x%X", address)
	}
	return 0
}

func (b *MemoryBus) Read16(address uint16) uint16 {
	return uint16(b.Read(address)) | uint16(b.Read(address+1))<<8
}

func (b *MemoryBus) Write(address uint16, value uint8) {
	switch true {
	case address < 0x8000:
		b.ctx.cart.Write(address, value)
	case address < 0xA000:
		b.ctx.ppu.vram.Write(address, value)
	case address < 0xC000:
		// cart ram
		b.ctx.cart.Write(address, value)
	case address < 0xE000:
		// working ram
		b.wram.Write(address, value)
	case address < 0xFE00:
		// Echo ram is unusable
	case address < 0xFEA0:
		if b.ctx.dma.Active {
			return
		}
		if address == 0xFE40 {
			log.Printf("w %d,%d", address, value)
		}
		b.ctx.ppu.oam.Write(address, value)
	case address < 0xFF00:
		// reserved and unusable
	case address < 0xFF80:
		// IO registers
		b.ctx.io.Write(address, value)
	case address < 0xFFFF:
		// high ram/zero page
		b.hram.Write(address, value)
	case address == 0xFFFF:
		// CPU ENABLE REIGSTER
		b.ctx.cpu.interruptRegister = value
	default:
		log.Printf("unsupported mem.write 0x%X", address)
	}
}

func (b *MemoryBus) Write16(address uint16, value uint16) {
	b.Write(address, uint8(value&0xFF))
	b.Write(address+1, uint8(value>>8))
}

type RamBank struct {
	offset uint16
	data   []byte
}

func (r *RamBank) Read(address uint16) uint8 {
	return r.data[address-r.offset]

}

func (r *RamBank) Write(address uint16, value uint8) {
	r.data[address-r.offset] = value
}

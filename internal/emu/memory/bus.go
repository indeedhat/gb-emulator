package memory

import (
	"bytes"
	"log"

	"github.com/indeedhat/gb-emulator/internal/emu/context"
	. "github.com/indeedhat/gb-emulator/internal/emu/types"
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

	ctx *context.Context
}

func NewBus(ctx *context.Context) {
	ctx.Bus = &MemoryBus{
		hram: NewRamBank(0xFF80, 0x80),
		wram: NewRamBank(0xC000, 0x2000),
		ctx:  ctx,
	}
}

func (b *MemoryBus) LoadState(data []byte) {
	r := bytes.NewReader(data)

	h := b.hram.Bytes()
	r.Read(h)
	b.hram.Fill(h)

	w := b.wram.Bytes()
	r.Read(w)
	b.wram.Fill(w)

}

func (b *MemoryBus) SaveState() []byte {
	var buf bytes.Buffer

	buf.Write(b.hram.Bytes())
	buf.Write(b.wram.Bytes())

	return buf.Bytes()
}

func (b *MemoryBus) Read(address uint16) uint8 {
	switch true {
	case address < 0x8000:
		return b.ctx.Cart.Read(address)
	case address < 0xA000:
		return b.ctx.Ppu.Read(address)
	case address < 0xC000:
		// cart ram
		return b.ctx.Cart.Read(address)
	case address < 0xE000:
		// working ram
		return b.wram.Read(address)
	case address < 0xFE00:
		// Echo ram is unusable
		return 0
	case address < 0xFEA0:
		if b.ctx.Dma.Active() {
			return 0xFF
		}

		value := b.ctx.Ppu.Read(address)
		// if address == 0xFE40 {
		// 	log.Fatalf("r %d,%d", address, value)
		// }
		return value
	case address < 0xFF00:
		// reserved and unusable
		return 0
	case address < 0xFF80:
		// IO registers
		return b.ctx.Io.Read(address)
	case address < 0xFFFF:
		// high ram/zero page
		return b.hram.Read(address)
	case address == 0xFFFF:
		// CPU ENABLE REIGSTER
		return b.ctx.Cpu.InterruptRegister()
	default:
		// log.Printf("unsupported mem.read 0x%X", address)
	}
	return 0
}

func (b *MemoryBus) Read16(address uint16) uint16 {
	return uint16(b.Read(address)) | uint16(b.Read(address+1))<<8
}

func (b *MemoryBus) Write(address uint16, value uint8) {
	switch true {
	case address < 0x8000:
		b.ctx.Cart.Write(address, value)
	case address < 0xA000:
		b.ctx.Ppu.Write(address, value)
	case address < 0xC000:
		// cart ram
		b.ctx.Cart.Write(address, value)
	case address < 0xE000:
		// working ram
		b.wram.Write(address, value)
	case address < 0xFE00:
		// Echo ram is unusable
	case address < 0xFEA0:
		if b.ctx.Dma.Active() {
			return
		}
		if address == 0xFE40 {
			log.Printf("w %d,%d", address, value)
		}
		b.ctx.Ppu.Write(address, value)
	case address < 0xFF00:
		// reserved and unusable
	case address < 0xFF80:
		// IO registers
		b.ctx.Io.Write(address, value)
	case address < 0xFFFF:
		// high ram/zero page
		b.hram.Write(address, value)
	case address == 0xFFFF:
		// CPU ENABLE REIGSTER
		b.ctx.Cpu.SetInterruptRegister(value)
	default:
		// log.Printf("unsupported mem.write 0x%X", address)
	}
}

func (b *MemoryBus) Write16(address uint16, value uint16) {
	b.Write(address, uint8(value&0xFF))
	b.Write(address+1, uint8(value>>8))
}

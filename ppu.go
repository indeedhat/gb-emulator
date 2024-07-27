package main

type Ppu struct {
	oam  OamRam
	vram *RamBank
}

func NewPpu(ctx *Context) {
	ctx.ppu = &Ppu{
		make(OamRam, 40),
		&RamBank{
			offset: 0x8000,
			data:   make([]byte, 0x2000),
		},
	}
}

type OamEntry struct {
	x       uint8
	y       uint8
	tileIdx uint8
	// 7   Priority:    0 = No, 1 = BG and Window colors 1–3 are drawn over this OBJ
	// 6   Y flip:      0 = Normal, 1 = Entire OBJ is vertically mirrored
	// 5   X flip:      0 = Normal, 1 = Entire OBJ is horizontally mirrored
	// 4   DMG palette: 0 = OBP0, 1 = OBP1
	// 3   Bank:        0 = Fetch tile from VRAM bank 0, 1 = Fetch tile from VRAM bank 1
	// 0-2 CGB palette: Which of OBP0–7 to use
	flags uint8
}

func (e *OamEntry) Read(address uint8) uint8 {
	switch address {
	case 0:
		return e.x
	case 1:
		return e.y
	case 2:
		return e.tileIdx
	case 3:
		return e.flags
	default:
		panic("invalid address for oam entry")
	}
}

func (e *OamEntry) Write(address uint8, value uint8) {
	switch address {
	case 0:
		e.x = value
	case 1:
		e.y = value
	case 2:
		e.tileIdx = value
	case 3:
		e.flags = value
	default:
		panic("invalid address for oam entry")
	}
}

type OamRam []*OamEntry

func (o *OamRam) Read(address uint16) uint8 {
	address -= 0xFE00
	i := address % 4
	return (*o)[(address-i)/4].Read(uint8(i))
}

func (o *OamRam) Write(address uint16, value uint8) {
	address -= 0xFE00
	i := address % 4
	(*o)[(address-i)/4].Write(uint8(i), value)
}

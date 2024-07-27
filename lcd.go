package main

type Lcdc uint8

const (
	LcdcBgwEnable = 1 < iota
	LcdcObjecteEnable
	LcdcObjecteSize
	LcdcBgTileArea
	LcdcBgwTileArea
	LcdcWindowEnable
	LcdcWindowTileArea
	LcdcLcdPpuEnable
)

type Lcd struct {
	// Bit 7: LCD & PPU enable:           0 = Off; 1 = On
	// Bit 6: Window tile map area:       0 = 9800–9BFF; 1 = 9C00–9FFF
	// Bit 5: Window enable:              0 = Off; 1 = On
	// Bit 4: BG & Window tile data area: 0 = 8800–97FF; 1 = 8000–8FFF
	// Bit 3: BG tile map area:           0 = 9800–9BFF; 1 = 9C00–9FFF
	// Bit 2: OBJ size:                   0 = 8×8; 1 = 8×16
	// Bit 1: OBJ enable:                 0 = Off; 1 = On
	// Bit 0: BG & Window enable / priority [Different meaning in CGB Mode]: 0 = Off; 1 = On
	control uint8
	ly      uint8
	lyc     uint8
	// bit 6: LYC int select (Read/Write):    If set, selects the LYC == LY condition for the STAT interrupt.
	// bit 5: Mode 2 int select (Read/Write): If set, selects the Mode 2 condition for the STAT interrupt.
	// bit 4: Mode 1 int select (Read/Write): If set, selects the Mode 1 condition for the STAT interrupt.
	// bit 3: Mode 0 int select (Read/Write): If set, selects the Mode 0 condition for the STAT interrupt.
	// bit 2: LYC == LY (Read-only):          Set when LY contains the same value as LYC; it is constantly updated.
	// bit 1-0: PPU mode (Read-only):         Indicates the PPU’s current status.
	status  uint8
	scrollY uint8
	scrollX uint8
	windowY uint8
	windoX  uint8
	// Bit 7-6: ID 3
	// Bit 5-4: ID 2
	// Bit 3-2: ID 1
	// Bit 1-0: ID 0
	// Color 0: white
	// Color 1: light grey
	// Color 2: dark grex
	// Color 3: black
	backgroundPallet uint8
	objectPallet0    uint8
	objectPallet1    uint8

	ctx *Context
}

func (l *Lcd) GetControl(code Lcdc) bool {
	return Lcdc(l.control)&code == code
}

func (l *Lcd) SetControl(code Lcdc, set bool) {
	if set {
		l.control |= uint8(code)
	} else {
		l.control &= ^uint8(code)
	}
}

func NewLcd(ctx *Context) {
	ctx.lcd = &Lcd{
		ctx: ctx,
	}
}

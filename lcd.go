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

type LcdStatus uint8

const (
	LcdStatusHblank = (1 << 3)
	LcdStatusVblank = (1 << 4)
	LcdStatusOam    = (1 << 5)
	LcdStatusLyc    = (1 << 6)
)

type LcdMode uint8

const (
	LcdModeHblank LcdMode = iota
	LcdModeVblank
	LcdModeOam
	LcdModeDrawLine
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
	dma     uint8
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
	windowX uint8
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

func NewLcd(ctx *Context) {
	ctx.lcd = &Lcd{
		ctx:              ctx,
		control:          0x91,
		status:           uint8(LcdModeOam),
		dma:              0xFF,
		backgroundPallet: 0xFC,
		objectPallet0:    0xFF,
		objectPallet1:    0xFF,
	}
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

func (l *Lcd) GetMode() LcdMode {
	return LcdMode(l.status & 0x3)
}

func (l *Lcd) SetMode(mode LcdMode) {
	l.status &= ^uint8(0x3)
	l.status |= uint8(mode)
}

func (l *Lcd) IncrementLine() {
	l.ly++

	if l.ly == l.lyc {
		l.status |= 0b10
		if l.status&LcdStatusLyc == LcdStatusLyc {
			l.ctx.cpu.requestInterrupt(InterruptLcdStat)
		}
	} else {
		l.status &= ^uint8(0b10)
	}
}

func (l *Lcd) Read(addr uint16) uint8 {
	switch addr {
	case 0xFF40:
		return l.control
	case 0xFF41:
		return l.status
	case 0xFF42:
		return l.scrollY
	case 0xFF43:
		return l.scrollX
	case 0xFF44:
		return l.ly
	case 0xFF45:
		return l.lyc
	case 0xFF46:
		return l.dma
	case 0xFF47:
		return l.backgroundPallet
	case 0xFF48:
		return l.objectPallet0
	case 0xFF49:
		return l.objectPallet1
	case 0xFF4A:
		return l.windowY
	case 0xFF4B:
		return l.windowX
	default:
		return 0xFF
	}
}

func (l *Lcd) Write(addr uint16, value uint8) {
	switch addr {
	case 0xFF40:
		l.control = value
	case 0xFF41:
		l.status = value
	case 0xFF42:
		l.scrollY = value
	case 0xFF43:
		l.scrollX = value
	case 0xFF44:
		l.ly = value
	case 0xFF45:
		l.lyc = value
	case 0xFF46:
		l.dma = value
		l.ctx.dma.Start(value)
	case 0xFF47:
		l.backgroundPallet = value
	case 0xFF48:
		l.objectPallet0 = value
	case 0xFF49:
		l.objectPallet1 = value
	case 0xFF4A:
		l.windowY = value
	case 0xFF4B:
		l.windowX = value
	}
}

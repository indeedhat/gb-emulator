package enum

type Lcdc uint8

const (
	LcdcBgwEnable Lcdc = 1 << iota
	LcdcObjecteEnable
	LcdcObjecteDoubleHeight
	LcdcBgTileArea
	LcdcBgwTileArea
	LcdcWindowEnable
	LcdcWindowTileArea
	LcdcLcdPpuEnable
)

type LcdStatus uint8

const (
	LcdStatusHblank LcdStatus = (1 << 3)
	LcdStatusVblank LcdStatus = (1 << 4)
	LcdStatusOam    LcdStatus = (1 << 5)
	LcdStatusLyc    LcdStatus = (1 << 6)
)

type LcdMode uint8

const (
	LcdModeHblank LcdMode = iota
	LcdModeVblank
	LcdModeOam
	LcdModeDraw
)

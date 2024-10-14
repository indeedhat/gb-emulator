package enum

const (
	InterruptVBlank uint8 = 1 << iota
	InterruptLcdStat
	InterruptTimer
	InterruptSerial
	InterruptJoyPad
)

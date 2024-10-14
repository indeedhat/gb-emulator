package config

import "time"

const (
	PpuLinesPerFrame = 154
	PpuTicksPerLine  = 456
	PpuYRes          = 144
	PpuXRes          = 160
	TargetFrameTime  = time.Second / 60
)

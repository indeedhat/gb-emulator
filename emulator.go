package main

type Emulator struct {
	Running bool
	Paused  bool
	Ticks   uint64
}

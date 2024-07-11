package main

type Cpu struct {
	Registers struct {
		AF [2]byte
		BC [2]byte
		DE [2]byte
		HL [2]byte
		SP [2]byte
		PC [2]byte
	}
}

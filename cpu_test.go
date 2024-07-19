package main

//
// import (
// 	"log"
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// )
//
// var incTestTable = []struct {
// 	name        string
// 	inMem       uint8
// 	inReg       cpuRegisters
// 	instruction CpuInstriction
// 	expectedReg cpuRegisters
// 	expectedMem uint8
// }{
// 	{
// 		name:        "8bit register",
// 		inMem:       0x00,
// 		inReg:       cpuRegisters{},
// 		instruction: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeA, 0, 0, 0},
// 		expectedReg: cpuRegisters{A: 0x01},
// 		expectedMem: 0x00,
// 	},
// 	{
// 		name:        "8bit register roll",
// 		inMem:       0x00,
// 		inReg:       cpuRegisters{A: 0xFF},
// 		instruction: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeA, 0, 0, 0},
// 		expectedReg: cpuRegisters{F: 0b10100000, A: 0x00},
// 		expectedMem: 0x00,
// 	},
// 	{
// 		name:        "16bit register",
// 		inMem:       0x00,
// 		inReg:       cpuRegisters{},
// 		instruction: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeBC, 0, 0, 0},
// 		expectedReg: cpuRegisters{B: 0x01},
// 		expectedMem: 0x00,
// 	},
// 	{
// 		name:        "16bit register roll",
// 		inMem:       0x00,
// 		inReg:       cpuRegisters{B: 0xFF, C: 0xFF},
// 		instruction: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeBC, 0, 0, 0},
// 		expectedReg: cpuRegisters{B: 0x00, C: 0x00},
// 		expectedMem: 0x00,
// 	},
// 	{
// 		name:        "memory",
// 		inMem:       0x00,
// 		inReg:       cpuRegisters{H: 0x00, L: 0xC0},
// 		instruction: CpuInstriction{InstructionTypeINC, AddressModeMR, RegisterTypeHL, 0, 0, 0},
// 		expectedReg: cpuRegisters{H: 0x00, L: 0xC0},
// 		expectedMem: 0x01,
// 	},
// 	{
// 		name:        "memory roll",
// 		inMem:       0xFF,
// 		inReg:       cpuRegisters{H: 0x00, L: 0xC0},
// 		instruction: CpuInstriction{InstructionTypeINC, AddressModeMR, RegisterTypeHL, 0, 0, 0},
// 		expectedReg: cpuRegisters{F: 0b10100000, H: 0x00, L: 0xC0},
// 		expectedMem: 0x00,
// 	},
// }
//
// func TestExecInc(t *testing.T) {
// 	for _, test := range incTestTable {
// 		t.Run(test.name, func(t *testing.T) {
// 			c := NewCpu(NewMemoryBus(&Cartridge{}))
// 			c.registers = &test.inReg
// 			c.membus.Write(0xC000, test.inMem)
//
// 			log.Print(c.registers.A)
// 			data, addr := c.fetchData(test.instruction)
// 			c.execINC(test.instruction, data, addr)
//
// 			assert.Equal(t, test.expectedReg.A, c.registers.A)
// 			assert.Equal(t, test.expectedReg.B, c.registers.B)
// 			assert.Equal(t, test.expectedReg.C, c.registers.C)
// 			assert.Equal(t, test.expectedReg.D, c.registers.D)
// 			assert.Equal(t, test.expectedReg.E, c.registers.E)
// 			assert.Equal(t, test.expectedReg.F, c.registers.F, "0b%08b", c.registers.F)
// 			assert.Equal(t, test.expectedReg.H, c.registers.H)
// 			assert.Equal(t, test.expectedReg.L, c.registers.L)
// 			assert.Equal(t, test.expectedMem, c.membus.wram.Read(0xC000))
// 		})
// 	}
// }
//
// var decTestTable = []struct {
// 	name        string
// 	inMem       uint8
// 	inReg       cpuRegisters
// 	instruction CpuInstriction
// 	expectedReg cpuRegisters
// 	expectedMem uint8
// }{
// 	{
// 		name:        "8bit register",
// 		inMem:       0x00,
// 		inReg:       cpuRegisters{A: 0x01},
// 		instruction: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeA, 0, 0, 0},
// 		expectedReg: cpuRegisters{A: 0x00, F: 0b11000000},
// 		expectedMem: 0x00,
// 	},
// 	{
// 		name:        "8bit register roll",
// 		inMem:       0x00,
// 		inReg:       cpuRegisters{},
// 		instruction: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeA, 0, 0, 0},
// 		expectedReg: cpuRegisters{F: 0b01100000, A: 0xFF},
// 		expectedMem: 0x00,
// 	},
// 	{
// 		name:        "16bit register",
// 		inMem:       0x00,
// 		inReg:       cpuRegisters{C: 0x01},
// 		instruction: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeBC, 0, 0, 0},
// 		expectedReg: cpuRegisters{B: 0xFF},
// 		expectedMem: 0x00,
// 	},
// 	{
// 		name:        "16bit register roll",
// 		inMem:       0x00,
// 		inReg:       cpuRegisters{},
// 		instruction: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeBC, 0, 0, 0},
// 		expectedReg: cpuRegisters{B: 0xFF, C: 0xFF},
// 		expectedMem: 0x00,
// 	},
// 	{
// 		name:        "memory",
// 		inMem:       0x01,
// 		inReg:       cpuRegisters{H: 0x00, L: 0xC0},
// 		instruction: CpuInstriction{InstructionTypeDEC, AddressModeMR, RegisterTypeHL, 0, 0, 0},
// 		expectedReg: cpuRegisters{H: 0x00, L: 0xC0, F: 0b11000000},
// 		expectedMem: 0x00,
// 	},
// 	{
// 		name:        "memory roll",
// 		inMem:       0x00,
// 		inReg:       cpuRegisters{H: 0x00, L: 0xC0},
// 		instruction: CpuInstriction{InstructionTypeDEC, AddressModeMR, RegisterTypeHL, 0, 0, 0},
// 		expectedReg: cpuRegisters{F: 0b01100000, H: 0x00, L: 0xC0},
// 		expectedMem: 0xFF,
// 	},
// }
//
// func TestExecDec(t *testing.T) {
// 	for _, test := range decTestTable {
// 		t.Run(test.name, func(t *testing.T) {
// 			c := NewCpu(NewMemoryBus(&Cartridge{}))
// 			c.registers = &test.inReg
// 			c.membus.Write(0xC000, test.inMem)
//
// 			log.Print(c.registers.A)
// 			data, addr := c.fetchData(test.instruction)
// 			c.execDEC(test.instruction, data, addr)
//
// 			assert.Equal(t, test.expectedReg.A, c.registers.A)
// 			assert.Equal(t, test.expectedReg.B, c.registers.B)
// 			assert.Equal(t, test.expectedReg.C, c.registers.C)
// 			assert.Equal(t, test.expectedReg.D, c.registers.D)
// 			assert.Equal(t, test.expectedReg.E, c.registers.E)
// 			assert.Equal(t, test.expectedReg.F, c.registers.F, "0b%08b", c.registers.F)
// 			assert.Equal(t, test.expectedReg.H, c.registers.H)
// 			assert.Equal(t, test.expectedReg.L, c.registers.L)
// 			assert.Equal(t, test.expectedMem, c.membus.wram.Read(0xC000))
// 		})
// 	}
// }

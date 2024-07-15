package main

// https://gbdev.io/gb-opcodes/optables/

import "fmt"

type InstructionType byte

func (i InstructionType) String() string {
	return []string{
		"None", "NOP", "LD",
		"INC", "DEC", "RLCA",
		"STOP",
		"RLA", "JR", "RRA",
		"DAA", "CPL", "SCF",
		"CCF", "HALT", "ADD",
		"SUB", "ADC", "SBC",
		"AND", "XOR", "OR",
		"CP", "POP", "JP",
		"PUSH", "RET", "CB",
		"CALL", "RETI", "LDH",
		"JPHL", "DI", "EI",
		"RST",
		"CB_RLC", "CB_RRC", "CB_RRCA",
		"CB_RL", "CB_RR", "CB_SLA",
		"CB_SRA", "CB_SWAP", "CB_SRL",
		"CB_BIT", "CB_RES", "CB_SET",
	}[i] + fmt.Sprintf("(0x%X)", uint8(i))
}

const (
	InstructionTypeNone InstructionType = iota
	InstructionTypeNOP
	InstructionTypeLD
	InstructionTypeINC
	InstructionTypeDEC
	InstructionTypeRLCA
	InstructionTypeSTOP
	InstructionTypeRLA
	InstructionTypeJR
	InstructionTypeRRA
	InstructionTypeDAA
	InstructionTypeCPL
	InstructionTypeSCF
	InstructionTypeCCF
	InstructionTypeHALT
	InstructionTypeADD
	InstructionTypeSUB
	InstructionTypeADC
	InstructionTypeSBC
	InstructionTypeAND
	InstructionTypeXOR
	InstructionTypeOR
	InstructionTypeCP
	InstructionTypePOP
	InstructionTypeJP
	InstructionTypePUSH
	InstructionTypeRET
	InstructionTypeCB
	InstructionTypeCALL
	InstructionTypeRETI
	InstructionTypeLDH
	InstructionTypeJPHL
	InstructionTypeDI
	InstructionTypeEI
	InstructionTypeRST

	// CB instructions...
	InstructionTypeRLC
	InstructionTypeRRC
	InstructionTypeRRCA
	InstructionTypeRL
	InstructionTypeRR
	InstructionTypeSLA
	InstructionTypeSRA
	InstructionTypeSWAP
	InstructionTypeSRL
	InstructionTypeBIT
	InstructionTypeRES
	InstructionTypeSET
)

type AddressMode byte

const (
	AddressModeImplied AddressMode = iota
	AddressModeR
	AddressModeR_N16
	AddressModeR_R
	AddressModeR_N8
	AddressModeR_MR
	AddressModeR_HLI
	AddressModeR_HLD
	AddressModeR_A16
	AddressModeR_A8
	AddressModeN8
	AddressModeN16
	AddressModeMR
	AddressModeMR_N8
	AddressModeMR_N16
	AddressModeMR_R
	AddressModeA8_R
	AddressModeA16_R
	AddressModeHLI_R
	AddressModeHLD_R
	AddressModeHL_SPR
)

type RegisterType byte

func (r RegisterType) Is16bit() bool {
	return r >= RegisterTypeAF
}

const (
	RegisterTypeNone RegisterType = iota
	RegisterTypeA
	RegisterTypeB
	RegisterTypeC
	RegisterTypeD
	RegisterTypeE
	RegisterTypeH
	RegisterTypeL
	RegisterTypeAF
	RegisterTypeBC
	RegisterTypeDE
	RegisterTypeHL
	RegisterTypeSP
	RegisterTypePC
)

type ConditionType byte

const (
	ConditionTypeNone ConditionType = iota
	ConditionTypeNZ
	ConditionTypeZ
	ConditionTypeNC
	ConditionTypeC
)

type CpuInstriction struct {
	Type        InstructionType
	AddressMode AddressMode
	Register1   RegisterType
	Register2   RegisterType
	Condition   ConditionType
	Param       byte
}

var CpuInstructions = []CpuInstriction{
	0x00: {InstructionTypeNOP, 0, 0, 0, 0, 0},
	0x01: {InstructionTypeLD, AddressModeR_N16, RegisterTypeBC, 0, 0, 0},
	0x02: {InstructionTypeLD, AddressModeMR_N16, RegisterTypeBC, RegisterTypeA, 0, 0},
	0x03: {InstructionTypeINC, AddressModeR, RegisterTypeBC, 0, 0, 0},
	0x04: {InstructionTypeINC, AddressModeR, RegisterTypeB, 0, 0, 0},
	0x05: {InstructionTypeDEC, AddressModeR, RegisterTypeB, 0, 0, 0},
	0x06: {InstructionTypeLD, AddressModeR_N8, RegisterTypeB, 0, 0, 0},
	0x07: {InstructionTypeRLCA, 0, 0, 0, 0, 0},
	0x08: {InstructionTypeLD, AddressModeA16_R, RegisterTypeNone, RegisterTypeSP, 0, 0},
	0x09: {InstructionTypeADD, AddressModeR_R, RegisterTypeHL, RegisterTypeBC, 0, 0},
	0x0A: {InstructionTypeLD, AddressModeR_MR, RegisterTypeA, RegisterTypeBC, 0, 0},
	0x0B: {InstructionTypeDEC, AddressModeR, RegisterTypeBC, 0, 0, 0},
	0x0C: {InstructionTypeINC, AddressModeR, RegisterTypeC, 0, 0, 0},
	0x0D: {InstructionTypeDEC, AddressModeR, RegisterTypeC, 0, 0, 0},
	0x0E: {InstructionTypeLD, AddressModeR_N8, RegisterTypeC, 0, 0, 0},
	0x0F: {InstructionTypeRRCA, 0, 0, 0, 0, 0},

	0x10: {InstructionTypeSTOP, AddressModeN8, 0, 0, 0, 0},
	0x11: {InstructionTypeLD, AddressModeR_N16, RegisterTypeDE, 0, 0, 0},
	0x12: {InstructionTypeLD, AddressModeMR_N16, RegisterTypeDE, RegisterTypeA, 0, 0},
	0x13: {InstructionTypeINC, AddressModeR, RegisterTypeDE, 0, 0, 0},
	0x14: {InstructionTypeINC, AddressModeR, RegisterTypeD, 0, 0, 0},
	0x15: {InstructionTypeDEC, AddressModeR, RegisterTypeD, 0, 0, 0},
	0x16: {InstructionTypeLD, AddressModeMR_N8, RegisterTypeD, 0, 0, 0},
	0x17: {InstructionTypeRLA, 0, 0, 0, 0, 0},
	0x18: {InstructionTypeJR, AddressModeN8, 0, 0, 0, 0},
	0x19: {InstructionTypeADD, AddressModeR_R, RegisterTypeHL, RegisterTypeDE, 0, 0},
	0x1A: {InstructionTypeLD, AddressModeR_MR, RegisterTypeA, RegisterTypeDE, 0, 0},
	0x1B: {InstructionTypeDEC, AddressModeR, RegisterTypeDE, 0, 0, 0},
	0x1C: {InstructionTypeINC, AddressModeR, RegisterTypeE, 0, 0, 0},
	0x1D: {InstructionTypeDEC, AddressModeR, RegisterTypeE, 0, 0, 0},
	0x1E: {InstructionTypeLD, AddressModeR_N8, RegisterTypeE, 0, 0, 0},
	0x1F: {InstructionTypeRRA, 0, 0, 0, 0, 0},

	0x20: {InstructionTypeJR, AddressModeN8, 0, 0, ConditionTypeNZ, 0},
	0x21: {InstructionTypeLD, AddressModeR_N16, RegisterTypeHL, 0, 0, 0},
	0x22: {InstructionTypeLD, AddressModeHLI_R, RegisterTypeHL, RegisterTypeA, 0, 0},
	0x23: {InstructionTypeINC, AddressModeR, RegisterTypeHL, 0, 0, 0},
	0x24: {InstructionTypeINC, AddressModeR, RegisterTypeH, 0, 0, 0},
	0x25: {InstructionTypeDEC, AddressModeR, RegisterTypeH, 0, 0, 0},
	0x26: {InstructionTypeLD, AddressModeR_N8, RegisterTypeH, 0, 0, 0},
	0x27: {InstructionTypeDAA, 0, 0, 0, 0, 0},
	0x28: {InstructionTypeJR, AddressModeN8, 0, 0, ConditionTypeZ, 0},
	0x29: {InstructionTypeADD, AddressModeR_R, RegisterTypeHL, RegisterTypeHL, 0, 0},
	0x2A: {InstructionTypeLD, AddressModeR_HLI, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x2B: {InstructionTypeDEC, AddressModeR, RegisterTypeHL, 0, 0, 0},
	0x2C: {InstructionTypeINC, AddressModeR, RegisterTypeL, 0, 0, 0},
	0x2D: {InstructionTypeDEC, AddressModeR, RegisterTypeL, 0, 0, 0},
	0x2E: {InstructionTypeLD, AddressModeR_N8, RegisterTypeL, 0, 0, 0},
	0x2F: {InstructionTypeCPL, 0, 0, 0, 0, 0},

	0x30: {InstructionTypeJR, AddressModeN8, 0, 0, ConditionTypeNC, 0},
	0x31: {InstructionTypeLD, AddressModeR_N16, RegisterTypeSP, 0, 0, 0},
	0x32: {InstructionTypeLD, AddressModeHLD_R, RegisterTypeHL, RegisterTypeA, 0, 0},
	0x33: {InstructionTypeINC, AddressModeR, RegisterTypeSP, 0, 0, 0},
	0x34: {InstructionTypeINC, AddressModeMR, RegisterTypeHL, 0, 0, 0},
	0x35: {InstructionTypeDEC, AddressModeMR, RegisterTypeHL, 0, 0, 0},
	0x36: {InstructionTypeLD, AddressModeMR_N8, RegisterTypeHL, 0, 0, 0},
	0x37: {InstructionTypeSCF, 0, 0, 0, 0, 0},
	0x38: {InstructionTypeJR, AddressModeN8, 0, 0, ConditionTypeC, 0},
	0x39: {InstructionTypeADD, AddressModeR_R, RegisterTypeHL, RegisterTypeSP, 0, 0},
	0x3A: {InstructionTypeLD, AddressModeR_HLD, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x3B: {InstructionTypeDEC, AddressModeR, RegisterTypeSP, 0, 0, 0},
	0x3C: {InstructionTypeINC, AddressModeR, RegisterTypeA, 0, 0, 0},
	0x3D: {InstructionTypeDEC, AddressModeR, RegisterTypeA, 0, 0, 0},
	0x3E: {InstructionTypeLD, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0x3F: {InstructionTypeCCF, 0, 0, 0, 0, 0},

	0x40: {InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeB, 0, 0},
	0x41: {InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeC, 0, 0},
	0x42: {InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeD, 0, 0},
	0x43: {InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeE, 0, 0},
	0x44: {InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeH, 0, 0},
	0x45: {InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeL, 0, 0},
	0x46: {InstructionTypeLD, AddressModeR_MR, RegisterTypeB, RegisterTypeHL, 0, 0},
	0x47: {InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeA, 0, 0},
	0x48: {InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeB, 0, 0},
	0x49: {InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeC, 0, 0},
	0x4A: {InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeD, 0, 0},
	0x4B: {InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeE, 0, 0},
	0x4C: {InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeH, 0, 0},
	0x4D: {InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeL, 0, 0},
	0x4E: {InstructionTypeLD, AddressModeR_MR, RegisterTypeC, RegisterTypeHL, 0, 0},
	0x4F: {InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeA, 0, 0},

	0x50: {InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeB, 0, 0},
	0x51: {InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeC, 0, 0},
	0x52: {InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeD, 0, 0},
	0x53: {InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeE, 0, 0},
	0x54: {InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeH, 0, 0},
	0x55: {InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeL, 0, 0},
	0x56: {InstructionTypeLD, AddressModeR_MR, RegisterTypeD, RegisterTypeHL, 0, 0},
	0x57: {InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeA, 0, 0},
	0x58: {InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeB, 0, 0},
	0x59: {InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeC, 0, 0},
	0x5A: {InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeD, 0, 0},
	0x5B: {InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeE, 0, 0},
	0x5C: {InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeH, 0, 0},
	0x5D: {InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeL, 0, 0},
	0x5E: {InstructionTypeLD, AddressModeR_MR, RegisterTypeE, RegisterTypeHL, 0, 0},
	0x5F: {InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeA, 0, 0},

	0x60: {InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeB, 0, 0},
	0x61: {InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeC, 0, 0},
	0x62: {InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeD, 0, 0},
	0x63: {InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeE, 0, 0},
	0x64: {InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeH, 0, 0},
	0x65: {InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeL, 0, 0},
	0x66: {InstructionTypeLD, AddressModeR_MR, RegisterTypeH, RegisterTypeHL, 0, 0},
	0x67: {InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeA, 0, 0},
	0x68: {InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeB, 0, 0},
	0x69: {InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeC, 0, 0},
	0x6A: {InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeD, 0, 0},
	0x6B: {InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeE, 0, 0},
	0x6C: {InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeH, 0, 0},
	0x6D: {InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeL, 0, 0},
	0x6E: {InstructionTypeLD, AddressModeR_MR, RegisterTypeL, RegisterTypeHL, 0, 0},
	0x6F: {InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeA, 0, 0},

	0x70: {InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeB, 0, 0},
	0x71: {InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeC, 0, 0},
	0x72: {InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeD, 0, 0},
	0x73: {InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeE, 0, 0},
	0x74: {InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeH, 0, 0},
	0x75: {InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeL, 0, 0},
	0x76: {InstructionTypeHALT, 0, 0, 0, 0, 0},
	0x77: {InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeA, 0, 0},
	0x78: {InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0x79: {InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0x7A: {InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0x7B: {InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0x7C: {InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0x7D: {InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0x7E: {InstructionTypeLD, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x7F: {InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},

	0x80: {InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0x81: {InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0x82: {InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0x83: {InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0x84: {InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0x85: {InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0x86: {InstructionTypeADD, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x87: {InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},
	0x88: {InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0x89: {InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0x8A: {InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0x8B: {InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0x8C: {InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0x8D: {InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0x8E: {InstructionTypeADC, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x8F: {InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},

	0x90: {InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0x91: {InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0x92: {InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0x93: {InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0x94: {InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0x95: {InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0x96: {InstructionTypeSUB, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x97: {InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},
	0x98: {InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0x99: {InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0x9A: {InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0x9B: {InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0x9C: {InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0x9D: {InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0x9E: {InstructionTypeSBC, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x9F: {InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},

	0xA0: {InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0xA1: {InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0xA2: {InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0xA3: {InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0xA4: {InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0xA5: {InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0xA6: {InstructionTypeAND, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0xA7: {InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},
	0xA8: {InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0xA9: {InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0xAA: {InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0xAB: {InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0xAC: {InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0xAD: {InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0xAE: {InstructionTypeXOR, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0xAF: {InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},

	0xB0: {InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0xB1: {InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0xB2: {InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0xB3: {InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0xB4: {InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0xB5: {InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0xB6: {InstructionTypeOR, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0xB7: {InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},
	0xB8: {InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0xB9: {InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0xBA: {InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0xBB: {InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0xBC: {InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0xBD: {InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0xBE: {InstructionTypeCP, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0xBF: {InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},

	0xC0: {InstructionTypeRET, 0, 0, 0, ConditionTypeNZ, 0},
	0xC1: {InstructionTypePOP, AddressModeR, RegisterTypeBC, 0, 0, 0},
	0xC2: {InstructionTypeJP, AddressModeN16, 0, 0, ConditionTypeNZ, 0},
	0xC3: {InstructionTypeJP, AddressModeN16, 0, 0, 0, 0},
	0xC4: {InstructionTypeCALL, AddressModeN16, 0, 0, ConditionTypeNZ, 0},
	0xC5: {InstructionTypePUSH, AddressModeR, RegisterTypeBC, 0, 0, 0},
	0xC6: {InstructionTypeADD, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xC7: {InstructionTypeRST, 0, 0, 0, 0, 0x00},
	0xC8: {InstructionTypeRET, 0, 0, 0, ConditionTypeZ, 0},
	0xC9: {InstructionTypeRET, 0, 0, 0, 0, 0},
	0xCA: {InstructionTypeJP, AddressModeN16, 0, 0, ConditionTypeZ, 0},
	0xCB: {InstructionTypeCB, AddressModeN8, 0, 0, 0, 0},
	0xCC: {InstructionTypeCALL, AddressModeN16, 0, 0, ConditionTypeZ, 0},
	0xCD: {InstructionTypeCALL, AddressModeN16, 0, 0, 0, 0},
	0xCE: {InstructionTypeSUB, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xCF: {InstructionTypeRST, 0, 0, 0, 0, 0x08},

	0xD0: {InstructionTypeRET, 0, 0, 0, ConditionTypeNC, 0},
	0xD1: {InstructionTypePOP, AddressModeR, RegisterTypeDE, 0, 0, 0},
	0xD2: {InstructionTypeJP, AddressModeN16, 0, 0, ConditionTypeNC, 0},
	0xD4: {InstructionTypeCALL, AddressModeN16, 0, 0, ConditionTypeNC, 0},
	0xD5: {InstructionTypePUSH, AddressModeR, RegisterTypeDE, 0, 0, 0},
	0xD6: {InstructionTypeSUB, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xD7: {InstructionTypeRST, 0, 0, 0, 0, 0x10},
	0xD8: {InstructionTypeRET, 0, 0, 0, ConditionTypeC, 0},
	0xD9: {InstructionTypeRETI, 0, 0, 0, 0, 0},
	0xDA: {InstructionTypeJP, AddressModeN16, 0, 0, ConditionTypeC, 0},
	0xDC: {InstructionTypeCALL, AddressModeN16, 0, 0, ConditionTypeC, 0},
	0xDE: {InstructionTypeADD, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xDF: {InstructionTypeRST, 0, 0, 0, 0, 0x18},

	0xE0: {InstructionTypeLDH, AddressModeA8_R, RegisterTypeNone, RegisterTypeA, 0, 0},
	0xE1: {InstructionTypePOP, AddressModeR, RegisterTypeHL, 0, 0, 0},
	0xE2: {InstructionTypeLD, AddressModeMR_R, RegisterTypeC, RegisterTypeA, 0, 0},
	0xE5: {InstructionTypePUSH, AddressModeR, RegisterTypeHL, 0, 0, 0},
	0xE6: {InstructionTypeAND, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xE7: {InstructionTypeRST, 0, 0, 0, 0, 0x20},
	0xE8: {InstructionTypeADD, AddressModeR_N8, RegisterTypeSP, 0, 0, 0},
	0xE9: {InstructionTypeJP, AddressModeR, RegisterTypeHL, 0, 0, 0},
	0xEA: {InstructionTypeLD, AddressModeA8_R, RegisterTypeNone, RegisterTypeA, 0, 0},
	0xEE: {InstructionTypeXOR, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xEF: {InstructionTypeRST, 0, 0, 0, 0, 0x28},

	0xF0: {InstructionTypeLDH, AddressModeR_A8, RegisterTypeA, 0, 0, 0},
	0xF1: {InstructionTypePOP, AddressModeR, RegisterTypeAF, 0, 0, 0},
	0xF2: {InstructionTypeLD, AddressModeR_MR, RegisterTypeA, RegisterTypeC, 0, 0},
	0xF3: {InstructionTypeDI, 0, 0, 0, 0, 0},
	0xF5: {InstructionTypePUSH, AddressModeR, RegisterTypeAF, 0, 0, 0},
	0xF6: {InstructionTypeOR, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xF7: {InstructionTypeRST, 0, 0, 0, 0, 0x30},
	0xF8: {InstructionTypeLD, AddressModeHL_SPR, RegisterTypeHL, RegisterTypeSP, 0, 0},
	0xF9: {InstructionTypeLD, AddressModeR_R, RegisterTypeSP, RegisterTypeHL, 0, 0},
	0xFA: {InstructionTypeLD, AddressModeR_A16, RegisterTypeA, 0, 0, 0},
	0xFB: {InstructionTypeEI, 0, 0, 0, 0, 0},
	0xFE: {InstructionTypeCP, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xFF: {InstructionTypeRST, 0, 0, 0, 0, 0x38},
}

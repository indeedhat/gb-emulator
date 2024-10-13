package emu

// https://gbdev.io/gb-opcodes/optables/

type InstructionType byte

func (i InstructionType) String() string {
	return []string{
		"None", "NOP", "LD",
		"INC", "DEC", "RLCA",
		"STOP", "RLA", "JR",
		"RRA", "DAA", "CPL",
		"SCF", "CCF", "HALT",
		"ADD", "SUB", "ADC",
		"SBC", "AND", "XOR",
		"OR", "CP", "POP",
		"JP", "PUSH", "RET", "CB",
		"CALL", "RETI", "LDH",
		"JPHL", "DI", "EI",
		"RST", "RRCA",
	}[i]
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
	InstructionTypeRRCA
)

type AddressMode byte

const (
	AddressModeNone AddressMode = iota
	AddressModeR
	AddressModeR_R
	AddressModeR_N16
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
	Type          InstructionType
	CyclesTaken   uint8
	CyclesUntaken uint8
	AddressMode   AddressMode
	Register1     RegisterType
	Register2     RegisterType
	Condition     ConditionType
	Param         byte
}

var CpuInstructions = []CpuInstriction{
	0x00: {InstructionTypeNOP, 1, 1, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0x01: {InstructionTypeLD, 3, 3, AddressModeR_N16, RegisterTypeBC, RegisterTypeNone, ConditionTypeNone, 0},
	0x02: {InstructionTypeLD, 2, 2, AddressModeMR_R, RegisterTypeBC, RegisterTypeA, ConditionTypeNone, 0},
	0x03: {InstructionTypeINC, 2, 2, AddressModeR, RegisterTypeBC, RegisterTypeNone, ConditionTypeNone, 0},
	0x04: {InstructionTypeINC, 1, 1, AddressModeR, RegisterTypeB, RegisterTypeNone, ConditionTypeNone, 0},
	0x05: {InstructionTypeDEC, 1, 1, AddressModeR, RegisterTypeB, RegisterTypeNone, ConditionTypeNone, 0},
	0x06: {InstructionTypeLD, 2, 2, AddressModeR_N8, RegisterTypeB, RegisterTypeNone, ConditionTypeNone, 0},
	0x07: {InstructionTypeRLCA, 1, 1, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0x08: {InstructionTypeLD, 5, 5, AddressModeA16_R, RegisterTypeNone, RegisterTypeSP, ConditionTypeNone, 0},
	0x09: {InstructionTypeADD, 2, 2, AddressModeR_R, RegisterTypeHL, RegisterTypeBC, ConditionTypeNone, 0},
	0x0A: {InstructionTypeLD, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeBC, ConditionTypeNone, 0},
	0x0B: {InstructionTypeDEC, 2, 2, AddressModeR, RegisterTypeBC, RegisterTypeNone, ConditionTypeNone, 0},
	0x0C: {InstructionTypeINC, 1, 1, AddressModeR, RegisterTypeC, RegisterTypeNone, ConditionTypeNone, 0},
	0x0D: {InstructionTypeDEC, 1, 1, AddressModeR, RegisterTypeC, RegisterTypeNone, ConditionTypeNone, 0},
	0x0E: {InstructionTypeLD, 2, 2, AddressModeR_N8, RegisterTypeC, RegisterTypeNone, ConditionTypeNone, 0},
	0x0F: {InstructionTypeRRCA, 1, 1, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},

	0x10: {InstructionTypeSTOP, 1, 1, AddressModeN8, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0x11: {InstructionTypeLD, 3, 3, AddressModeR_N16, RegisterTypeDE, RegisterTypeNone, ConditionTypeNone, 0},
	0x12: {InstructionTypeLD, 2, 2, AddressModeMR_R, RegisterTypeDE, RegisterTypeA, ConditionTypeNone, 0},
	0x13: {InstructionTypeINC, 2, 2, AddressModeR, RegisterTypeDE, RegisterTypeNone, ConditionTypeNone, 0},
	0x14: {InstructionTypeINC, 1, 1, AddressModeR, RegisterTypeD, RegisterTypeNone, ConditionTypeNone, 0},
	0x15: {InstructionTypeDEC, 1, 1, AddressModeR, RegisterTypeD, RegisterTypeNone, ConditionTypeNone, 0},
	0x16: {InstructionTypeLD, 2, 2, AddressModeR_N8, RegisterTypeD, RegisterTypeNone, ConditionTypeNone, 0},
	0x17: {InstructionTypeRLA, 1, 1, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0x18: {InstructionTypeJR, 3, 3, AddressModeN8, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0x19: {InstructionTypeADD, 2, 2, AddressModeR_R, RegisterTypeHL, RegisterTypeDE, ConditionTypeNone, 0},
	0x1A: {InstructionTypeLD, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeDE, ConditionTypeNone, 0},
	0x1B: {InstructionTypeDEC, 2, 2, AddressModeR, RegisterTypeDE, RegisterTypeNone, ConditionTypeNone, 0},
	0x1C: {InstructionTypeINC, 1, 1, AddressModeR, RegisterTypeE, RegisterTypeNone, ConditionTypeNone, 0},
	0x1D: {InstructionTypeDEC, 1, 1, AddressModeR, RegisterTypeE, RegisterTypeNone, ConditionTypeNone, 0},
	0x1E: {InstructionTypeLD, 2, 2, AddressModeR_N8, RegisterTypeE, RegisterTypeNone, ConditionTypeNone, 0},
	0x1F: {InstructionTypeRRA, 1, 1, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},

	0x20: {InstructionTypeJR, 3, 2, AddressModeN8, RegisterTypeNone, RegisterTypeNone, ConditionTypeNZ, 0},
	0x21: {InstructionTypeLD, 3, 3, AddressModeR_N16, RegisterTypeHL, RegisterTypeNone, ConditionTypeNone, 0},
	0x22: {InstructionTypeLD, 2, 2, AddressModeHLI_R, RegisterTypeHL, RegisterTypeA, ConditionTypeNone, 0},
	0x23: {InstructionTypeINC, 2, 2, AddressModeR, RegisterTypeHL, RegisterTypeNone, ConditionTypeNone, 0},
	0x24: {InstructionTypeINC, 1, 1, AddressModeR, RegisterTypeH, RegisterTypeNone, ConditionTypeNone, 0},
	0x25: {InstructionTypeDEC, 1, 1, AddressModeR, RegisterTypeH, RegisterTypeNone, ConditionTypeNone, 0},
	0x26: {InstructionTypeLD, 2, 2, AddressModeR_N8, RegisterTypeH, RegisterTypeNone, ConditionTypeNone, 0},
	0x27: {InstructionTypeDAA, 1, 1, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0x28: {InstructionTypeJR, 3, 2, AddressModeN8, RegisterTypeNone, RegisterTypeNone, ConditionTypeZ, 0},
	0x29: {InstructionTypeADD, 2, 2, AddressModeR_R, RegisterTypeHL, RegisterTypeHL, ConditionTypeNone, 0},
	0x2A: {InstructionTypeLD, 2, 2, AddressModeR_HLI, RegisterTypeA, RegisterTypeHL, ConditionTypeNone, 0},
	0x2B: {InstructionTypeDEC, 2, 2, AddressModeR, RegisterTypeHL, RegisterTypeNone, ConditionTypeNone, 0},
	0x2C: {InstructionTypeINC, 1, 1, AddressModeR, RegisterTypeL, RegisterTypeNone, ConditionTypeNone, 0},
	0x2D: {InstructionTypeDEC, 1, 1, AddressModeR, RegisterTypeL, RegisterTypeNone, ConditionTypeNone, 0},
	0x2E: {InstructionTypeLD, 2, 2, AddressModeR_N8, RegisterTypeL, RegisterTypeNone, ConditionTypeNone, 0},
	0x2F: {InstructionTypeCPL, 1, 1, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},

	0x30: {InstructionTypeJR, 3, 2, AddressModeN8, RegisterTypeNone, RegisterTypeNone, ConditionTypeNC, 0},
	0x31: {InstructionTypeLD, 3, 3, AddressModeR_N16, RegisterTypeSP, RegisterTypeNone, ConditionTypeNone, 0},
	0x32: {InstructionTypeLD, 2, 2, AddressModeHLD_R, RegisterTypeHL, RegisterTypeA, ConditionTypeNone, 0},
	0x33: {InstructionTypeINC, 2, 2, AddressModeR, RegisterTypeSP, RegisterTypeNone, ConditionTypeNone, 0},
	0x34: {InstructionTypeINC, 3, 3, AddressModeMR, RegisterTypeHL, RegisterTypeNone, ConditionTypeNone, 0},
	0x35: {InstructionTypeDEC, 3, 3, AddressModeMR, RegisterTypeHL, RegisterTypeNone, ConditionTypeNone, 0},
	0x36: {InstructionTypeLD, 3, 3, AddressModeMR_N8, RegisterTypeHL, RegisterTypeNone, ConditionTypeNone, 0},
	0x37: {InstructionTypeSCF, 1, 1, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0x38: {InstructionTypeJR, 3, 2, AddressModeN8, RegisterTypeNone, RegisterTypeNone, ConditionTypeC, 0},
	0x39: {InstructionTypeADD, 2, 2, AddressModeR_R, RegisterTypeHL, RegisterTypeSP, ConditionTypeNone, 0},
	0x3A: {InstructionTypeLD, 2, 2, AddressModeR_HLD, RegisterTypeA, RegisterTypeHL, ConditionTypeNone, 0},
	0x3B: {InstructionTypeDEC, 2, 2, AddressModeR, RegisterTypeSP, RegisterTypeNone, ConditionTypeNone, 0},
	0x3C: {InstructionTypeINC, 1, 1, AddressModeR, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0x3D: {InstructionTypeDEC, 1, 1, AddressModeR, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0x3E: {InstructionTypeLD, 2, 2, AddressModeR_N8, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0x3F: {InstructionTypeCCF, 1, 1, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},

	0x40: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeB, RegisterTypeB, ConditionTypeNone, 0},
	0x41: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeB, RegisterTypeC, ConditionTypeNone, 0},
	0x42: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeB, RegisterTypeD, ConditionTypeNone, 0},
	0x43: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeB, RegisterTypeE, ConditionTypeNone, 0},
	0x44: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeB, RegisterTypeH, ConditionTypeNone, 0},
	0x45: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeB, RegisterTypeL, ConditionTypeNone, 0},
	0x46: {InstructionTypeLD, 2, 2, AddressModeR_MR, RegisterTypeB, RegisterTypeHL, ConditionTypeNone, 0},
	0x47: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeB, RegisterTypeA, ConditionTypeNone, 0},
	0x48: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeC, RegisterTypeB, ConditionTypeNone, 0},
	0x49: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeC, RegisterTypeC, ConditionTypeNone, 0},
	0x4A: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeC, RegisterTypeD, ConditionTypeNone, 0},
	0x4B: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeC, RegisterTypeE, ConditionTypeNone, 0},
	0x4C: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeC, RegisterTypeH, ConditionTypeNone, 0},
	0x4D: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeC, RegisterTypeL, ConditionTypeNone, 0},
	0x4E: {InstructionTypeLD, 2, 2, AddressModeR_MR, RegisterTypeC, RegisterTypeHL, ConditionTypeNone, 0},
	0x4F: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeC, RegisterTypeA, ConditionTypeNone, 0},

	0x50: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeD, RegisterTypeB, ConditionTypeNone, 0},
	0x51: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeD, RegisterTypeC, ConditionTypeNone, 0},
	0x52: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeD, RegisterTypeD, ConditionTypeNone, 0},
	0x53: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeD, RegisterTypeE, ConditionTypeNone, 0},
	0x54: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeD, RegisterTypeH, ConditionTypeNone, 0},
	0x55: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeD, RegisterTypeL, ConditionTypeNone, 0},
	0x56: {InstructionTypeLD, 2, 2, AddressModeR_MR, RegisterTypeD, RegisterTypeHL, ConditionTypeNone, 0},
	0x57: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeD, RegisterTypeA, ConditionTypeNone, 0},
	0x58: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeE, RegisterTypeB, ConditionTypeNone, 0},
	0x59: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeE, RegisterTypeC, ConditionTypeNone, 0},
	0x5A: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeE, RegisterTypeD, ConditionTypeNone, 0},
	0x5B: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeE, RegisterTypeE, ConditionTypeNone, 0},
	0x5C: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeE, RegisterTypeH, ConditionTypeNone, 0},
	0x5D: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeE, RegisterTypeL, ConditionTypeNone, 0},
	0x5E: {InstructionTypeLD, 2, 2, AddressModeR_MR, RegisterTypeE, RegisterTypeHL, ConditionTypeNone, 0},
	0x5F: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeE, RegisterTypeA, ConditionTypeNone, 0},

	0x60: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeH, RegisterTypeB, ConditionTypeNone, 0},
	0x61: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeH, RegisterTypeC, ConditionTypeNone, 0},
	0x62: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeH, RegisterTypeD, ConditionTypeNone, 0},
	0x63: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeH, RegisterTypeE, ConditionTypeNone, 0},
	0x64: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeH, RegisterTypeH, ConditionTypeNone, 0},
	0x65: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeH, RegisterTypeL, ConditionTypeNone, 0},
	0x66: {InstructionTypeLD, 2, 2, AddressModeR_MR, RegisterTypeH, RegisterTypeHL, ConditionTypeNone, 0},
	0x67: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeH, RegisterTypeA, ConditionTypeNone, 0},
	0x68: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeL, RegisterTypeB, ConditionTypeNone, 0},
	0x69: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeL, RegisterTypeC, ConditionTypeNone, 0},
	0x6A: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeL, RegisterTypeD, ConditionTypeNone, 0},
	0x6B: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeL, RegisterTypeE, ConditionTypeNone, 0},
	0x6C: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeL, RegisterTypeH, ConditionTypeNone, 0},
	0x6D: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeL, RegisterTypeL, ConditionTypeNone, 0},
	0x6E: {InstructionTypeLD, 2, 2, AddressModeR_MR, RegisterTypeL, RegisterTypeHL, ConditionTypeNone, 0},
	0x6F: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeL, RegisterTypeA, ConditionTypeNone, 0},

	0x70: {InstructionTypeLD, 2, 2, AddressModeMR_R, RegisterTypeHL, RegisterTypeB, ConditionTypeNone, 0},
	0x71: {InstructionTypeLD, 2, 2, AddressModeMR_R, RegisterTypeHL, RegisterTypeC, ConditionTypeNone, 0},
	0x72: {InstructionTypeLD, 2, 2, AddressModeMR_R, RegisterTypeHL, RegisterTypeD, ConditionTypeNone, 0},
	0x73: {InstructionTypeLD, 2, 2, AddressModeMR_R, RegisterTypeHL, RegisterTypeE, ConditionTypeNone, 0},
	0x74: {InstructionTypeLD, 2, 2, AddressModeMR_R, RegisterTypeHL, RegisterTypeH, ConditionTypeNone, 0},
	0x75: {InstructionTypeLD, 2, 2, AddressModeMR_R, RegisterTypeHL, RegisterTypeL, ConditionTypeNone, 0},
	0x76: {InstructionTypeHALT, 1, 1, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0x77: {InstructionTypeLD, 2, 2, AddressModeMR_R, RegisterTypeHL, RegisterTypeA, ConditionTypeNone, 0},
	0x78: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeB, ConditionTypeNone, 0},
	0x79: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeC, ConditionTypeNone, 0},
	0x7A: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeD, ConditionTypeNone, 0},
	0x7B: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeE, ConditionTypeNone, 0},
	0x7C: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeH, ConditionTypeNone, 0},
	0x7D: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeL, ConditionTypeNone, 0},
	0x7E: {InstructionTypeLD, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, ConditionTypeNone, 0},
	0x7F: {InstructionTypeLD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeA, ConditionTypeNone, 0},

	0x80: {InstructionTypeADD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeB, ConditionTypeNone, 0},
	0x81: {InstructionTypeADD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeC, ConditionTypeNone, 0},
	0x82: {InstructionTypeADD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeD, ConditionTypeNone, 0},
	0x83: {InstructionTypeADD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeE, ConditionTypeNone, 0},
	0x84: {InstructionTypeADD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeH, ConditionTypeNone, 0},
	0x85: {InstructionTypeADD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeL, ConditionTypeNone, 0},
	0x86: {InstructionTypeADD, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, ConditionTypeNone, 0},
	0x87: {InstructionTypeADD, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeA, ConditionTypeNone, 0},
	0x88: {InstructionTypeADC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeB, ConditionTypeNone, 0},
	0x89: {InstructionTypeADC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeC, ConditionTypeNone, 0},
	0x8A: {InstructionTypeADC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeD, ConditionTypeNone, 0},
	0x8B: {InstructionTypeADC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeE, ConditionTypeNone, 0},
	0x8C: {InstructionTypeADC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeH, ConditionTypeNone, 0},
	0x8D: {InstructionTypeADC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeL, ConditionTypeNone, 0},
	0x8E: {InstructionTypeADC, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, ConditionTypeNone, 0},
	0x8F: {InstructionTypeADC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeA, ConditionTypeNone, 0},

	0x90: {InstructionTypeSUB, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeB, ConditionTypeNone, 0},
	0x91: {InstructionTypeSUB, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeC, ConditionTypeNone, 0},
	0x92: {InstructionTypeSUB, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeD, ConditionTypeNone, 0},
	0x93: {InstructionTypeSUB, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeE, ConditionTypeNone, 0},
	0x94: {InstructionTypeSUB, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeH, ConditionTypeNone, 0},
	0x95: {InstructionTypeSUB, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeL, ConditionTypeNone, 0},
	0x96: {InstructionTypeSUB, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, ConditionTypeNone, 0},
	0x97: {InstructionTypeSUB, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeA, ConditionTypeNone, 0},
	0x98: {InstructionTypeSBC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeB, ConditionTypeNone, 0},
	0x99: {InstructionTypeSBC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeC, ConditionTypeNone, 0},
	0x9A: {InstructionTypeSBC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeD, ConditionTypeNone, 0},
	0x9B: {InstructionTypeSBC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeE, ConditionTypeNone, 0},
	0x9C: {InstructionTypeSBC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeH, ConditionTypeNone, 0},
	0x9D: {InstructionTypeSBC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeL, ConditionTypeNone, 0},
	0x9E: {InstructionTypeSBC, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, ConditionTypeNone, 0},
	0x9F: {InstructionTypeSBC, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeA, ConditionTypeNone, 0},

	0xA0: {InstructionTypeAND, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeB, ConditionTypeNone, 0},
	0xA1: {InstructionTypeAND, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeC, ConditionTypeNone, 0},
	0xA2: {InstructionTypeAND, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeD, ConditionTypeNone, 0},
	0xA3: {InstructionTypeAND, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeE, ConditionTypeNone, 0},
	0xA4: {InstructionTypeAND, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeH, ConditionTypeNone, 0},
	0xA5: {InstructionTypeAND, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeL, ConditionTypeNone, 0},
	0xA6: {InstructionTypeAND, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, ConditionTypeNone, 0},
	0xA7: {InstructionTypeAND, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeA, ConditionTypeNone, 0},
	0xA8: {InstructionTypeXOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeB, ConditionTypeNone, 0},
	0xA9: {InstructionTypeXOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeC, ConditionTypeNone, 0},
	0xAA: {InstructionTypeXOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeD, ConditionTypeNone, 0},
	0xAB: {InstructionTypeXOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeE, ConditionTypeNone, 0},
	0xAC: {InstructionTypeXOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeH, ConditionTypeNone, 0},
	0xAD: {InstructionTypeXOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeL, ConditionTypeNone, 0},
	0xAE: {InstructionTypeXOR, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, ConditionTypeNone, 0},
	0xAF: {InstructionTypeXOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeA, ConditionTypeNone, 0},

	0xB0: {InstructionTypeOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeB, ConditionTypeNone, 0},
	0xB1: {InstructionTypeOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeC, ConditionTypeNone, 0},
	0xB2: {InstructionTypeOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeD, ConditionTypeNone, 0},
	0xB3: {InstructionTypeOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeE, ConditionTypeNone, 0},
	0xB4: {InstructionTypeOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeH, ConditionTypeNone, 0},
	0xB5: {InstructionTypeOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeL, ConditionTypeNone, 0},
	0xB6: {InstructionTypeOR, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, ConditionTypeNone, 0},
	0xB7: {InstructionTypeOR, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeA, ConditionTypeNone, 0},
	0xB8: {InstructionTypeCP, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeB, ConditionTypeNone, 0},
	0xB9: {InstructionTypeCP, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeC, ConditionTypeNone, 0},
	0xBA: {InstructionTypeCP, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeD, ConditionTypeNone, 0},
	0xBB: {InstructionTypeCP, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeE, ConditionTypeNone, 0},
	0xBC: {InstructionTypeCP, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeH, ConditionTypeNone, 0},
	0xBD: {InstructionTypeCP, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeL, ConditionTypeNone, 0},
	0xBE: {InstructionTypeCP, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, ConditionTypeNone, 0},
	0xBF: {InstructionTypeCP, 1, 1, AddressModeR_R, RegisterTypeA, RegisterTypeA, ConditionTypeNone, 0},

	0xC0: {InstructionTypeRET, 5, 2, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNZ, 0},
	0xC1: {InstructionTypePOP, 3, 3, AddressModeR, RegisterTypeBC, RegisterTypeNone, ConditionTypeNone, 0},
	0xC2: {InstructionTypeJP, 4, 3, AddressModeN16, RegisterTypeNone, RegisterTypeNone, ConditionTypeNZ, 0},
	0xC3: {InstructionTypeJP, 4, 4, AddressModeN16, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0xC4: {InstructionTypeCALL, 6, 3, AddressModeN16, RegisterTypeNone, RegisterTypeNone, ConditionTypeNZ, 0},
	0xC5: {InstructionTypePUSH, 4, 4, AddressModeR, RegisterTypeBC, RegisterTypeNone, ConditionTypeNone, 0},
	0xC6: {InstructionTypeADD, 2, 2, AddressModeR_N8, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0xC7: {InstructionTypeRST, 4, 4, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0x00},
	0xC8: {InstructionTypeRET, 5, 2, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeZ, 0},
	0xC9: {InstructionTypeRET, 4, 4, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0xCA: {InstructionTypeJP, 4, 3, AddressModeN16, RegisterTypeNone, RegisterTypeNone, ConditionTypeZ, 0},
	0xCB: {InstructionTypeCB, 1, 1, AddressModeN8, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0xCC: {InstructionTypeCALL, 6, 3, AddressModeN16, RegisterTypeNone, RegisterTypeNone, ConditionTypeZ, 0},
	0xCD: {InstructionTypeCALL, 6, 6, AddressModeN16, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0xCE: {InstructionTypeADC, 2, 2, AddressModeR_N8, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0xCF: {InstructionTypeRST, 4, 4, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0x08},

	0xD0: {InstructionTypeRET, 5, 2, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNC, 0},
	0xD1: {InstructionTypePOP, 3, 3, AddressModeR, RegisterTypeDE, RegisterTypeNone, ConditionTypeNone, 0},
	0xD2: {InstructionTypeJP, 4, 3, AddressModeN16, RegisterTypeNone, RegisterTypeNone, ConditionTypeNC, 0},
	0xD4: {InstructionTypeCALL, 6, 3, AddressModeN16, RegisterTypeNone, RegisterTypeNone, ConditionTypeNC, 0},
	0xD5: {InstructionTypePUSH, 4, 4, AddressModeR, RegisterTypeDE, RegisterTypeNone, ConditionTypeNone, 0},
	0xD6: {InstructionTypeSUB, 2, 2, AddressModeR_N8, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0xD7: {InstructionTypeRST, 4, 4, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0x10},
	0xD8: {InstructionTypeRET, 5, 2, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeC, 0},
	0xD9: {InstructionTypeRETI, 4, 4, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0xDA: {InstructionTypeJP, 4, 3, AddressModeN16, RegisterTypeNone, RegisterTypeNone, ConditionTypeC, 0},
	0xDC: {InstructionTypeCALL, 6, 3, AddressModeN16, RegisterTypeNone, RegisterTypeNone, ConditionTypeC, 0},
	0xDE: {InstructionTypeSBC, 2, 2, AddressModeR_N8, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0xDF: {InstructionTypeRST, 4, 4, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0x18},

	0xE0: {InstructionTypeLDH, 3, 3, AddressModeA8_R, RegisterTypeNone, RegisterTypeA, ConditionTypeNone, 0},
	0xE1: {InstructionTypePOP, 3, 3, AddressModeR, RegisterTypeHL, RegisterTypeNone, ConditionTypeNone, 0},
	0xE2: {InstructionTypeLD, 2, 2, AddressModeMR_R, RegisterTypeC, RegisterTypeA, ConditionTypeNone, 0},
	0xE5: {InstructionTypePUSH, 4, 4, AddressModeR, RegisterTypeHL, RegisterTypeNone, ConditionTypeNone, 0},
	0xE6: {InstructionTypeAND, 2, 2, AddressModeR_N8, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0xE7: {InstructionTypeRST, 4, 4, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0x20},
	0xE8: {InstructionTypeADD, 4, 4, AddressModeR_N8, RegisterTypeSP, RegisterTypeNone, ConditionTypeNone, 0},
	0xE9: {InstructionTypeJP, 1, 1, AddressModeR, RegisterTypeHL, RegisterTypeNone, ConditionTypeNone, 0},
	0xEA: {InstructionTypeLD, 4, 4, AddressModeA16_R, RegisterTypeNone, RegisterTypeA, ConditionTypeNone, 0},
	0xEE: {InstructionTypeXOR, 2, 2, AddressModeR_N8, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0xEF: {InstructionTypeRST, 4, 4, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0x28},

	0xF0: {InstructionTypeLDH, 3, 3, AddressModeR_A8, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0xF1: {InstructionTypePOP, 3, 3, AddressModeR, RegisterTypeAF, RegisterTypeNone, ConditionTypeNone, 0},
	0xF2: {InstructionTypeLD, 2, 2, AddressModeR_MR, RegisterTypeA, RegisterTypeC, ConditionTypeNone, 0},
	0xF3: {InstructionTypeDI, 1, 1, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0xF5: {InstructionTypePUSH, 4, 4, AddressModeR, RegisterTypeAF, RegisterTypeNone, ConditionTypeNone, 0},
	0xF6: {InstructionTypeOR, 2, 2, AddressModeR_N8, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0xF7: {InstructionTypeRST, 4, 4, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0x30},
	0xF8: {InstructionTypeLD, 3, 3, AddressModeHL_SPR, RegisterTypeHL, RegisterTypeSP, ConditionTypeNone, 0},
	0xF9: {InstructionTypeLD, 2, 2, AddressModeR_R, RegisterTypeSP, RegisterTypeHL, ConditionTypeNone, 0},
	0xFA: {InstructionTypeLD, 4, 4, AddressModeR_A16, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0xFB: {InstructionTypeEI, 1, 2, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0},
	0xFE: {InstructionTypeCP, 2, 2, AddressModeR_N8, RegisterTypeA, RegisterTypeNone, ConditionTypeNone, 0},
	0xFF: {InstructionTypeRST, 4, 4, AddressModeNone, RegisterTypeNone, RegisterTypeNone, ConditionTypeNone, 0x38},
}

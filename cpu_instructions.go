package main

import "image/jpeg"

type InstructionType byte

const (
	InstructionTypeNone = iota
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
	InstructionTypeERR
	InstructionTypePREFIX

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
	AddressModeImplied = iota
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
	AddressModeN16_R
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

const (
	RegisterTypeNone = iota
	RegisterTypeA
	RegisterTypeF
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
	ConditionTypeNone = iota
	ConditionTypeNZ
	ConditionTypeZ
	ConditionTypeNC
	ConditionTypeC
)

type CpuInstriction struct {
	Type        byte
	AddressMode byte
	Register1   byte
	Register2   byte
	Condition   byte
	Param       byte
}

var CpuInstructions = []CpuInstriction{
	0x00: CpuInstriction{InstructionTypeNOP, 0, 0, 0, 0, 0},
	0x01: CpuInstriction{InstructionTypeLD, AddressModeR_N16, RegisterTypeBC, 0, 0, 0},
	0x02: CpuInstriction{InstructionTypeLD, AddressModeMR_N16, RegisterTypeBC, RegisterTypeA, 0, 0},
	0x03: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeBC, 0, 0, 0},
	0x04: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeB, 0, 0, 0},
	0x05: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeB, 0, 0, 0},
	0x06: CpuInstriction{InstructionTypeLD, AddressModeR_N8, RegisterTypeB, 0, 0, 0},
	0x07: CpuInstriction{InstructionTypeRLCA, 0, 0, 0, 0, 0},
	0x08: CpuInstriction{InstructionTypeLD, AddressModeA16_R, RegisterTypeNone, RegisterTypeSP, 0, 0},
	0x09: CpuInstriction{InstructionTypeADD, AddressModeR_R, RegisterTypeHL, RegisterTypeBC, 0, 0},
	0x0A: CpuInstriction{InstructionTypeLD, AddressModeR_MR, RegisterTypeA, RegisterTypeBC, 0, 0},
	0x0B: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeBC, 0, 0, 0},
	0x0C: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeC, 0, 0, 0},
	0x0D: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeC, 0, 0, 0},
	0x0E: CpuInstriction{InstructionTypeLD, AddressModeR_N8, RegisterTypeC, 0, 0, 0},
	0x0F: CpuInstriction{InstructionTypeRRCA, 0, 0, 0, 0, 0},

	0x10: CpuInstriction{InstructionTypeSTOP, AddressModeN8, 0, 0, 0, 0},
	0x11: CpuInstriction{InstructionTypeLD, AddressModeR_N16, RegisterTypeDE, 0, 0, 0},
	0x12: CpuInstriction{InstructionTypeLD, AddressModeMR_N16, RegisterTypeDE, RegisterTypeA, 0, 0},
	0x13: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeDE, 0, 0, 0},
	0x14: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeD, 0, 0, 0},
	0x15: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeD, 0, 0, 0},
	0x16: CpuInstriction{InstructionTypeLD, AddressModeMR_N8, RegisterTypeD, 0, 0, 0},
	0x17: CpuInstriction{InstructionTypeRLA, 0, 0, 0, 0, 0},
	0x18: CpuInstriction{InstructionTypeJR, AddressModeN8, 0, 0, 0, 0},
	0x19: CpuInstriction{InstructionTypeADD, AddressModeR_R, RegisterTypeHL, RegisterTypeDE, 0, 0},
	0x1A: CpuInstriction{InstructionTypeLD, AddressModeR_MR, RegisterTypeA, RegisterTypeDE, 0, 0},
	0x1B: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeDE, 0, 0, 0},
	0x1C: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeE, 0, 0, 0},
	0x1D: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeE, 0, 0, 0},
	0x1E: CpuInstriction{InstructionTypeLD, AddressModeR_N8, RegisterTypeE, 0, 0, 0},
	0x1F: CpuInstriction{InstructionTypeRRA, 0, 0, 0, 0, 0},

	0x20: CpuInstriction{InstructionTypeJR, AddressModeN8, 0, 0, 0, ConditionTypeNZ},
	0x21: CpuInstriction{InstructionTypeLD, AddressModeR_N16, RegisterTypeHL, 0, 0, 0},
	0x22: CpuInstriction{InstructionTypeLD, AddressModeHLI_R, RegisterTypeHL, RegisterTypeA, 0, 0},
	0x23: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeHL, 0, 0, 0},
	0x24: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeH, 0, 0, 0},
	0x25: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeH, 0, 0, 0},
	0x26: CpuInstriction{InstructionTypeLD, AddressModeR_N8, RegisterTypeH, 0, 0, 0},
	0x27: CpuInstriction{InstructionTypeDAA, 0, 0, 0, 0, 0},
	0x28: CpuInstriction{InstructionTypeJR, AddressModeN8, 0, 0, 0, ConditionTypeZ},
	0x29: CpuInstriction{InstructionTypeADD, AddressModeR_R, RegisterTypeHL, RegisterTypeHL, 0, 0},
	0x2A: CpuInstriction{InstructionTypeLD, AddressModeR_HLI, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x2B: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeHL, 0, 0, 0},
	0x2C: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeL, 0, 0, 0},
	0x2D: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeL, 0, 0, 0},
	0x2E: CpuInstriction{InstructionTypeLD, AddressModeR_N8, RegisterTypeL, 0, 0, 0},
	0x2F: CpuInstriction{InstructionTypeCPL, 0, 0, 0, 0, 0},

	0x30: CpuInstriction{InstructionTypeJR, AddressModeN8, 0, 0, 0, ConditionTypeNC},
	0x31: CpuInstriction{InstructionTypeLD, AddressModeR_N16, RegisterTypeSP, 0, 0, 0},
	0x32: CpuInstriction{InstructionTypeLD, AddressModeHLD_R, RegisterTypeHL, RegisterTypeA, 0, 0},
	0x33: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeSP, 0, 0, 0},
	0x34: CpuInstriction{InstructionTypeINC, AddressModeMR, RegisterTypeHL, 0, 0, 0},
	0x35: CpuInstriction{InstructionTypeDEC, AddressModeMR, RegisterTypeHL, 0, 0, 0},
	0x36: CpuInstriction{InstructionTypeLD, AddressModeMR_N8, RegisterTypeHL, 0, 0, 0},
	0x37: CpuInstriction{InstructionTypeSCF, 0, 0, 0, 0, 0},
	0x38: CpuInstriction{InstructionTypeJR, AddressModeN8, 0, 0, 0, ConditionTypeC},
	0x39: CpuInstriction{InstructionTypeADD, AddressModeR_R, RegisterTypeHL, RegisterTypeSP, 0, 0},
	0x3A: CpuInstriction{InstructionTypeLD, AddressModeR_HLD, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x3B: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeSP, 0, 0, 0},
	0x3C: CpuInstriction{InstructionTypeINC, AddressModeR, RegisterTypeA, 0, 0, 0},
	0x3D: CpuInstriction{InstructionTypeDEC, AddressModeR, RegisterTypeA, 0, 0, 0},
	0x3E: CpuInstriction{InstructionTypeLD, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0x3F: CpuInstriction{InstructionTypeCCF, 0, 0, 0, 0, 0},

	0x40: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeB, 0, 0},
	0x41: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeC, 0, 0},
	0x42: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeD, 0, 0},
	0x43: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeE, 0, 0},
	0x44: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeH, 0, 0},
	0x45: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeL, 0, 0},
	0x46: CpuInstriction{InstructionTypeLD, AddressModeR_MR, RegisterTypeB, RegisterTypeHL, 0, 0},
	0x47: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeB, RegisterTypeA, 0, 0},
	0x48: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeB, 0, 0},
	0x49: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeC, 0, 0},
	0x4A: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeD, 0, 0},
	0x4B: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeE, 0, 0},
	0x4C: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeH, 0, 0},
	0x4D: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeL, 0, 0},
	0x4E: CpuInstriction{InstructionTypeLD, AddressModeR_MR, RegisterTypeC, RegisterTypeHL, 0, 0},
	0x4F: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeC, RegisterTypeA, 0, 0},

	0x50: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeB, 0, 0},
	0x51: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeC, 0, 0},
	0x52: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeD, 0, 0},
	0x53: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeE, 0, 0},
	0x54: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeH, 0, 0},
	0x55: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeL, 0, 0},
	0x56: CpuInstriction{InstructionTypeLD, AddressModeR_MR, RegisterTypeD, RegisterTypeHL, 0, 0},
	0x57: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeD, RegisterTypeA, 0, 0},
	0x58: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeB, 0, 0},
	0x59: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeC, 0, 0},
	0x5A: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeD, 0, 0},
	0x5B: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeE, 0, 0},
	0x5C: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeH, 0, 0},
	0x5D: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeL, 0, 0},
	0x5E: CpuInstriction{InstructionTypeLD, AddressModeR_MR, RegisterTypeE, RegisterTypeHL, 0, 0},
	0x5F: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeE, RegisterTypeA, 0, 0},

	0x60: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeB, 0, 0},
	0x61: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeC, 0, 0},
	0x62: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeD, 0, 0},
	0x63: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeE, 0, 0},
	0x64: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeH, 0, 0},
	0x65: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeL, 0, 0},
	0x66: CpuInstriction{InstructionTypeLD, AddressModeR_MR, RegisterTypeH, RegisterTypeHL, 0, 0},
	0x67: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeH, RegisterTypeA, 0, 0},
	0x68: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeB, 0, 0},
	0x69: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeC, 0, 0},
	0x6A: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeD, 0, 0},
	0x6B: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeE, 0, 0},
	0x6C: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeH, 0, 0},
	0x6D: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeL, 0, 0},
	0x6E: CpuInstriction{InstructionTypeLD, AddressModeR_MR, RegisterTypeL, RegisterTypeHL, 0, 0},
	0x6F: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeL, RegisterTypeA, 0, 0},

	0x70: CpuInstriction{InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeB, 0, 0},
	0x71: CpuInstriction{InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeC, 0, 0},
	0x72: CpuInstriction{InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeD, 0, 0},
	0x73: CpuInstriction{InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeE, 0, 0},
	0x74: CpuInstriction{InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeH, 0, 0},
	0x75: CpuInstriction{InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeL, 0, 0},
	0x76: CpuInstriction{InstructionTypeHALT, 0, 0, 0, 0, 0},
	0x77: CpuInstriction{InstructionTypeLD, AddressModeMR_R, RegisterTypeHL, RegisterTypeA, 0, 0},
	0x78: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0x79: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0x7A: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0x7B: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0x7C: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0x7D: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0x7E: CpuInstriction{InstructionTypeLD, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x7F: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},

	0x80: CpuInstriction{InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0x81: CpuInstriction{InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0x82: CpuInstriction{InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0x83: CpuInstriction{InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0x84: CpuInstriction{InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0x85: CpuInstriction{InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0x86: CpuInstriction{InstructionTypeADD, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x87: CpuInstriction{InstructionTypeADD, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},
	0x88: CpuInstriction{InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0x89: CpuInstriction{InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0x8A: CpuInstriction{InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0x8B: CpuInstriction{InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0x8C: CpuInstriction{InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0x8D: CpuInstriction{InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0x8E: CpuInstriction{InstructionTypeADC, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x8F: CpuInstriction{InstructionTypeADC, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},

	0x90: CpuInstriction{InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0x91: CpuInstriction{InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0x92: CpuInstriction{InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0x93: CpuInstriction{InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0x94: CpuInstriction{InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0x95: CpuInstriction{InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0x96: CpuInstriction{InstructionTypeSUB, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x97: CpuInstriction{InstructionTypeSUB, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},
	0x98: CpuInstriction{InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0x99: CpuInstriction{InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0x9A: CpuInstriction{InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0x9B: CpuInstriction{InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0x9C: CpuInstriction{InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0x9D: CpuInstriction{InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0x9E: CpuInstriction{InstructionTypeSBC, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0x9F: CpuInstriction{InstructionTypeSBC, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},

	0xA0: CpuInstriction{InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0xA1: CpuInstriction{InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0xA2: CpuInstriction{InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0xA3: CpuInstriction{InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0xA4: CpuInstriction{InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0xA5: CpuInstriction{InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0xA6: CpuInstriction{InstructionTypeAND, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0xA7: CpuInstriction{InstructionTypeAND, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},
	0xA8: CpuInstriction{InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0xA9: CpuInstriction{InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0xAA: CpuInstriction{InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0xAB: CpuInstriction{InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0xAC: CpuInstriction{InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0xAD: CpuInstriction{InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0xAE: CpuInstriction{InstructionTypeXOR, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0xAF: CpuInstriction{InstructionTypeXOR, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},

	0xB0: CpuInstriction{InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0xB1: CpuInstriction{InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0xB2: CpuInstriction{InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0xB3: CpuInstriction{InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0xB4: CpuInstriction{InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0xB5: CpuInstriction{InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0xB6: CpuInstriction{InstructionTypeOR, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0xB7: CpuInstriction{InstructionTypeOR, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},
	0xB8: CpuInstriction{InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeB, 0, 0},
	0xB9: CpuInstriction{InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeC, 0, 0},
	0xBA: CpuInstriction{InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeD, 0, 0},
	0xBB: CpuInstriction{InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeE, 0, 0},
	0xBC: CpuInstriction{InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeH, 0, 0},
	0xBD: CpuInstriction{InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeL, 0, 0},
	0xBE: CpuInstriction{InstructionTypeCP, AddressModeR_MR, RegisterTypeA, RegisterTypeHL, 0, 0},
	0xBF: CpuInstriction{InstructionTypeCP, AddressModeR_R, RegisterTypeA, RegisterTypeA, 0, 0},

	0xC0: CpuInstriction{InstructionTypeRET, 0, 0, 0, 0, ConditionTypeNZ},
	0xC1: CpuInstriction{InstructionTypePOP, AddressModeR, RegisterTypeBC, 0, 0, 0},
	0xC2: CpuInstriction{InstructionTypeJP, AddressModeN16, 0, 0, 0, ConditionTypeNZ},
	0xC3: CpuInstriction{InstructionTypeJP, AddressModeN16, 0, 0, 0, 0},
	0xC4: CpuInstriction{InstructionTypeCALL, AddressModeN16, 0, 0, 0, ConditionTypeNZ},
	0xC5: CpuInstriction{InstructionTypePUSH, AddressModeR, RegisterTypeBC, 0, 0, 0},
	0xC6: CpuInstriction{InstructionTypeADD, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xC7: CpuInstriction{InstructionTypeRST, 0, 0, 0, 0, 0x00},
	0xC8: CpuInstriction{InstructionTypeRET, 0, 0, 0, 0, ConditionTypeZ},
	0xC9: CpuInstriction{InstructionTypeRET, 0, 0, 0, 0, 0},
	0xCA: CpuInstriction{InstructionTypeJP, AddressModeN16, 0, 0, 0, ConditionTypeZ},
	0xCB: CpuInstriction{InstructionTypePREFIX, 0, 0, 0, 0, 0},
	0xCC: CpuInstriction{InstructionTypeCALL, AddressModeN16, 0, 0, 0, ConditionTypeZ},
	0xCD: CpuInstriction{InstructionTypeCALL, AddressModeN16, 0, 0, 0, 0},
	0xCE: CpuInstriction{InstructionTypeSUB, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xCF: CpuInstriction{InstructionTypeRST, 0, 0, 0, 0, 0x08},

	0xD0: CpuInstriction{InstructionTypeRET, 0, 0, 0, 0, ConditionTypeNC},
	0xD1: CpuInstriction{InstructionTypePOP, AddressModeR, RegisterTypeDE, 0, 0, 0},
	0xD2: CpuInstriction{InstructionTypeJP, AddressModeN16, 0, 0, 0, ConditionTypeNC},
	0xD4: CpuInstriction{InstructionTypeCALL, AddressModeN16, 0, 0, 0, ConditionTypeNC},
	0xD5: CpuInstriction{InstructionTypePUSH, AddressModeR, RegisterTypeDE, 0, 0, 0},
	0xD6: CpuInstriction{InstructionTypeSUB, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xD7: CpuInstriction{InstructionTypeRST, 0, 0, 0, 0, 0x10},
	0xD8: CpuInstriction{InstructionTypeRET, 0, 0, 0, 0, ConditionTypeC},
	0xD9: CpuInstriction{InstructionTypeRETI, 0, 0, 0, 0, 0},
	0xDA: CpuInstriction{InstructionTypeJP, AddressModeN16, 0, 0, 0, ConditionTypeC},
	0xDC: CpuInstriction{InstructionTypeCALL, AddressModeN16, 0, 0, 0, ConditionTypeC},
	0xDE: CpuInstriction{InstructionTypeADD, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xDF: CpuInstriction{InstructionTypeRST, 0, 0, 0, 0, 0x18},

	0xE0: CpuInstriction{InstructionTypeLDH, AddressModeA8_R, RegisterTypeNone, RegisterTypeA, 0, 0},
	0xE1: CpuInstriction{InstructionTypePOP, AddressModeR, RegisterTypeHL, 0, 0, 0},
	0xE2: CpuInstriction{InstructionTypeLD, AddressModeMR_R, RegisterTypeC, RegisterTypeA, 0, 0},
	0xE5: CpuInstriction{InstructionTypePUSH, AddressModeR, RegisterTypeHL, 0, 0, 0},
	0xE6: CpuInstriction{InstructionTypeAND, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xE7: CpuInstriction{InstructionTypeRST, 0, 0, 0, 0, 0x20},
	0xE8: CpuInstriction{InstructionTypeADD, AddressModeR_N8, RegisterTypeSP, 0, 0, 0},
	0xE9: CpuInstriction{InstructionTypeJP, AddressModeR, RegisterTypeHL, 0, 0, 0},
	0xEA: CpuInstriction{InstructionTypeLD, AddressModeA8_R, RegisterTypeNone, RegisterTypeA, 0, 0},
	0xEE: CpuInstriction{InstructionTypeXOR, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xEF: CpuInstriction{InstructionTypeRST, 0, 0, 0, 0, 0x28},

	0xF0: CpuInstriction{InstructionTypeLDH, AddressModeR_A8, RegisterTypeA, 0, 0, 0},
	0xF1: CpuInstriction{InstructionTypePOP, AddressModeR, RegisterTypeAF, 0, 0, 0},
	0xF2: CpuInstriction{InstructionTypeLD, AddressModeR_MR, RegisterTypeA, RegisterTypeC, 0, 0},
	0xF3: CpuInstriction{InstructionTypeDI, 0, 0, 0, 0, 0},
	0xF5: CpuInstriction{InstructionTypePUSH, AddressModeR, RegisterTypeAF, 0, 0, 0},
	0xF6: CpuInstriction{InstructionTypeOR, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xF7: CpuInstriction{InstructionTypeRST, 0, 0, 0, 0, 0x30},
	0xF8: CpuInstriction{InstructionTypeLD, AddressModeHL_SPR, RegisterTypeHL, RegisterTypeSP, 0, 0},
	0xF9: CpuInstriction{InstructionTypeLD, AddressModeR_R, RegisterTypeSP, RegisterTypeHL, 0, 0},
	0xFA: CpuInstriction{InstructionTypeLD, AddressModeR_A16, RegisterTypeA, 0, 0, 0},
	0xFB: CpuInstriction{InstructionTypeEI, 0, 0, 0, 0, 0},
	0xFE: CpuInstriction{InstructionTypeCP, AddressModeR_N8, RegisterTypeA, 0, 0, 0},
	0xFF: CpuInstriction{InstructionTypeRST, 0, 0, 0, 0, 0x38},
}

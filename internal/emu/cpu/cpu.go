package cpu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/indeedhat/gb-emulator/internal/emu/context"
)

type Cpu struct {
	registers *cpuRegisters

	halted bool

	// interrupts
	ime               bool
	enablingIME       bool
	interruptFlags    uint8
	interruptRegister uint8

	ctx *context.Context
}

func New(ctx *context.Context) {
	ctx.Cpu = &Cpu{
		registers: &cpuRegisters{
			PC: 0x100,
			SP: 0xFFFE,
			A:  0x01,
			F:  0xB0,
			B:  0x00,
			C:  0x13,
			D:  0x00,
			E:  0xD8,
			H:  0x01,
			L:  0x4D,
		},
		ctx: ctx,
	}
}

func (c *Cpu) LoadState(data []byte) {
	r := bytes.NewReader(data)

	binary.Read(r, binary.BigEndian, &c.halted)
	binary.Read(r, binary.BigEndian, &c.ime)
	binary.Read(r, binary.BigEndian, &c.enablingIME)
	binary.Read(r, binary.BigEndian, &c.interruptFlags)
	binary.Read(r, binary.BigEndian, &c.interruptRegister)

	binary.Read(r, binary.BigEndian, &c.registers.A)
	binary.Read(r, binary.BigEndian, &c.registers.F)
	binary.Read(r, binary.BigEndian, &c.registers.B)
	binary.Read(r, binary.BigEndian, &c.registers.C)
	binary.Read(r, binary.BigEndian, &c.registers.D)
	binary.Read(r, binary.BigEndian, &c.registers.E)
	binary.Read(r, binary.BigEndian, &c.registers.H)
	binary.Read(r, binary.BigEndian, &c.registers.L)
	binary.Read(r, binary.BigEndian, &c.registers.SP)
	binary.Read(r, binary.BigEndian, &c.registers.PC)
}

func (c *Cpu) SaveState() []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, c.halted)
	binary.Write(&buf, binary.BigEndian, c.ime)
	binary.Write(&buf, binary.BigEndian, c.enablingIME)
	binary.Write(&buf, binary.BigEndian, c.interruptFlags)
	binary.Write(&buf, binary.BigEndian, c.interruptRegister)

	binary.Write(&buf, binary.BigEndian, c.registers.A)
	binary.Write(&buf, binary.BigEndian, c.registers.F)
	binary.Write(&buf, binary.BigEndian, c.registers.B)
	binary.Write(&buf, binary.BigEndian, c.registers.C)
	binary.Write(&buf, binary.BigEndian, c.registers.D)
	binary.Write(&buf, binary.BigEndian, c.registers.E)
	binary.Write(&buf, binary.BigEndian, c.registers.H)
	binary.Write(&buf, binary.BigEndian, c.registers.L)
	binary.Write(&buf, binary.BigEndian, c.registers.SP)
	binary.Write(&buf, binary.BigEndian, c.registers.PC)

	return buf.Bytes()
}

func (c *Cpu) String(pc uint16) string {
	zf, hf, nf, cf := "-", "-", "-", "-"
	if c.registers.GetFlag(CpuFlagZ) == 1 {
		zf = "Z"
	}
	if c.registers.GetFlag(CpuFlagN) == 1 {
		nf = "N"
	}
	if c.registers.GetFlag(CpuFlagH) == 1 {
		hf = "H"
	}
	if c.registers.GetFlag(CpuFlagC) == 1 {
		cf = "C"
	}

	opcode := CpuInstructions[c.ctx.Bus.Read(pc)]

	return fmt.Sprintf("%08X - %04X: %-7s (%02X %02X %02X) A: %02X F: %s%s%s%s BC: %02X%02X DE: %02X%02X HL: %02X%02X SP: %04X LY: %02X\n",
		c.ctx.Ticks(),
		pc,
		opcode.Type,
		c.ctx.Bus.Read(pc),
		c.ctx.Bus.Read(pc+1),
		c.ctx.Bus.Read(pc+2),
		c.registers.A,
		zf,
		nf,
		hf,
		cf,
		c.registers.B,
		c.registers.C,
		c.registers.D,
		c.registers.E,
		c.registers.H,
		c.registers.L,
		c.registers.SP,
		c.ctx.Lcd.Ly(),
	)
}

func (c *Cpu) Step() error {
	if c.halted {
		c.ctx.EmuCycle(1)
		if c.interruptFlags != 0 {
			c.halted = false
		}
	} else {
		pc := c.registers.PC
		_, instruction := c.fetchIsntruction()
		data, destAddress := c.fetchData(instruction)

		if c.ctx.Debug.Enabled() {
			log.Print(c.ctx.Lcd.String(pc))
			c.ctx.Debug.Update()
			c.ctx.Debug.Print()
		}

		if c.executeInstruction(instruction, data, destAddress) {
			c.ctx.EmuCycle(instruction.CyclesTaken)
		} else {
			c.ctx.EmuCycle(instruction.CyclesUntaken)
		}
	}

	if c.ime {
		c.handleInterrupts()
		c.enablingIME = false
	}

	if c.enablingIME {
		c.ime = true
	}

	return nil
}

func (c *Cpu) stackPop() uint16 {
	val := c.ctx.Bus.Read16(c.registers.SP)
	c.registers.SP += 2

	return val
}

func (c *Cpu) stackPush(value uint16) {
	c.registers.SP--
	c.ctx.Bus.Write(c.registers.SP, uint8(value>>8))

	c.registers.SP--
	c.ctx.Bus.Write(c.registers.SP, uint8(value))
}

func (c *Cpu) fetchIsntruction() (uint8, CpuInstriction) {
	opcode := c.ctx.Bus.Read(c.registers.PC)
	c.registers.PC++

	return opcode, CpuInstructions[opcode]
}

type CpuDestAddress struct {
	Address uint16
}

func (c *Cpu) fetchData(instruction CpuInstriction) (data uint16, destAddr *CpuDestAddress) {
	switch instruction.AddressMode {
	case AddressModeR:
		data = c.readFromRegister(instruction.Register1)

	case AddressModeR_R:
		data = c.readFromRegister(instruction.Register2)

	case AddressModeR_N16,
		AddressModeN16:

		data = c.ctx.Bus.Read16(c.registers.PC)
		c.registers.PC += 2

	case AddressModeHL_SPR,
		AddressModeR_A8,
		AddressModeR_N8,
		AddressModeN8:

		data = uint16(c.ctx.Bus.Read(c.registers.PC))
		c.registers.PC++

	case AddressModeMR_R:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = c.readFromRegister(instruction.Register2)
		if instruction.Register1 == RegisterTypeC {
			destAddr.Address |= 0xFF00
		}

	case AddressModeR_MR:
		address := c.readFromRegister(instruction.Register2)
		if instruction.Register2 == RegisterTypeC {
			address |= 0xFF00
		}
		data = uint16(c.ctx.Bus.Read(address))

	case AddressModeA8_R:
		destAddr = &CpuDestAddress{uint16(c.ctx.Bus.Read(c.registers.PC)) | 0xFF00}
		c.registers.PC++

	case AddressModeMR:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = uint16(c.ctx.Bus.Read(c.readFromRegister(instruction.Register1)))

	case AddressModeMR_N8:
		destAddr = &CpuDestAddress{c.readFromRegister(instruction.Register1)}
		data = uint16(c.ctx.Bus.Read(c.registers.PC))
		c.registers.PC++

	case AddressModeR_HLI:
		hl := c.readFromRegister(RegisterTypeHL)
		data = uint16(c.ctx.Bus.Read(hl))
		c.writeToRegister(RegisterTypeHL, hl+1)

	case AddressModeR_HLD:
		hl := c.readFromRegister(RegisterTypeHL)
		data = uint16(c.ctx.Bus.Read(hl))
		c.writeToRegister(RegisterTypeHL, hl-1)

	case AddressModeHLI_R:
		hl := c.readFromRegister(RegisterTypeHL)
		destAddr = &CpuDestAddress{hl}
		c.writeToRegister(RegisterTypeHL, hl+1)
		data = c.readFromRegister(instruction.Register2)

	case AddressModeHLD_R:
		hl := c.readFromRegister(RegisterTypeHL)
		destAddr = &CpuDestAddress{hl}
		c.writeToRegister(RegisterTypeHL, hl-1)
		data = c.readFromRegister(instruction.Register2)

	case AddressModeA16_R:
		destAddr = &CpuDestAddress{
			c.ctx.Bus.Read16(c.registers.PC),
		}
		c.registers.PC += 2
		data = c.readFromRegister(instruction.Register2)

	case AddressModeR_A16:
		addr := c.ctx.Bus.Read16(c.registers.PC)
		c.registers.PC += 2
		data = uint16(c.ctx.Bus.Read(addr))
	}

	return data, destAddr
}

func (c *Cpu) executeInstruction(instruction CpuInstriction, data uint16, destAddress *CpuDestAddress) bool {
	switch instruction.Type {
	case InstructionTypeADC:
		return c.execADC(instruction, data)
	case InstructionTypeADD:
		return c.execADD(instruction, data)
	case InstructionTypeAND:
		return c.execAND(instruction, data)
	case InstructionTypeCALL:
		return c.execCALL(instruction, data)
	case InstructionTypeCCF:
		return c.execCCF()
	case InstructionTypeCP:
		return c.execCP(instruction, data)
	case InstructionTypeCPL:
		return c.execCPL()
	case InstructionTypeDAA:
		return c.execDAA()
	case InstructionTypeDEC:
		return c.execDEC(instruction, data, destAddress)
	case InstructionTypeDI:
		return c.execDI()
	case InstructionTypeEI:
		return c.execEI()
	case InstructionTypeHALT:
		return c.execHALT()
	case InstructionTypeINC:
		return c.execINC(instruction, data, destAddress)
	case InstructionTypeJP:
		return c.execJP(instruction, data)
	case InstructionTypeJR:
		return c.execJR(instruction, data)
	case InstructionTypeLD:
		return c.execLD(instruction, data, destAddress)
	case InstructionTypeLDH:
		return c.execLDH(instruction, data, destAddress)
	case InstructionTypeOR:
		return c.execOR(instruction, data)
	case InstructionTypePOP:
		return c.execPOP(instruction)
	case InstructionTypePUSH:
		return c.execPUSH(instruction)
	case InstructionTypeRET:
		return c.execRET(instruction)
	case InstructionTypeRETI:
		return c.execRETI(instruction)
	case InstructionTypeRLA:
		return c.execRLA()
	case InstructionTypeRLCA:
		return c.execRLCA()
	case InstructionTypeRRA:
		return c.execRRA()
	case InstructionTypeRRCA:
		return c.execRRCA()
	case InstructionTypeRST:
		return c.execRST(instruction)
	case InstructionTypeSBC:
		return c.execSBC(instruction, data)
	case InstructionTypeSCF:
		return c.execSCF()
	case InstructionTypeSTOP:
		return c.execSTOP(data)
	case InstructionTypeSUB:
		return c.execSUB(instruction, data)
	case InstructionTypeXOR:
		return c.execXOR(data)
	case InstructionTypeCB:
		return c.execCB(instruction, data)
	case InstructionTypeNOP:
		return true

	default:
		panic("instruction not implemented")
	}
}

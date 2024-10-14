package cart

import (
	"errors"
	"io/fs"
	"os"

	. "github.com/indeedhat/gb-emulator/internal/emu/types"
)

type MBC1 struct {
	path string

	romBanks uint16
	romData  []byte

	ramBanks   uint16
	ramData    []byte
	ramEnabled bool

	// registers
	romBank  uint8
	romBank2 uint8
	ramBank  uint8
	mode     uint8

	hasBattery bool
}

func NewMBC1(path string, data []byte, header *CartHeader) (*MBC1, error) {
	m := &MBC1{
		path:       path,
		romBanks:   header.RomBanks(),
		romData:    data,
		ramBanks:   header.RamBanks(),
		ramData:    make([]byte, 0x2000*uint32(header.RamBanks())),
		hasBattery: CartTypeMbc1RamBattery == header.CartType,
	}

	if err := m.Load(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *MBC1) Read(address uint16) byte {
	var offset uint32
	switch true {
	case address < 0x4000:
		if m.mode == 1 && m.romBanks >= 32 {
			offset = uint32(m.romBank2) * 0x4000
		}

		return m.romData[offset+uint32(address)]

	case address < 0x8000:
		if m.romBank == 0 {
			m.romBank = 1
		}
		offset = uint32(m.romBank) * 0x4000
		return m.romData[offset+uint32(address-0x4000)]

	case address <= 0xC000:
		if !m.ramEnabled {
			return 0xFF
		}

		offset = uint32(m.ramBank) * 0x2000
		return m.ramData[offset+uint32(address-0xA000)]
	}

	panic("bad cart read")
}

func (m *MBC1) Write(address uint16, value byte) {
	switch true {
	case address < 0x2000:
		if m.ramBanks > 0 {
			m.ramEnabled = value&0xF == 0xA
		}

	case address < 0x4000:
		value &= 0b11111
		if value&0b11111 == 0x0 {
			value |= 0x1
		}

		m.romBank = value

	case address <= 0x6000:
		if m.ramBanks > 0 {
			m.ramBank = (value & 0b11)
		} else {
			m.romBank2 = value & 0b11
		}

	case address <= 0x8000:
		m.mode = value & 0x1

	case address <= 0xC000:
		if m.ramBanks == 0 {
			return
		}

		offset := uint32(m.ramBank) * 0x2000
		m.ramData[offset+uint32(address-0xA000)] = value
	}
}

// Load implements MBC.
func (m *MBC1) Load() error {
	if !m.hasBattery {
		return nil
	}

	data, err := os.ReadFile(m.path + ".gbsav")
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if data != nil {
		m.ramData = data
	}

	return nil
}

// Save implements MBC.
func (m *MBC1) Save() error {
	if !m.hasBattery {
		return nil
	}

	return os.WriteFile(m.path+".gbsav", m.ramData, 0644)
}

var _ MBC = (*MBC1)(nil)

package emu

import (
	"errors"
	"io/fs"
	"os"
)

type MBC3 struct {
	path string

	romBanks uint16
	romData  []byte

	ramBanks      uint16
	ramData       []byte
	ramRtcEnabled bool

	rtcData        []byte
	rtcLatchedData []byte
	rtcLatched     bool

	// registers
	romBank     uint8
	ramBank     uint8
	rtcRegister uint8
	mode        uint8

	hasBattery bool
}

func NewMBC3(path string, data []byte, header *CartHeader) (*MBC3, error) {
	m := &MBC3{
		path:           path,
		romBanks:       header.RomBanks(),
		romData:        data,
		ramBanks:       header.RamBanks(),
		ramData:        make([]byte, 0x2000*uint32(header.RamBanks())),
		rtcData:        make([]byte, 5),
		rtcLatchedData: make([]byte, 5),
		hasBattery:     CartTypeMbc1RamBattery == header.CartType,
	}

	if err := m.Load(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *MBC3) Read(address uint16) byte {
	var offset uint32
	switch true {
	case address < 0x4000:
		return m.romData[address]

	case address < 0x8000:
		offset = uint32(m.romBank) * 0x4000
		return m.romData[offset+uint32(address-0x4000)]

	case address <= 0xC000:
		if !m.ramRtcEnabled {
			return 0xFF
		}

		if m.ramBank <= 0x3 {
			offset = uint32(m.ramBank) * 0x2000
			return m.ramData[offset+uint32(address-0xA000)]
		}

		if m.rtcLatched {
			return m.rtcLatchedData[m.rtcRegister]
		}

		return m.rtcData[m.rtcRegister]
	}

	panic("bad cart read")
}

func (m *MBC3) Write(address uint16, value byte) {
	switch true {
	case address < 0x2000:
		m.ramRtcEnabled = value == 0x0A
		if m.ramBanks > 0 {
		}

	case address < 0x4000:
		address &= 0b01111111
		m.romBank = value

	case address <= 0x6000:
		if !m.ramRtcEnabled {
			return
		}

		if value <= 0x3 {
			m.ramBank = value
		} else if value >= 0x8 && value <= 0xC {
			m.rtcRegister = value - 0x8
		}

	case address <= 0x8000:
		if value == 0x1 && !m.rtcLatched {
			m.rtcLatched = true
			copy(m.rtcLatchedData, m.rtcData)
		} else if value == 0x0 {
			m.rtcLatched = false
		}

	case address <= 0xC000:
		// TODO
		if m.ramBanks == 0 {
			return
		}

		offset := uint32(m.ramBank) * 0x2000
		m.ramData[offset+uint32(address-0xA000)] = value
	}
}

// Load implements MBC.
func (m *MBC3) Load() error {
	if !m.hasBattery {
		return nil
	}

	data, err := os.ReadFile(m.path + ".gbsav")
	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if data != nil {
		m.rtcData = data[:5]
		m.ramData = data[5:]
	}

	return nil
}

// Save implements MBC.
func (m *MBC3) Save() error {
	if !m.hasBattery {
		return nil
	}

	return os.WriteFile(m.path+".gbsav", append(m.rtcData, m.ramData...), 0644)
}

func (m *MBC3) RtcTick() {
	m.rtcData[0]++
	if m.rtcData[0] >= 60 {
		m.rtcData[0] = 0
		m.rtcData[1]++
	}

	if m.rtcData[1] >= 60 {
		m.rtcData[1] = 0
		m.rtcData[2]++
	}

	if m.rtcData[2] >= 24 {
		m.rtcData[2] = 0
		m.rtcData[3]++
	}
}

var _ MBC = (*MBC3)(nil)

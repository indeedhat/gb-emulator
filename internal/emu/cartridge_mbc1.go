package emu

type MBC1 struct {
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

	variant uint8
}

func NewMBC1(data []byte, romBanks, ramBanks uint16) *MBC1 {
	return &MBC1{
		romBanks: romBanks,
		romData:  data,
		ramBanks: ramBanks,
		ramData:  make([]byte, 0x2000*uint32(ramBanks)),
	}
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
		address &= 0b11111
		if address&0b11111 == 0x0 {
			address |= 0x1
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

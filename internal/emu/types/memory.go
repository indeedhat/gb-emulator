package types

type RamBank struct {
	offset uint16
	data   []byte
}

func NewRamBank(offset uint16, size uint) *RamBank {
	return &RamBank{
		offset: offset,
		data:   make([]byte, size),
	}
}

func (r *RamBank) Read(address uint16) uint8 {
	return r.data[address-r.offset]
}

func (r *RamBank) Write(address uint16, value uint8) {
	r.data[address-r.offset] = value
}

func (r *RamBank) Offset() uint16 {
	return r.offset
}

func (r *RamBank) Bytes() []byte {
	return r.data
}

func (r *RamBank) Fill(data []byte) {
	r.data = data
}

package ppu

import (
	"sort"
)

type OamFlag uint8

const (
	OamFlagBank = 1 << (iota + 3)
	OamFlagDmgPalette
	OamFlagXFlip
	OamFlagYflip
	OamFlagPriority
)

type OamEntry struct {
	x       uint8
	y       uint8
	tileIdx uint8
	// 7   Priority:    0 = No, 1 = BG and Window colors 1–3 are drawn over this OBJ
	// 6   Y flip:      0 = Normal, 1 = Entire OBJ is vertically mirrored
	// 5   X flip:      0 = Normal, 1 = Entire OBJ is horizontally mirrored
	// 4   DMG palette: 0 = OBP0, 1 = OBP1
	// 3   Bank:        0 = Fetch tile from VRAM bank 0, 1 = Fetch tile from VRAM bank 1
	// 0-2 CGB palette: Which of OBP0–7 to use
	flags uint8
}

func (e OamEntry) Check(flag OamFlag) bool {
	return e.flags&uint8(flag) == uint8(flag)
}

type OamRam struct {
	data []byte
}

func (o *OamRam) Bytes() []byte {
	return o.data
}

func (o *OamRam) Fill(data []byte) {
	o.data = data
}

func (o *OamRam) Read(address uint16) uint8 {
	if address >= 0xFE00 {
		address -= 0xFE00
	}

	return o.data[address]
}

func (o *OamRam) Write(address uint16, value uint8) {
	if address >= 0xFE00 {
		address -= 0xFE00
	}

	o.data[address] = value
}

func (o *OamRam) SelectObjects(y uint8, doubleHight bool) []OamEntry {
	var (
		height   uint8 = 8
		selected       = make([]OamEntry, 0, 10)
	)
	if doubleHight {
		height = 16
	}

	for _, entry := range o.iterate() {
		if len(selected) == 10 {
			break
		}

		if entry.x == 0 {
			continue
		}

		if entry.y > y+16 || entry.y+height <= y+16 {
			continue
		}

		selected = append(selected, entry)
	}

	sort.SliceStable(selected, func(i, j int) bool {
		return selected[i].x < selected[j].x
	})

	return selected
}

func (o *OamRam) iterate() func(func(int, OamEntry) bool) {
	return func(yield func(int, OamEntry) bool) {
		var i int

		for c := 0; c < len(o.data); c += 4 {
			entry := OamEntry{
				o.data[c+1],
				o.data[c],
				o.data[c+2],
				o.data[c+3],
			}

			if !yield(i, entry) {
				return
			}

			i++
		}
	}
}

package cart

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	. "github.com/indeedhat/gb-emulator/internal/emu/types"
)

type Cartridge struct {
	data   MBC
	header *CartHeader
	path   string
}

func Load(path string) (*Cartridge, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	c := &Cartridge{
		header: &CartHeader{},
		path:   path,
	}

	data, err := io.ReadAll(fh)
	if err != nil {
		return nil, err
	}

	if err = c.header.Parse(data); err != nil {
		return nil, err
	}

	c.initMbc(path, data)

	return c, nil
}

func (c *Cartridge) LoadState(data []byte) {
	c.data.LoadState(data)
}

func (c *Cartridge) SaveState() []byte {
	return c.data.SaveState()
}

func (c *Cartridge) Read(address uint16) byte {
	return c.data.Read(address)
}

func (c *Cartridge) Write(address uint16, value byte) {
	c.data.Write(address, value)
}

func (c *Cartridge) Save() error {
	return c.data.Save()
}

func (c *Cartridge) Load() error {
	return c.data.Load()
}

func (c *Cartridge) Mbc() MBC {
	return c.data
}

func (c *Cartridge) Filepath() string {
	return c.path
}

func (c *Cartridge) initMbc(path string, data []byte) {
	var err error

	switch c.header.CartType {
	case CartTypeRomOnly:
		c.data = MBCNone(data)
	case CartTypeMbc1, CartTypeMbc1Ram, CartTypeMbc1RamBattery:
		c.data, err = NewMBC1(path, data, c.header)
		if err != nil {
			log.Fatalf("failed to init MBC1: %s", err)
		}
	case CartTypeMbc3, CartTypeMbc3Ram, CartTypeMbc3RamBattery, CartTypeMbc3TimerBattery, CartTypeMbc3TimerRamBattery:
		c.data, err = NewMBC3(path, data, c.header)
		if err != nil {
			log.Fatalf("failed to init MBC3: %s", err)
		}
	default:
		spew.Dump(c.header)
		panic("mbc type not implemented")
	}
}

type MBCNone []byte

func (m MBCNone) Read(address uint16) byte {
	return m[address]
}

func (m MBCNone) Write(address uint16, value byte) {
	log.Print("cart.write not implemented")
}

func (m MBCNone) Save() error {
	// NB: no battery backing for this nbc controller
	return nil
}

func (m MBCNone) Load() error {
	// NB: no battery backing for this nbc controller
	return nil
}

func (m MBCNone) LoadState(_ []byte) {
}

func (m MBCNone) SaveState() []byte {
	return nil
}

type CartHeader struct {
	// 0x0100 - 0x0103
	EntryPoint [4]byte
	// 0x0104 - 0x0133
	Logo [48]byte
	// 0x0134 - 0x0143
	Title string
	// 0x0144 - 0x0145
	NewLicensee string
	// 0x0146
	SgbFlag byte
	// 0x0147
	CartType byte
	// 0x0148
	RomSize byte
	// 0x0149
	RamSize byte
	// 0x014A
	RegionCode byte
	// 0x014B
	OldLicense byte
	// 0x014C
	Version byte
	// 0x014D
	HeaderChecksum byte
	// 0x014E - 0x014F
	GlobalChecksum [2]byte
}

func (h *CartHeader) Parse(data []byte) error {
	if len(data) < 0x014F {
		return errors.New("rom payload is not large enough to contain the header")
	}

	copy(h.EntryPoint[:], data[0x0100:0x0103])
	copy(h.Logo[:], data[0x0104:0x0133])
	h.Title = string(data[0x0134:0x0143])
	h.NewLicensee = string(data[0x0144:0x0145])
	h.SgbFlag = data[0x0146]
	h.CartType = data[0x0147]
	h.RomSize = data[0x0148]
	h.RamSize = data[0x0149]
	h.RegionCode = data[0x014A]
	h.OldLicense = data[0x014B]
	h.Version = data[0x014C]
	h.HeaderChecksum = data[0x014D]
	copy(h.GlobalChecksum[:], data[0x014E:0x014F])

	var checksum uint8
	for i := 0x0134; i < 0x014D; i++ {
		checksum = checksum - data[i] - 1
	}

	if checksum != h.HeaderChecksum {
		return errors.New("invalid header checksum")
	}

	return nil
}

func (h *CartHeader) RomBanks() uint16 {
	switch h.RomSize {
	case 0x00:
		return 2
	case 0x01:
		return 4
	case 0x02:
		return 8
	case 0x03:
		return 16
	case 0x04:
		return 32
	case 0x05:
		return 64
	case 0x06:
		return 128
	case 0x07:
		return 256
	case 0x08:
		return 512
	}
	return 0
}

func (h *CartHeader) RamBanks() uint16 {
	switch h.RomSize {
	case 0x02:
		return 1
	case 0x03:
		return 4
	case 0x04:
		return 16
	case 0x05:
		return 8
	}
	return 0
}

const (
	CartTypeRomOnly                    = 0x00
	CartTypeMbc1                       = 0x01
	CartTypeMbc1Ram                    = 0x02
	CartTypeMbc1RamBattery             = 0x03
	CartTypeMbc2                       = 0x05
	CartTypeMbc2RamBattery             = 0x06
	CartTypeRomRam                     = 0x08
	CartTypeRomRamBattery              = 0x09
	CartTypeMmm01                      = 0x0B
	CartTypeMMM01Ram                   = 0x0C
	CartTypeMMM01RamBattery            = 0x0D
	CartTypeMbc3TimerBattery           = 0x0F
	CartTypeMbc3TimerRamBattery        = 0x10
	CartTypeMbc3                       = 0x11
	CartTypeMbc3Ram                    = 0x12
	CartTypeMbc3RamBattery             = 0x13
	CartTypeMbc5                       = 0x19
	CartTypeMbc5Ram                    = 0x1A
	CartTypeMbc5RamBattery             = 0x1B
	CartTypeMbc5Rumble                 = 0x1C
	CartTypeMbc5RumbleRam              = 0x1D
	CartTypeMbc5RumbleRamBattery       = 0x1E
	CartTypeMbc6                       = 0x20
	CartTypeMbc7SensorRumbleRamBattery = 0x22
	CartTypePocketCamera               = 0xFC
	CartTypeBandaiTatma5               = 0xFD
	CartTypeHuC3                       = 0xFE
	CartTypeHuC1RamBattery             = 0xFF
)

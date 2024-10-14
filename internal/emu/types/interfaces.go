package types

type ReadWriter interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
}

type ReadWriter16 interface {
	ReadWriter

	Read16(address uint16) uint16
	Write16(address, value uint16)
}

type SaveLoader interface {
	Save() error
	Load() error
}

type Ticker interface {
	Tick()
}

type Cart interface {
	ReadWriter
	SaveLoader

	Mbc() MBC
}

type MBC interface {
	ReadWriter
	SaveLoader
}

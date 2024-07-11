package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	cart, err := LoadCartridge("roms/Pokemon Red.gb")
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(cart.Header)
}

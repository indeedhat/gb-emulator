package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must pass a path to a .gb rom")
	}

	e := Emulator{}

	log.Print(e.Run(os.Args[1]))
}

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

	fh, err := os.Create("mem-t.log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(fh)
	log.Print(e.Run(os.Args[1]))
}

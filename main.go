package main

import (
	"flag"
	"log"
	"os"
)

var logFile string

func main() {
	flag.StringVar(&logFile, "log", "", "save log to file")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("must pass a path to a .gb rom")
	}

	e := Emulator{}

	if logFile != "" {
		fh, err := os.Create(logFile)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(fh)
	}
	log.Print(e.Run(flag.Arg(0)))
}

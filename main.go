package main

import (
	"flag"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	var (
		logFile   string
		debugMode bool
	)

	flag.StringVar(&logFile, "log", "", "save log to file")
	flag.BoolVar(&debugMode, "debug", false, "Print out debug logs")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("must pass a path to a .gb rom")
	}

	emu, err := NewEmulator(flag.Arg(0), debugMode)
	if err != nil {
		log.Fatal(err)
	}

	if debugMode && logFile != "" {
		fh, err := os.Create(logFile)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(fh)
	}

	go emu.Run()

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GB emulator")
	renderer := NewLcdRenderer(emu.ctx)
	if err := ebiten.RunGame(renderer); err != nil {
		log.Fatal(err)
	}
}

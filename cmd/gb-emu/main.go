package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/indeedhat/gb-emulator/internal/emu"
	"github.com/indeedhat/gb-emulator/internal/ui"
)

func main() {
	var (
		logFile    string
		debugMode  bool
		cpuProfile bool
	)

	flag.StringVar(&logFile, "log", "", "save log to file")
	flag.BoolVar(&debugMode, "debug", false, "Print out debug logs")
	flag.BoolVar(&cpuProfile, "profile-cpu", false, "generate a cpu profile")
	flag.Parse()

	if cpuProfile {
		fh, err := os.Create("cpu.pprof")
		if err != nil {
			log.Fatal("failed to create pprof file: ", err)
		}

		if err := pprof.StartCPUProfile(fh); err != nil {
			log.Fatal("failed to start cpu profiler: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if len(flag.Args()) < 1 {
		log.Fatal("must pass a path to a .gb rom")
	}

	engine, ctx, err := emu.NewEmulator(flag.Arg(0), debugMode)
	if err != nil {
		log.Fatal(err)
	}

	if logFile != "" {
		fh, err := os.Create(logFile)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(fh)
	}

	go engine.Run()

	_, window := ui.NewFyneRenderer(ctx)
	window.ShowAndRun()

	// ebiten.SetWindowSize(640, 480)
	// ebiten.SetWindowTitle("GB emulator")
	// renderer := emu.NewLcdRenderer(ctx)
	// if err := ebiten.RunGame(renderer); err != nil {
	// 	log.Fatal(err)
	// }
}

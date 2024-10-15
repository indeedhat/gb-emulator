package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

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

	if logFile != "" {
		fh, err := os.Create(logFile)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(fh)
	}

	_, window := ui.NewFyneRenderer()
	window.ShowAndRun()
}

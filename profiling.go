package main

import (
	"log"
	"os"
	"runtime/pprof"
)

var (
	fprof *os.File
	mprof *os.File
)

func startCPUProfiling(fn string) {
	fp, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	fprof = fp
	if err := pprof.StartCPUProfile(fp); err != nil {
		log.Fatal(err)
	}
	log.Println("CPU profiling started")
}

func stopCPUProfiling() {
	if fprof != nil {
		pprof.StopCPUProfile()
		fprof.Close()
		log.Println("CPU profiling stopped")
	}
}

func startMemoryProfiling(fn string) {
	fp, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	fprof = fp
	if err := pprof.WriteHeapProfile(fp); err != nil {
		log.Fatal(err)
	}
	log.Println("Memory profiling started")
}

func stopMemoryProfiling() {
	if mprof != nil {
		mprof.Close()
		log.Println("Memory profiling stopped")
	}
}

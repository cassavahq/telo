package main

import (
	"flag"
	"os"
)

var (
	// DevMode setting
	DevMode    bool
	configPath string
	cpuProfile string
	memProfile string
	port       string
	host       string
)

func initConfiguration() error {

	flag.BoolVar(&DevMode, "d", true, "Development Mode")
	flag.StringVar(&port, "port", os.Getenv("PORT"), "Port")
	flag.StringVar(&host, "host", "", "hostname")
	flag.StringVar(&cpuProfile, "cprof", "", "cpu profile file")
	flag.StringVar(&memProfile, "memprof", "", "memory profile file")
	flag.Parse()

	return nil
}

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
	tokenSecret string
)

func initConfiguration() error {

	flag.BoolVar(&DevMode, "d", true, "Development Mode")
	flag.StringVar(&port, "port", strDefault(os.Getenv("PORT"), "5000"), "Port")
	flag.StringVar(&tokenSecret, "token", strDefault(os.Getenv("TOKEN_SECRET"), "97d6867beb4d4726c5958e03e9337b4599ae3e43a49433a07eea569eb473fcfb"), "Token Secret")
	flag.StringVar(&host, "host", "", "hostname")
	flag.StringVar(&cpuProfile, "cprof", "", "cpu profile file")
	flag.StringVar(&memProfile, "memprof", "", "memory profile file")
	flag.Parse()

	return nil
}

func strDefault(input, alternate string) string{
	if input == ""{
		return alternate
	}

	return input
}
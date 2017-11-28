//go:generate yarn build
package main

import (
	"log"
	"net/http"

	_ "github.com/h4ckm03d/telo/statik"
	"github.com/rakyll/statik/fs"
	"github.com/labstack/echo"
)

func main() {
	// load configuration file or default
	if err := initConfiguration(); err != nil {
		panic(err)
	}

	// check if profiling enabled
	if cpuProfile != "" {
		startCPUProfiling(cpuProfile)
		defer stopCPUProfiling()
	}

	if memProfile != "" {
		startMemoryProfiling(memProfile)
		defer stopMemoryProfiling()
	}
	
	e := echo.New()

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	staticHandler := http.FileServer(statikFS)

	e.GET("/", echo.WrapHandler(staticHandler))
	e.GET("/statics/*", echo.WrapHandler(staticHandler))
		
	e.Logger.Fatal(e.Start(":1323"))
}
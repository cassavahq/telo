//go:generate yarn build
package main

import (
	"log"

	"github.com/cassavahq/telo/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
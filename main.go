package main

import (
	"log"

	"github.com/SoenkeD/sc/src/cmd"
)

// set logging flags
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	cmd.Execute()
}

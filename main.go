package main

import (
	"log"

	"github.com/SoenkeD/sc/src/cmd"
)

// will be populated at build time
// set by adding a build flag
// for go install it needs to be passed as an environment variable
// -ldflags "-X main.commitHash=$(git rev-parse --short HEAD)"
var commitHash string

// set logging flags
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	cmd.SetCommitHash(commitHash)
	cmd.Execute()
}

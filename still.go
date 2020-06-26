package main

import (
	"github.com/danielecook/still/src/commands"
)

// Version - set dynamically during build
var Version = "dev"

func main() {
	commands.Run(Version)
}

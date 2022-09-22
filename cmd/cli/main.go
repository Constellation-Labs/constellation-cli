package main

import (
	"constellation/cmd/cli/commands"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	commands.Execute()
}

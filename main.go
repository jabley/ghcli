package main

import (
	"os"

	"github.com/jabley/ghcli/commands"
)

func main() {
	err := commands.CmdRunner.Execute()
	os.Exit(err.ExitCode)
}

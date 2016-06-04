package commands

import (
	"os"

	"github.com/google/go-github/github"
	"github.com/jabley/ghcli/ui"
	"github.com/jabley/ghcli/version"
)

var cmdVersion = &Command{
	Run:   runVersion,
	Usage: "version",
	Long:  "Shows ghcli version",
}

func init() {
	CmdRunner.Use(cmdVersion, "--version")
}

func runVersion(client *github.Client, cmd *Command, args *Args) {
	ui.Println(version.FullVersion())
	os.Exit(0)
}

package utils

import (
	"os"

	"github.com/google/go-github/github"
	"github.com/jabley/ghcli/ui"
)

func Check(err error) {
	if err != nil {
		ui.Errorln(err)
		os.Exit(1)
	}
}

func CheckClient(client *github.Client) {
	if client == nil {
		ui.Errorln("Need to provide OAuth2 token via GH_OAUTH_TOKEN env var")
		os.Exit(1)
	}
}

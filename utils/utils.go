package utils

import (
	"os"

	"github.com/google/go-github/github"
	"github.com/jabley/ghcli/ui"
)

// Check checks if there error is non-nil. If it is, show the error and exit the process with an error code
func Check(err error) {
	if err != nil {
		ui.Errorln(err)
		os.Exit(1)
	}
}

// CheckClient checks if the client is non-nil. If it is, display an error and exit the process with an error code
func CheckClient(client *github.Client) {
	if client == nil {
		ui.Errorln("Need to provide OAuth2 token via GH_OAUTH_TOKEN env var")
		os.Exit(1)
	}
}

package version

import "fmt"

var Version = "0.0.6"

func FullVersion() string {
	return fmt.Sprintf("ghcli %s", Version)
}

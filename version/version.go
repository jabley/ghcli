package version

import "fmt"

var Version = "0.0.5"

func FullVersion() string {
	return fmt.Sprintf("ghcli %s", Version)
}

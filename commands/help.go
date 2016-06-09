package commands

import (
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/jabley/ghcli/ui"
)

var cmdHelp = &Command{
	Run: runHelp,
	Usage: `
help ghcli
help <COMMAND>
`,
	Long: `Show the help page for a command.

`,
}

func init() {
	CmdRunner.Use(cmdHelp, "--help")
}

func runHelp(client *github.Client, helpCmd *Command, args *Args) {
	if args.IsParamsEmpty() {
		printUsage()
		os.Exit(0)
	}

	command := args.FirstParam()

	if c := lookupCmd(command); c != nil {
		ui.Println(c.HelpText())
		os.Exit(0)
	}
}

func lookupCmd(name string) *Command {
	return CmdRunner.Lookup(name)
}

var helpText = `
These commands are provided by ghcli:

    events   Work with events
    members  Work with members
    teams    Work with teams
`

func printUsage() {
	fmt.Print(helpText)
}

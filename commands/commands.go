package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/jabley/ghcli/ui"
	flag "github.com/ogier/pflag"
)

var (
	NameRe          = "[\\w.][\\w.-]*"
	OwnerRe         = "[a-zA-Z0-9][a-zA-Z0-9-]*"
	NameWithOwnerRe = fmt.Sprintf("^(?:%s|%s\\/%s)$", NameRe, OwnerRe, NameRe)

	CmdRunner = NewRunner()

	ErrNoOrganisationSpecified = errors.New("No organisation specified")

	ErrNoUserSpecified = errors.New("No user specified")
)

type Command struct {
	Run  func(client *github.Client, cmd *Command, args *Args)
	Flag flag.FlagSet

	Key   string
	Usage string
	Long  string

	subCommands map[string]*Command
}

func (c *Command) Call(args *Args) (err error) {
	oauthToken := os.Getenv("GH_OAUTH_TOKEN")

	var client *github.Client

	if oauthToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: oauthToken},
		)

		tc := oauth2.NewClient(oauth2.NoContext, ts)

		client = github.NewClient(tc)
	} else {
		client = nil
	}

	runCommand, err := c.lookupSubCommand(args)
	if err != nil {
		ui.Errorln(err)
		return
	}

	err = runCommand.parseArguments(args)
	if err != nil {
		return
	}

	runCommand.Run(client, runCommand, args)

	return
}

func (c *Command) parseArguments(args *Args) (err error) {
	c.Flag.SetInterspersed(true)
	c.Flag.Init(c.Name(), flag.ContinueOnError)
	c.Flag.Usage = func() {
		ui.Errorln("")
		ui.Errorln(c.Synopsis())
	}
	if err = c.Flag.Parse(args.Params); err == nil {
		for _, arg := range args.Params {
			if arg == "--" {
				args.Terminator = true
			}
		}
		args.Params = c.Flag.Args()
	}

	return
}

func (c *Command) Use(subCommand *Command) {
	if c.subCommands == nil {
		c.subCommands = make(map[string]*Command)
	}
	c.subCommands[subCommand.Name()] = subCommand
}

func (c *Command) Synopsis() string {
	lines := []string{}
	usagePrefix := "Usage:"

	for _, line := range strings.Split(c.Usage, "\n") {
		if line != "" {
			usage := fmt.Sprintf("%s ghcli %s", usagePrefix, line)
			usagePrefix = "      "
			lines = append(lines, usage)
		}
	}
	return strings.Join(lines, "\n")
}

func (c *Command) HelpText() string {
	return fmt.Sprintf("%s\n\n%s", c.Synopsis(), strings.Replace(c.Long, "'", "`", -1))
}

func (c *Command) Name() string {
	if c.Key != "" {
		return c.Key
	}
	usageLine := strings.Split(strings.TrimSpace(c.Usage), "\n")[0]
	return strings.Split(usageLine, " ")[0]
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (c *Command) lookupSubCommand(args *Args) (runCommand *Command, err error) {
	if len(c.subCommands) > 0 && args.HasSubcommand() {
		subCommandName := args.FirstParam()
		if subCommand, ok := c.subCommands[subCommandName]; ok {
			runCommand = subCommand
			args.Params = args.Params[1:]
		} else {
			err = fmt.Errorf("error: Unknown subcommand: %s", subCommandName)
		}
	} else {
		runCommand = c
	}

	return
}

func GetOrg(org string) (string, error) {
	if org == "" {
		return "", ErrNoOrganisationSpecified
	}
	return org, nil
}

func GetUser(user string) (string, error) {
	if user == "" {
		return "", ErrNoUserSpecified
	}
	return user, nil
}

func HttpCleanup(resp *github.Response) {
	if resp == nil {
		return
	}

	if resp.Body == nil {
		return
	}

	resp.Body.Close()
}

func ToJSON(v interface{}) (bytes.Buffer, error) {
	var doc bytes.Buffer
	enc := json.NewEncoder(&doc)
	err := enc.Encode(v)

	return doc, err
}

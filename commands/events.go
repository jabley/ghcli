package commands

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/jabley/ghcli/ui"
	"github.com/jabley/ghcli/utils"
)

var (
	cmdEvents = &Command{
		Run:   listEvents,
		Usage: "events -o USER",
		Long: `List the events that the specified user has performed

## Options:
	-o, --org <USER>
        The user name
`,
	}

	flagEventsUser string
)

func init() {
	cmdEvents.Flag.StringVarP(&flagEventsUser, "user", "u", "", "USER")

	CmdRunner.Use(cmdEvents)
}

func listEvents(client *github.Client, cmd *Command, args *Args) {
	user, err := GetUser(flagEventsUser)
	utils.Check(err)

	if args.Noop {
		ui.Println(fmt.Sprintf("Listing events for user %s", user))
	} else {
		utils.CheckClient(client)

		opt := &github.ListOptions{
			PerPage: 40,
		}

		var allEvents []github.Event
		for {
			events, resp, err := client.Activity.ListEventsPerformedByUser(user, false, opt)

			if err != nil {
				ui.Errorln(err)
				return
			}
			allEvents = append(allEvents, events...)
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
			HttpCleanup(resp)
		}

		doc, err := ToJSON(allEvents)
		utils.Check(err)
		ui.Println(doc.String())
	}
}

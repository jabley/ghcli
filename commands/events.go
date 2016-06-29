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
		Usage: "events -o ORGANISATION -u USER",
		Long: `List the events in the organisation that the specified user can view.

The user needs to be the same as the one associated with the GH_OAUTH_TOKEN, otherwise
you will get an error.

## Options:
	-o, --org <ORGANISATION>
        The organisation name

	-u, --user <USER>
        The user name. This needs to be the user associated with the GH_OAUTH_TOKEN.
`,
	}
	flagEventsOrganisation,
	flagEventsUser string
)

func init() {
	cmdEvents.Flag.StringVarP(&flagEventsUser, "user", "u", "", "USER")
	cmdEvents.Flag.StringVarP(&flagEventsOrganisation, "organisation", "o", "", "ORGANISATION")

	CmdRunner.Use(cmdEvents)
}

func listEvents(client *github.Client, cmd *Command, args *Args) {
	user, err := GetUser(flagEventsUser)
	utils.Check(err)

	org, err := GetOrg(flagEventsOrganisation)
	utils.Check(err)

	if args.Noop {
		ui.Println(fmt.Sprintf("Listing events for user %s in organisation %s", user, org))
	} else {
		utils.CheckClient(client)

		opt := &github.ListOptions{
			PerPage: 40,
		}

		allEvents := make([]*github.Event, 0)
		for {
			var (
				events []*github.Event
				resp   *github.Response
				err    error
			)

			events, resp, err = client.Activity.ListUserEventsForOrganization(org, user, opt)

			if err != nil {
				if resp.StatusCode != 422 {
					ui.Errorln(err)
					return
				}

				// 422 "In order to keep the API fast for everyone, pagination is limited for
				// this resource. Check the rel=last link relation in the Link response header
				// to see how far back you can traverse."

				resp.NextPage = resp.LastPage
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

package commands

import (
	"github.com/google/go-github/github"
	"github.com/jabley/ghcli/ui"
	"github.com/jabley/ghcli/utils"
)

var (
	cmdTeams = &Command{
		Run:   listTeams,
		Usage: "teams -o ORGANISATION",
		Long: `List the members of the organisation

## Options:
	-o, --org <ORGANISATION>
        The organisation name

`,
	}

	flagTeamOrganisation string
)

func init() {
	cmdTeams.Flag.StringVarP(&flagTeamOrganisation, "organisation", "o", "", "ORGANISATION")
	CmdRunner.Use(cmdTeams)
}

func listTeams(client *github.Client, cmd *Command, args *Args) {
	org, err := GetOrg(flagTeamOrganisation)
	utils.Check(err)

	if args.Noop {

	} else {
		utils.CheckClient(client)

		opt := &github.ListOptions{
			PerPage: 40,
		}

		allTeams := make([]*github.Team, 0)
		for {
			teams, resp, err := client.Organizations.ListTeams(org, opt)

			if err != nil {
				ui.Errorln(err)
				return
			}
			allTeams = append(allTeams, teams...)
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
			HttpCleanup(resp)
		}

		doc, err := ToJSON(allTeams)
		utils.Check(err)

		ui.Println(doc.String())
	}
}

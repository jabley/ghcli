package commands

import (
	"bytes"
	"encoding/json"

	"github.com/google/go-github/github"
	"github.com/jabley/ghcli/ui"
	"github.com/jabley/ghcli/utils"
)

var (
	cmdTeams = &Command{
		Run:   listTeams,
		Usage: "teams ORGANISATION",
		Long: `List the members of the organisation

## Options:

`,
	}

	flagMemberOrganisation string
)

func init() {
	CmdRunner.Use(cmdTeams)
}

func listTeams(client *github.Client, cmd *Command, args *Args) {
	utils.CheckClient(client)

	org, err := GetOrg(args)
	utils.Check(err)

	opt := &github.ListOptions{
		PerPage: 40,
	}

	var allTeams []github.Team
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

	var doc bytes.Buffer
	enc := json.NewEncoder(&doc)
	err = enc.Encode(allTeams)

	utils.Check(err)

	ui.Println(doc.String())
}

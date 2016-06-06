package commands

import (
	"bytes"
	"encoding/json"

	"github.com/google/go-github/github"
	"github.com/jabley/ghcli/ui"
	"github.com/jabley/ghcli/utils"
)

var (
	cmdMembers = &Command{
		Run:   listMembers,
		Usage: "members ORGANISATION",
		Long: `List the members of the organisation

## Options:

`,
	}
)

func init() {
	CmdRunner.Use(cmdMembers)
}

func listMembers(client *github.Client, cmd *Command, args *Args) {
	utils.CheckClient(client)

	org, err := GetOrg(args)
	utils.Check(err)

	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 40},
	}

	var allUsers []github.User
	for {
		users, resp, err := client.Organizations.ListMembers(org, opt)

		if err != nil {
			ui.Errorln(err)
			return
		}
		allUsers = append(allUsers, users...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
		HttpCleanup(resp)
	}

	var doc bytes.Buffer
	enc := json.NewEncoder(&doc)
	err = enc.Encode(allUsers)

	utils.Check(err)

	ui.Println(doc.String())
}

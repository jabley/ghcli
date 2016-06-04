package commands

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/google/go-github/github"
	"github.com/jabley/ghcli/ui"
	"github.com/jabley/ghcli/utils"
)

var errNoOrganisation = errors.New("No organisation specified")

var (
	cmdMembers = &Command{
		Run:   listMembers,
		Usage: "members -o <ORGANISATION> [-f <FORMAT>]",
		Long: `List the members of the organisation

## Options:
	-f, --format <FORMAT>
		Sets the output format to be format.
		Supported formats are:
		json

	-o, --org <ORGANISATION>
		The organisation name
`,
	}

	flagMemberFormat,
	flagMemberOrganisation string
)

func init() {
	cmdMembers.Flag.StringVarP(&flagMemberOrganisation, "organisation", "o", "", "ORGANISATION")
	cmdMembers.Flag.StringVarP(&flagMemberFormat, "format", "f", "json", "FORMAT")

	CmdRunner.Use(cmdMembers)
}

func listMembers(client *github.Client, cmd *Command, args *Args) {
	utils.CheckClient(client)

	org, err := getOrg()
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
		httpCleanup(resp)
	}

	var doc bytes.Buffer
	enc := json.NewEncoder(&doc)
	err = enc.Encode(allUsers)

	utils.Check(err)

	ui.Println(doc.String())
}

func getOrg() (string, error) {
	if flagMemberOrganisation == "" {
		return "", errNoOrganisation
	}
	return flagMemberOrganisation, nil
}

func httpCleanup(resp *github.Response) {
	if resp == nil {
		return
	}

	if resp.Body == nil {
		return
	}

	resp.Body.Close()
}

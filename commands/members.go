package commands

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/jabley/ghcli/ui"
	"github.com/jabley/ghcli/utils"
)

var (
	cmdMembers = &Command{
		Run:   listMembers,
		Usage: "members -o ORGANISATION",
		Long: `List the members of the organisation

## Options:
	-o, --org <ORGANISATION>
        The organisation name
`,
	}

	cmdAddMember = &Command{
		Key:   "add",
		Run:   addMember,
		Usage: "members add -o ORGANISATION -u USERNAME",
		Long: `Add a member to the organisation

## Options:
	-o, --org <ORGANISATION>
        The organisation name

    -u, --user <USERNAME>
    	The user name

`,
	}

	flagMemberOrganisation,
	flagMemberUser string
)

func init() {
	cmdMembers.Flag.StringVarP(&flagMemberOrganisation, "organisation", "o", "", "ORGANISATION")

	cmdAddMember.Flag.StringVarP(&flagMemberOrganisation, "organisation", "o", "", "ORGANISATION")
	cmdAddMember.Flag.StringVarP(&flagMemberUser, "user", "u", "", "USER")

	cmdMembers.Use(cmdAddMember)
	CmdRunner.Use(cmdMembers)
}

func listMembers(client *github.Client, cmd *Command, args *Args) {
	org, err := GetOrg(flagMemberOrganisation)
	utils.Check(err)

	if args.Noop {
		ui.Println(fmt.Sprintf("Listing members for organisation %s", org))
	} else {
		utils.CheckClient(client)

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

		doc, err := ToJSON(allUsers)
		utils.Check(err)
		ui.Println(doc.String())
	}
}

func addMember(client *github.Client, cmd *Command, args *Args) {
	org, err := GetOrg(flagMemberOrganisation)
	utils.Check(err)

	user, err := getUser(flagMemberUser)
	utils.Check(err)

	if args.Noop {
		ui.Println(fmt.Sprintf("Adding user %s to organisation %s", user, org))
	} else {
		utils.CheckClient(client)

		membership := new(github.Membership)
		membership.Role = github.String("member")

		_, res, err := client.Organizations.EditOrgMembership(user, org, membership)

		defer HttpCleanup(res)
		utils.Check(err)

		ui.Println(fmt.Sprintf("Added user %s to organisation %s", user, org))
	}
}

func getUser(user string) (string, error) {
	if user == "" {
		return "", ErrNoUserSpecified
	}
	return user, nil
}

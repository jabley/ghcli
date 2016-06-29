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

	cmdRemoveMember = &Command{
		Key:   "remove",
		Run:   removeMember,
		Usage: "members remove -o ORGANISATION -u USERNAME",
		Long: `Remove a member from the organisation

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

	cmdRemoveMember.Flag.StringVarP(&flagMemberOrganisation, "organisation", "o", "", "ORGANISATION")
	cmdRemoveMember.Flag.StringVarP(&flagMemberUser, "user", "u", "", "USER")

	cmdMembers.Use(cmdAddMember)
	cmdMembers.Use(cmdRemoveMember)
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

		allUsers := make([]*github.User, 0)
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

	user, err := GetUser(flagMemberUser)
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

func removeMember(client *github.Client, cmd *Command, args *Args) {
	org, err := GetOrg(flagMemberOrganisation)
	utils.Check(err)

	user, err := GetUser(flagMemberUser)
	utils.Check(err)

	if args.Noop {
		ui.Println(fmt.Sprintf("Removing user %s from organisation %s", user, org))
	} else {
		utils.CheckClient(client)

		res, err := client.Organizations.RemoveOrgMembership(user, org)

		defer HttpCleanup(res)
		utils.Check(err)

		ui.Println(fmt.Sprintf("Removed user %s from organisation %s", user, org))
	}
}

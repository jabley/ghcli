package commands

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/jabley/ghcli/ui"
	"github.com/jabley/ghcli/utils"
)

var (
	cmdRepositories = &Command{
		Run:   listRepositories,
		Usage: "repositories -o ORGANISATION",
		Long: `List the repositories of the organisation

## Options:
	-o, --org <ORGANISATION>
        The organisation name
`,
	}

	cmdRepositoryCommits = &Command{
		Key:   "commits",
		Run:   listCommits,
		Usage: "repositories commits -o ORGANISATION -r REPOSITORY",
		Long: `List the commits for a given repository

## Options:
	-o, --org <ORGANISATION>
        The organisation name

    -r, --repository <REPOSITORY>
    	The repository name

    -s, --SHA <SHA>
    	The SHA or branch to start listing commits from.
    	Default: the repository’s default branch (usually master).

`,
	}

	flagRepositoryOrganisation,
	flagRepository,
	flagRepositorySHA string
)

func init() {
	cmdRepositories.Flag.StringVarP(&flagRepositoryOrganisation, "organisation", "o", "", "ORGANISATION")

	cmdRepositoryCommits.Flag.StringVarP(&flagRepositoryOrganisation, "organisation", "o", "", "ORGANISATION")
	cmdRepositoryCommits.Flag.StringVarP(&flagRepository, "repository", "r", "", "REPOSITORY")
	cmdRepositoryCommits.Flag.StringVarP(&flagRepositorySHA, "sha", "s", "", "SHA")

	cmdRepositories.Use(cmdRepositoryCommits)
	CmdRunner.Use(cmdRepositories)
}

func listRepositories(client *github.Client, cmd *Command, args *Args) {
	org, err := GetOrg(flagRepositoryOrganisation)
	utils.Check(err)

	if args.Noop {
		ui.Println(fmt.Sprintf("Listing repositories for organisation %s", org))
	} else {
		utils.CheckClient(client)

		opt := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 40},
		}

		allRepos := make([]*github.Repository, 0)

		for {
			repos, resp, err := client.Repositories.ListByOrg(org, opt)

			if err != nil {
				ui.Errorln(err)
				return
			}
			allRepos = append(allRepos, repos...)

			if resp.NextPage == 0 {
				break
			}

			opt.ListOptions.Page = resp.NextPage
			HttpCleanup(resp)
		}

		doc, err := ToJSON(allRepos)
		utils.Check(err)
		ui.Println(doc.String())
	}
}

func listCommits(client *github.Client, cmd *Command, args *Args) {
	org, err := GetOrg(flagRepositoryOrganisation)
	utils.Check(err)

	repository, err := GetRepository(flagRepository)
	utils.Check(err)

	if args.Noop {
		ui.Println(fmt.Sprintf("Listing commits for repository %s in organisation %s", repository, org))
	} else {
		utils.CheckClient(client)

		opt := &github.CommitsListOptions{
			ListOptions: github.ListOptions{PerPage: 40},
		}

		if flagRepositorySHA != "" {
			opt.SHA = flagRepositorySHA
		}

		allCommits := make([]*github.RepositoryCommit, 0)

		for {
			commits, resp, err := client.Repositories.ListCommits(org, repository, opt)

			if err != nil {
				ui.Errorln(err)
				return
			}
			allCommits = append(allCommits, commits...)

			if resp.NextPage == 0 {
				break
			}

			if len(allCommits) > 200 {
				break
			}

			opt.ListOptions.Page = resp.NextPage
			HttpCleanup(resp)
		}

		doc, err := ToJSON(allCommits)
		utils.Check(err)
		ui.Println(doc.String())
	}
}

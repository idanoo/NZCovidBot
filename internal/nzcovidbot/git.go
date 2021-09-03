package nzcovidbot

import (
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/waigani/diffparser"
)

// Repo URL
var gitRepo *git.Repository

// Location to store tmp data
var tmpDirectory = "./tmp"

// loadRepo Load repo into gitRepo var
func loadRepo(repository string) {
	r, err := git.PlainOpen(tmpDirectory + "/repo")
	if err != nil {
		if err.Error() == "repository does not exist" {
			r = cloneRepo(repository)
		} else {
			log.Fatal(err)
		}
	}

	// Preload cache data of current rows
	loadRepoIntoCache(tmpDirectory + "/repo")
	gitRepo = r
}

// cloneRepo Clone git repo
func cloneRepo(repository string) *git.Repository {
	if _, err := os.Stat(tmpDirectory); os.IsNotExist(err) {
		err = os.Mkdir(tmpDirectory, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	r, err := git.PlainClone(tmpDirectory+"/repo", false, &git.CloneOptions{
		URL:      repository,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Succesfully cloned repo")
	return r
}

func checkForUpdates() {
	// Fetch updates from remote
	pullOptions := git.PullOptions{
		RemoteName: "origin",
	}

	// Get current commit hash PRE PULL
	currentHead, err := gitRepo.Head()
	if err != nil {
		log.Fatal(err)
	}
	currentHash := currentHead.Hash()

	// Pull latest changes if exist
	worktree, err := gitRepo.Worktree()
	if err != nil {
		log.Fatal(err)
	}
	err = worktree.Pull(&pullOptions)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Println("No updates")
			return
		} else {
			log.Fatal(err)
		}
	}

	// Get current commit hash POST PULL
	head, err := gitRepo.Head()
	if err != nil {
		log.Fatal(err)
	}

	// Get latest commit
	latestCommit, err := gitRepo.CommitObject(head.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// Get current commit
	currentCommit, err := gitRepo.CommitObject(currentHash)
	if err != nil {
		log.Fatal(err)
	}

	// Get patch of changes
	patch, err := currentCommit.Patch(latestCommit)
	if err != nil {
		log.Fatal(err)
	}

	// Parse git diff
	diff, err := diffparser.Parse(patch.String())
	if err != nil {
		log.Fatal(err)
	}

	// Loop through file changes
	for _, file := range diff.Files {
		if strings.HasSuffix(file.NewName, ".csv") {
			if strings.HasPrefix(file.NewName, "locations-of-interest/") {
				for _, hunk := range file.Hunks {
					newRange := hunk.WholeRange
					for _, line := range newRange.Lines {
						if strings.Contains(line.Content, "Start,End,Advice") {
							continue
						}

						if line.Mode == diffparser.ADDED {
							parseCsvRow(line.Content)
						}

					}
				}
			}
		}
	}

	// SEND IT o/---->
	postTheUpdates()
}

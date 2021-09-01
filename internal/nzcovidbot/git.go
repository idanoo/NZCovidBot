package nzcovidbot

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/waigani/diffparser"
)

// Repo URL
var gitRepo *git.Repository

// Current hash we have parsed
var commitHash string

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

	commitHash := getCommitHash()
	log.Printf("Last reported hash: %s", commitHash)
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

	ref, err := r.Head()
	if err != nil {
		log.Fatal(err)
	}
	setCommitHash(ref.Hash().String())

	log.Println("Succesfully cloned repo")
	return r
}

// Set it in memory + write to disk
func setCommitHash(hash string) error {
	log.Printf("Update hash to: %s", hash)
	commitHash = hash
	return os.WriteFile(tmpDirectory+"/commithash", []byte(hash), 0644)
}

// Read from memory, or disk
func getCommitHash() string {
	if commitHash != "" {
		return commitHash
	}

	hash, err := ioutil.ReadFile(tmpDirectory + "/commithash")
	if err != nil {
		log.Fatal(err)
	}

	commitHash = string(hash)
	return commitHash
}

func checkForUpdates() {
	// Fetch updates from remote
	pullOptions := git.PullOptions{
		RemoteName: "origin",
	}

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

	// Get current commit hash
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
	currentCommit, err := gitRepo.CommitObject(plumbing.NewHash(getCommitHash()))
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
			for _, hunk := range file.Hunks {
				newRange := hunk.WholeRange
				for _, line := range newRange.Lines {
					if line.Mode == diffparser.ADDED {
						parseCsvRow("ADDED", line.Content)
					}
					// switch changeType := line.Mode; changeType {
					// case diffparser.UNCHANGED:
					// 	continue
					// case diffparser.ADDED:
					// 	parseCsvRow("ADDED", line.Content)
					// case diffparser.REMOVED:
					// 	continue
					// 	// To re-add in future?
					// 	// parseCsvRow("REMOVED", line.Content)
					// }
				}
			}
		}
	}

	// Store current hash for future comparisons
	setCommitHash(head.Hash().String())

	// SEND IT o/---->
	postTheUpdates()
}

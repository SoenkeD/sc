package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/SoenkeD/sc/src/utils"
	"github.com/spf13/cobra"
)

// will be populated with SetCommitHash()
var commitHash string

func SetCommitHash(hash string) {
	commitHash = hash
}

var update bool

func init() {
	versionCmd.Flags().BoolVar(&update, "update", false, "reinstall the sc binary")
	rootCmd.AddCommand(versionCmd)
}

type Commit struct {
	Sha string `json:"sha"`
}

func getLatestCommitHash(owner, repo, branch string) (string, error) {
	// Build the GitHub API URL for the latest commit on the given branch
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, branch)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch the commit: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get a valid response: %s", resp.Status)
	}

	// Parse the response body
	var commit Commit
	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return commit.Sha, nil
}

func updateBin(latestCommitHash string) {
	log.Println("INFO: There is a newer version of sc available")
	log.Println("INFO: Latest commit", latestCommitHash)
	log.Println("INFO: Current commit", commitHash)

	execBin := "go"
	args := []string{
		"install",
		"github.com/SoenkeD/sc@main",
	}
	env := []string{
		fmt.Sprintf("GOFLAGS=-ldflags=-X=main.commitHash=%s", latestCommitHash),
	}

	execCmd := strings.Join(append(append(env, execBin), args...), " ")
	msg := fmt.Sprintf("Do you want to update by executing '%s' (y/n)", execCmd)
	if confirm, _ := utils.UserConfirm(msg); !confirm {

		return
	}

	cmdOut, err := utils.ExecuteCommand(execBin, args, env, config.RepoRoot)
	if err != nil {
		log.Printf("Executing the command=%s failed with=%s and output %s", execCmd, err, cmdOut)

		return
	}

	log.Println("INFO: Successfully updated")
}

func handleUpdate() {

	if commitHash == "" && !update {
		// no commit was injected at build time
		// skipping the update check
		log.Println("No commit hash was injected at build time")

		return
	}

	latestCommitHash, err := getLatestCommitHash("SoenkeD", "sc", "main")
	if err != nil {
		log.Println("INFO: The update information could not be loaded")

		return
	}

	if latestCommitHash == commitHash {
		log.Println("INFO: The most recent version of sc is already installed", latestCommitHash)

		return
	}

	updateBin(latestCommitHash)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get version and update information",
	RunE: func(cmd *cobra.Command, args []string) error {

		handleUpdate()

		return nil
	},
}

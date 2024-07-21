package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	importCmd.Flags().StringVarP(&config.RepoRoot, "root", "r", "", "base path to generate ctl in")
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))

	importCmd.Flags().StringVarP(&config.Module, "module", "m", "", "name of the module e.g. github.com/SoenkeD/sc")
	viper.BindPFlag("module", rootCmd.PersistentFlags().Lookup("module"))

	importCmd.Flags().StringVarP(&ctlName, "name", "n", "", "name of the ctl")

	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import templates from sources",
	RunE: func(cmd *cobra.Command, args []string) error {

		for _, importItem := range config.Imports {

			err := downloadFolder(
				importItem.RepoOwner,
				importItem.RepoName,
				importItem.RepoPath,
				importItem.LocalPath,
				importItem.Token,
				importItem.RepoPath,
			)

			if err != nil {
				return err
			}
		}

		return nil
	},
}

type GithubFileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

func downloadFile(fileInfo GithubFileInfo, localTargetPath string, token, dirPrefix string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fileInfo.DownloadURL, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	relativeFilePath := strings.TrimPrefix(fileInfo.Path, dirPrefix)

	out, err := os.Create(filepath.Join(localTargetPath, relativeFilePath))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func downloadFolder(repoOwner, repoName, folderPath, localPath, token, dirPrefix string) error {
	files, err := downloadFileInfos(repoOwner, repoName, folderPath, token)
	if err != nil {
		return err
	}

	for _, file := range files {
		relativeFilePath := strings.TrimPrefix(file.Path, dirPrefix)
		localFilePath := filepath.Join(localPath, relativeFilePath)
		if file.Type == "dir" {
			if err := os.MkdirAll(localFilePath, os.ModePerm); err != nil {
				return err
			}
			err = downloadFolder(repoOwner, repoName, file.Path, localPath, token, dirPrefix)
			if err != nil {
				return err
			}
		} else {
			err = downloadFile(file, localPath, token, dirPrefix)
			if err != nil {
				log.Println(err, file.Name, file.DownloadURL)
				return err
			}
		}
	}

	return nil
}

func downloadFileInfos(repoOwner, repoName, folderPath, token string) ([]GithubFileInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, folderPath)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch folder: %s", resp.Status)
	}

	files := []GithubFileInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	return files, nil
}

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

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

func DownloadFolder(repoOwner, repoName, folderPath, localPath, token, dirPrefix, localPathPrefix string) error {
	files, err := downloadFileInfos(repoOwner, repoName, folderPath, token)
	if err != nil {
		return err
	}

	for idx, file := range files {
		relativeFilePath := strings.TrimPrefix(file.Path, dirPrefix)
		localFilePath := filepath.Join(localPath, relativeFilePath)
		if file.Type == "dir" {
			if err := os.MkdirAll(filepath.Join(localPathPrefix, localFilePath), os.ModePerm); err != nil {
				return err
			}
			err = DownloadFolder(repoOwner, repoName, file.Path, localPath, token, dirPrefix, localPathPrefix)
			if err != nil {
				return err
			}
		} else {
			log.Printf("downloading (%d, %d) %s\n", idx+1, len(files), filepath.Join(file.Path, file.Name))
			err = downloadFile(file, filepath.Join(localPathPrefix, localPath), token, dirPrefix)
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

func ExtractGitHubInfo(gitHubURL string) (owner, repo, filePath string, err error) {

	parsedURL, err := url.Parse(gitHubURL)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse URL: %w", err)
	}

	pathSegments := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

	if len(pathSegments) < 4 {
		return "", "", "", fmt.Errorf("invalid GitHub URL format")
	}

	owner = pathSegments[0]
	repo = pathSegments[1]
	filePath = strings.Join(pathSegments[3:], "/")

	return owner, repo, filePath, nil
}

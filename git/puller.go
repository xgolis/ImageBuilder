package git

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Git struct {
	GitRepoPath string `json:"gitRepo,omitempty"`
	GitToken    string `json:"gitToken,omitempty"`
	Username    string `json:"username,omitempty"`
}

func Pull(gitStruct Git) (string, error) {
	_, err := os.Create("repo")
	if err != nil {
		return "", err
	}
	// Clones the repository into the worktree
	_, err = git.PlainClone("./repo", false, &git.CloneOptions{
		URL:  gitStruct.GitRepoPath,
		Auth: &gitHttp.BasicAuth{Username: gitStruct.Username, Password: gitStruct.GitToken},
	})
	if err != nil {
		fmt.Printf("error pulling repo: %v", err)
		return "", err
	}

	file, err := findDockerfile()
	if err != nil {
		return "", err
	}

	return file, nil
}

func findDockerfile() (string, error) {
	var files []string

	filepath.Walk("./repo", func(path string, info fs.FileInfo, err error) error {

		if info.Name() == "Dockerfile" {
			files = append(files, info.Name())
		}
		return nil
	})
	if len(files) == 0 {
		return "", fmt.Errorf("error finding Dockerfile")
	}
	return files[0], nil
}

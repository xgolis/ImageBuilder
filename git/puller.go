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
	_, err := os.ReadDir("./" + gitStruct.Username)
	if err == nil {
		os.RemoveAll("./" + gitStruct.Username)
		fmt.Println("FOUND DIR REMOVING")
	}
	err = os.Mkdir(gitStruct.Username, os.ModePerm)
	if err != nil {
		return "", err
	}
	// Clones the repository into the worktree
	_, err = git.PlainClone("./"+gitStruct.Username, false, &git.CloneOptions{
		URL:  gitStruct.GitRepoPath,
		Auth: &gitHttp.BasicAuth{Username: gitStruct.Username, Password: gitStruct.GitToken},
	})
	if err != nil {
		fmt.Printf("%v", err)
		return "", err
	}

	file, err := findDockerfile("./" + gitStruct.Username)
	if err != nil {
		return "", err
	}

	fmt.Println(file)
	return file, nil
}

func findDockerfile(dir string) (string, error) {
	var files []string = []string{}

	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		fmt.Println(path)
		if info.Name() == "Dockerfile" {
			files = append(files, path)
		}
		return nil
	})
	if len(files) == 0 {
		return "", fmt.Errorf("error finding Dockerfile")
	}
	// fmt.Println(files)
	return files[0], nil
}

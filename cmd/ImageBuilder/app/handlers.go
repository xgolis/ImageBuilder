package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xgolis/ImageBuilder/git"
)

func MakeHandlers() *http.ServeMux {
	mux := *http.NewServeMux()
	mux.HandleFunc("/", pullGit)
	return &mux
	// a.Server.HandleFunc("/", getGit(w http.ResponseWriter, req *http.Request)})
}

//	{
//	    "gitRepo":"localhost",
//	    "gitToken":"localhost",
//	    "username":"aha"
//	}
func pullGit(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(w, "error while reading request: %v", err)
	}

	var gitStruct git.Git
	err = json.Unmarshal(body, &gitStruct)
	if err != nil {
		fmt.Fprintf(w, "unmarshal error: %v", err)
	}

	// fmt.Print(body.GitRepoPath)
	err = git.Pull(gitStruct)
	if err != nil {
		fmt.Fprintf(w, "error while pulling git repository: %v", err)
	}

}

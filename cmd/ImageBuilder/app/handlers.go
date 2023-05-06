package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/xgolis/ImageBuilder/builder"
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
		return
	}

	var gitStruct git.Git
	err = json.Unmarshal(body, &gitStruct)
	if err != nil {
		fmt.Fprintf(w, "unmarshal error: %v", err)
		return
	}

	// fmt.Print(body.GitRepoPath)
	path, err := git.Pull(gitStruct)
	if err != nil {
		fmt.Fprintf(w, "error while pulling git repository: %v", err)
		return
	}

	fmt.Println(gitStruct, path)
	image, err := builder.BuildRepo(path, gitStruct.Username, gitStruct.AppName, gitStruct.Args)
	if err != nil {
		fmt.Fprintf(w, "error building image: %v", err)
		return
	}

	fmt.Printf("Builded image: %s\n", image)

	os.RemoveAll("./" + gitStruct.Username)

	statusOK := map[string]string{"status": "ok"}

	statusJson, err := json.Marshal(statusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(statusJson)
}

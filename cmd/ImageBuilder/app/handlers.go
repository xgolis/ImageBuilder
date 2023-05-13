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

type Response struct {
	Message string `json:"message"`
}

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
func sendError(w *http.ResponseWriter, err error) {
	(*w).Header().Set("Content-Type", "application/json")

	status := Response{
		Message: err.Error(),
	}
	fmt.Print(status.Message)
	statusJson, err := json.Marshal(status)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}

	(*w).Write(statusJson)
}

func pullGit(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		sendError(&w, fmt.Errorf("error while reading request: %v\n", err))
		return
	}

	var gitStruct git.Git
	err = json.Unmarshal(body, &gitStruct)
	if err != nil {
		sendError(&w, fmt.Errorf("unmarshal error: %v\n", err))
		return
	}

	// fmt.Print(body.GitRepoPath)
	path, err := git.Pull(gitStruct)
	if err != nil {
		sendError(&w, fmt.Errorf("error while pulling git repository: %v\n", err))
		return
	}

	fmt.Println(gitStruct, path)
	image, err := builder.BuildRepo(path, gitStruct.AppName, gitStruct.Args)
	if err != nil {
		sendError(&w, fmt.Errorf("error building image: %v\n", err))
		return
	}

	fmt.Printf("Builded image: %s\n", image)

	os.RemoveAll("./" + gitStruct.Username)

	w.Header().Set("Content-Type", "application/json")
	status := Response{
		Message: "Image xgolis/" + gitStruct.AppName + ":latest built",
	}
	statusJson, err := json.Marshal(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(statusJson)
}

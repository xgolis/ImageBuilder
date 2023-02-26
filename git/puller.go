package git

type Git struct {
	GitRepoPath string `json:"gitRepo,omitempty"`
	GitToken    string `json:"gitToken,omitempty"`
	Username    string `json:"username,omitempty"`
}

func Pull(gitStruct Git) error {
	// fmt.Fprintf(w, "haha %s", gitStruct.GitRepoPath)
	return nil
}

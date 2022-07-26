package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v45/github"
)

type repoState struct {
	repoOwner  string
	repoName   string
	currentRef string
}

var state = repoState{
	currentRef: "",
}

func getFileContents(w http.ResponseWriter, r *http.Request) {
	//get the pathInRepo from url
	path := strings.Split(r.URL.Path, "/file-contents/")
	opt := &github.RepositoryContentGetOptions{Ref: state.currentRef}

	client := github.NewClient(nil)

	// get the default branch from current repository and update the ref state
	if state.currentRef == "" {
		r, _, err := client.Repositories.Get(r.Context(), state.repoOwner, state.repoName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if r.DefaultBranch != nil {
			state.currentRef = *r.DefaultBranch
		}
	}

	if state.currentRef == "" {
		http.Error(w, "repository ref does not set", http.StatusBadRequest)
		return
	}

	// get the content file from current repository
	content, _, _, err := client.Repositories.GetContents(r.Context(), state.repoOwner, state.repoName, path[1], opt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if content != nil {
		rawDecodedText, err := base64.StdEncoding.DecodeString(*content.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(rawDecodedText)
	}
}

func checkoutRef(w http.ResponseWriter, r *http.Request) {
	//get the ref from url
	path := strings.Split(r.URL.Path, "/checkout-ref/")
	state.currentRef = path[1]
}

func main() {
	state.repoOwner = os.Getenv("REPO_OWNER")
	state.repoName = os.Getenv("REPO_NAME")

	if state.repoOwner == "" || state.repoName == "" {
		fmt.Printf("REPO_OWNER and REPO_NAME is required!")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/file-contents/", getFileContents)
	mux.HandleFunc("/checkout-ref/", checkoutRef)

	err := http.ListenAndServe(":3334", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

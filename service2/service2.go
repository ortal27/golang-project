package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

type config struct {
	service1Host string
}

type paths struct {
	Paths []string
}

var cnf = config{
	service1Host: "localhost:3334",
}

func getHashFiles(w http.ResponseWriter, r *http.Request) {
	// get the paths from body req
	var body paths
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// to supprt simetric array -> a,b equal to b,a
	sort.Strings(body.Paths)

	hash := ""
	h := sha256.New()

	// implement channels and go routines
	for i := 0; i < len(body.Paths); i++ {
		res, err := http.Get("http://" + cnf.service1Host + "/file-contents/" + body.Paths[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Write(bodyBytes)
		sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
		hash = hash + sha
	}

	// hash of resulting hashes
	hRes := sha256.New()
	hRes.Write([]byte(hash))
	hashRes := base64.URLEncoding.EncodeToString(hRes.Sum(nil))
	w.Write([]byte(hashRes))
}

func main() {
	if os.Getenv("SERVICE1_HOST") != "" {
		cnf.service1Host = os.Getenv("SERVICE1_HOST")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hash-files/", getHashFiles)
	err := http.ListenAndServe(":3333", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

const cacheDir = "./cache"

var objectStore = &DiskObjectStore{root: cacheDir}
var metadataStore *MetadataSqliteStore

func main() {
	var err error
	metadataStore, err = NewMetadataSqliteStore("./metadata.db")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/v8/artifacts/events", ArtifactExistHandler).Methods("POST")
	r.HandleFunc("/v8/artifacts", ArtifactQueryHandler).Methods("POST")

	r.HandleFunc("/v8/artifacts/status", StatusHandler).Methods("GET")
	r.HandleFunc("/v8/artifacts/{hash}", ArtifactFetchHandler).Methods("HEAD")
	r.HandleFunc("/v8/artifacts/{hash}", ArtifactExistHandler).Methods("GET")
	r.HandleFunc("/v8/artifacts/{hash}", ArtifactStoreHandler).Methods("PUT")

	log.Println("Listening on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("StatusHandler")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"enabled"}`))
}

func ArtifactExistHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ArtifactExistHandler")
	vars := mux.Vars(r)
	hash := vars["hash"]

	exists, err := metadataStore.Exists(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ArtifactFetchHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ArtifactFetchHandler")
	vars := mux.Vars(r)
	hash := vars["hash"]

	artifact, err := metadataStore.Get(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if artifact == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	o, err := objectStore.Open(hash)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer o.Close()

	w.Header().Set("Content-Length", strconv.FormatInt(artifact.Size, 10))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("x-artifact-tag", artifact.Tag)
	w.Header().Set("x-artifact-duration", strconv.FormatInt(artifact.Duration, 10))
	w.WriteHeader(http.StatusOK)

	// The artifact is downloaded as an octet-stream. The client should verify the content-length header and response body.
	io.Copy(w, o)
}

type storeResponse struct {
	URLS []string `json:"urls"`
}

func ArtifactStoreHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ArtifactStoreHandler")
	vars := mux.Vars(r)
	hash := vars["hash"]

	duration, err := strconv.ParseInt(r.Header.Get("x-artifact-duration"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	size, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	stored, err := objectStore.Store(hash, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if stored != size {
		// TODO: clean up store
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = metadataStore.Store(hash, r.Header.Get("x-artifact-tag"), duration, stored)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(storeResponse{
		URLS: []string{
			fmt.Sprintf("/v8/artifacts/%s", hash),
		},
	}) // TODO: FQDN
}

type queryRequest struct {
	Hashes []string `json:"hashes"`
}

type queryResponseError struct {
	Message string `json:"message"`
}

type queryResponseObject struct {
	Size     int64              `json:"size,omitempty"`
	Duration int64              `json:"taskDurationMs,omitempty"`
	Tag      string             `json:"tag,omitempty"`
	Error    queryResponseError `json:"error,omitempty"`
}

type queryResponse map[string]queryResponseObject

func ArtifactQueryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ArtifactQueryHandler")
	var req queryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := make(queryResponse)

	var (
		artifact *Artifact
		err      error
	)
	for _, hash := range req.Hashes {
		artifact, err = metadataStore.Get(hash)
		if err != nil {
			// TODO: Should this just push to an error?
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if artifact == nil {
			resp[hash] = queryResponseObject{
				Error: queryResponseError{
					Message: "Artifact not found",
				},
			}
		} else {
			resp[hash] = queryResponseObject{
				Size:     artifact.Size,
				Duration: artifact.Duration,
				Tag:      artifact.Tag,
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

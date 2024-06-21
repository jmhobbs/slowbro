package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/jmhobbs/slowbro/internal/metadata"
	"github.com/jmhobbs/slowbro/internal/object"
)

func ArtifactStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"enabled"}`))
}

func ArtifactEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	/*
		Records an artifacts cache usage event. The body of this request is an array of cache usage events.
		The supported event types are HIT and MISS. The source is either LOCAL the cache event was on the
		users filesystem cache or REMOTE if the cache event is for a remote cache. When the event is a HIT
		the request also accepts a number duration which is the time taken to generate the artifact in the cache.
		[
		{
			"sessionId": "string",
			"source": "LOCAL",
			"event": "HIT",
			"hash": "12HKQaOmR5t5Uy6vdcQsNIiZgHGB",
			"duration": 400
		}
		]
	*/
	w.WriteHeader(http.StatusOK)
}

func ArtifactExists(metadataStore metadata.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.With().Str("url", r.URL.String()).Str("method", r.Method).Logger()

		vars := mux.Vars(r)
		hash := vars["hash"]

		artifact, err := metadataStore.Get(hash)
		if err != nil {
			logger.Error().Err(err).Msg("Error fetching artifact metadata")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if artifact == nil {
			log.Debug().Msg("Artifact metadata not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("x-artifact-tag", artifact.Tag)
		w.Header().Set("x-artifact-duration", strconv.FormatInt(artifact.Duration, 10))
		w.Header().Set("x-ratelimit-remaining", "1000")
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-reset", "0")

		log.Debug().Msg("Artifact exists")
		w.WriteHeader(http.StatusOK)
	}
}

func ArtifactFetch(metadataStore metadata.Store, objectStore object.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.With().Str("url", r.URL.String()).Str("method", r.Method).Logger()

		vars := mux.Vars(r)
		hash := vars["hash"]

		artifact, err := metadataStore.Get(hash)
		if err != nil {
			logger.Error().Err(err).Msg("Error fetching artifact metadata")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if artifact == nil {
			log.Debug().Msg("Artifact metadata not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		o, err := objectStore.Open(hash)
		if err != nil {
			if os.IsNotExist(err) {
				log.Debug().Msg("Object not found")
				w.WriteHeader(http.StatusNotFound)
				return
			}
			logger.Error().Err(err).Msg("Error opening object")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer o.Close()

		w.Header().Set("Content-Length", strconv.FormatInt(artifact.Size, 10))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("x-artifact-tag", artifact.Tag)
		w.Header().Set("x-artifact-duration", strconv.FormatInt(artifact.Duration, 10))
		w.Header().Set("x-ratelimit-remaining", "1000")
		w.Header().Set("x-ratelimit-limit", "1000")
		w.Header().Set("x-ratelimit-reset", "0")
		w.Header().Set("cache-control", "public, max-age=0, must-revalidate")

		w.WriteHeader(http.StatusOK)

		written, err := io.Copy(w, o)
		if err != nil {
			log.Error().Err(err).Msg("Error copying object")
		}
		log.Debug().Int64("sent", written).Int64("size", artifact.Size).Msg("Sent object")
	}
}

type storeResponse struct {
	URLS []string `json:"urls"`
}

func ArtifactStore(metadataStore metadata.Store, objectStore object.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hash := vars["hash"]

		duration, err := strconv.ParseInt(r.Header.Get("x-artifact-duration"), 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing x-artifact-duration")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// TODO: never saw content-length, can't rely on it
		/*
			size, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
			if err != nil {
				log.Printf("Error parsing content-length header: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		*/
		defer r.Body.Close()

		stored, err := objectStore.Store(hash, r.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error storing object")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		/*
			if stored != size {
				log.Printf("Stored size does not match content-length header: %d != %d", stored, size)
				// TODO: clean up store
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		*/

		err = metadataStore.Store(hash, r.Header.Get("x-artifact-tag"), duration, stored)
		if err != nil {
			log.Error().Err(err).Msg("Error storing metadata")
			// TODO: Remove from object store
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)

		err = json.NewEncoder(w).Encode(map[string][]string{
			"urls": {
				fmt.Sprintf("http://%s/v8/artifacts/%s", r.Host, hash),
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Error encoding response")
		}
	}
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

func ArtifactQuery(metadataStore metadata.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req queryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := make(map[string]queryResponseObject, len(req.Hashes))

		var (
			artifact *metadata.Artifact
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Error().Err(err).Msg("Error encoding response")
		}
	}
}

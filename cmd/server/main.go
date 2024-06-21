package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/peterbourgon/ff"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jmhobbs/slowbro/internal/api"
	"github.com/jmhobbs/slowbro/internal/metadata"
	"github.com/jmhobbs/slowbro/internal/object"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	fs := flag.NewFlagSet("slowbro", flag.ContinueOnError)
	var (
		listenAddr   = fs.String("listen", "localhost:8080", "listen address")
		debug        = fs.Bool("debug", false, "enable debug middleware")
		token        = fs.String("token", "867-5309", "API token to accept from turbo")
		databasePath = fs.String("database", "./metadata.db", "sqlite database file")
		cacheDir     = fs.String("cache", "./cache", "cache directory")
		enableLogin  = fs.Bool("login", false, "enable login/link workflow")
		_            = fs.String("config", "", "config file (optional)")
	)

	ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	metadataStore, err := metadata.NewSqliteStore(*databasePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating metadata store")
	}

	objectStore, err := object.NewDiskStore(*cacheDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating object store")
	}

	r := mux.NewRouter()

	ar := r.NewRoute().Subrouter()
	ar.Use(api.AuthMiddleware(*token))

	if *enableLogin {
		r.HandleFunc("/turborepo/token", api.Login(*token))
		r.HandleFunc("/turborepo/success", api.LoginSuccess)

		ar.HandleFunc("/v5/user/tokens/{tokenId}", api.GetUserToken)
		ar.HandleFunc("/v2/user", api.GetUser)
		ar.HandleFunc("/v2/teams", api.GetTeams)
	}

	ar.HandleFunc("/v8/artifacts", api.ArtifactQuery(metadataStore)).Methods("POST")
	ar.HandleFunc("/v8/artifacts/events", api.ArtifactEvents).Methods("POST")
	ar.HandleFunc("/v8/artifacts/status", api.ArtifactStatus).Methods("GET")
	ar.HandleFunc("/v8/artifacts/{hash}", api.ArtifactExists(metadataStore)).Methods("HEAD")
	ar.HandleFunc("/v8/artifacts/{hash}", api.ArtifactFetch(metadataStore, objectStore)).Methods("GET")
	ar.HandleFunc("/v8/artifacts/{hash}", api.ArtifactStore(metadataStore, objectStore)).Methods("PUT")

	if *debug {
		r.Use(loggingMiddleware)
	}

	log.Info().Str("address", *listenAddr).Msg("Starting server")
	if err = http.ListenAndServe(*listenAddr, r); err != nil {
		log.Error().Err(err).Msg("Error running server")
	}
}

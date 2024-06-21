package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type AuthLimitedUser struct {
	Limited       bool    `json:"limited"`
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	Name          *string `json:"name"`
	Username      string  `json:"username"`
	Avatar        *string `json:"avatar"`
	DefaultTeamId *string `json:"defaultTeamId"`
	Version       *string `json:"version"`
}

type AuthUser struct {
	CreatedAt         int64     `json:"createdAt"`
	SoftBlock         *struct{} `json:"softBlock"`
	Billing           *struct{} `json:"billing"`
	ResourceConfig    struct{}  `json:"resourceConfig"`
	StagingPrefix     string    `json:"stagingPrefix"`
	HasTrialAvailable bool      `json:"hasTrialAvailable"`
	RemoteCacheing    struct {
		Enabled bool `json:"enabled"`
	} `json:"remoteCaching"`
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	Name          *string `json:"name"`
	Username      string  `json:"username"`
	Avatar        *string `json:"avatar"`
	DefaultTeamId *string `json:"defaultTeamId"`
	Version       *string `json:"version"`
}

type AuthToken struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Origin string `json:"origin"`
	// Scopes
	ActiveAt  int64 `json:"activeAt"`
	CreatedAt int64 `json:"createdAt"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]AuthUser{
		"user": {
			ID:                "1234567890",
			Email:             "slowbro@example.com",
			Name:              nil,
			Username:          "slowbro",
			Avatar:            nil,
			DefaultTeamId:     nil,
			Version:           nil,
			CreatedAt:         0,
			SoftBlock:         nil,
			Billing:           nil,
			ResourceConfig:    struct{}{},
			StagingPrefix:     "slowbro",
			HasTrialAvailable: false,
			RemoteCacheing: struct {
				Enabled bool `json:"enabled"`
			}{true},
		}})
}
func GetUserToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]AuthToken{
		"token": {
			ID:        "1234567890",
			Name:      "Slowbro",
			Type:      "API",
			Origin:    "CLI",
			ActiveAt:  0,
			CreatedAt: 0,
		}})
	if err != nil {
		log.Printf("error encoding auth token: %v", err)
	}
}

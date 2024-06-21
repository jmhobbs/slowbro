package api

import (
	"encoding/json"
	"net/http"
)

type Membership struct {
	Confirmed   bool   `json:"confirmed"`
	ConfirmedAt int64  `json:"confirmedAt"`
	Role        string `json:"role"`
	TeamID      string `json:"teamId"`
	UID         string `json:"uid"`
	CreatedAt   int64  `json:"createdAt"`
	Created     int64  `json:"created"`
}

type TeamLimited struct {
	Limited    bool       `json:"limited"`
	ID         string     `json:"id"`
	Name       *string    `json:"name"`
	Slug       string     `json:"slug"`
	Avatar     *string    `json:"avatar"`
	Membership Membership `json:"membership"`
	Created    string     `json:"created"`
	CreatedAt  int64      `json:"createdAt"`
}

type Pagination struct {
	Count int    `json:"count"`
	Next  *int64 `json:"next"`
	Prev  *int64 `json:"prev"`
}

func GetTeams(w http.ResponseWriter, r *http.Request) {
	name := "Slowbro"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pagination": Pagination{
			Count: 1,
			Next:  nil,
			Prev:  nil,
		},
		"teams": []TeamLimited{
			{
				Limited: true,
				ID:      "1234567890",
				Name:    &name,
				Slug:    "slowbro",
				Avatar:  nil,
				Membership: Membership{
					Confirmed:   true,
					ConfirmedAt: 0,
					Role:        "MEMBER",
					TeamID:      "1234567890",
					UID:         "1234567890",
					CreatedAt:   0,
					Created:     0,
				},
				Created:   "2021-01-01T00:00:00.000Z",
				CreatedAt: 0,
			},
		},
	})
}

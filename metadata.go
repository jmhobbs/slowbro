package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type MetadataSqliteStore struct {
	db *sql.DB
}

const create string = `
CREATE TABLE IF NOT EXISTS metadata (
	hash TEXT NOT NULL PRIMARY KEY,
	tag TEXT NOT NULL,
	duration INTEGER NOT NULL,
	size INTEGER NOT NULL
) WITHOUT ROWID;
`

func NewMetadataSqliteStore(path string) (*MetadataSqliteStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(create); err != nil {
		return nil, err
	}

	return &MetadataSqliteStore{db: db}, nil
}

func (m *MetadataSqliteStore) Exists(hash string) (bool, error) {
	a, err := m.Get(hash)
	if err != nil {
		return false, err
	}
	return a != nil, nil
}

type Artifact struct {
	Hash     string
	Tag      string
	Duration int64
	Size     int64
}

func (m *MetadataSqliteStore) Get(hash string) (*Artifact, error) {
	var tag string
	var duration int64
	var size int64
	err := m.db.QueryRow("SELECT tag, duration, size FROM metadata WHERE hash = ?", hash).Scan(&tag, &duration, &size)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &Artifact{
		Hash:     hash,
		Tag:      tag,
		Duration: duration,
		Size:     size,
	}, nil
}

func (m *MetadataSqliteStore) Store(hash, tag string, duration, size int64) error {
	_, err := m.db.Exec("INSERT INTO metadata (hash, tag, duration, size) VALUES (?, ?, ?, ?)", hash, tag, duration, size)
	return err
}

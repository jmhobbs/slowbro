package main

import (
	"io"
	"os"
	"path/filepath"
)

type DiskObjectStore struct {
	root string
}

func (d *DiskObjectStore) Open(hash string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(d.root, hash))
}

func (d *DiskObjectStore) Store(hash string, r io.Reader) (int64, error) {
	f, err := os.Create(filepath.Join(d.root, hash))
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return io.Copy(f, r)
}

func (d *DiskObjectStore) Exists(hash string) (bool, error) {
	_, err := os.Stat(filepath.Join(d.root, hash))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

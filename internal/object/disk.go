package object

import (
	"io"
	"os"
	"path/filepath"
)

type DiskStore struct {
	root string
}

func NewDiskStore(root string) (*DiskStore, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}
	return &DiskStore{root: root}, nil
}

func (d *DiskStore) Open(hash string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(d.root, hash))
}

func (d *DiskStore) Store(hash string, r io.Reader) (int64, error) {
	f, err := os.Create(filepath.Join(d.root, hash))
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return io.Copy(f, r)
}

func (d *DiskStore) Exists(hash string) (bool, error) {
	_, err := os.Stat(filepath.Join(d.root, hash))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

package object

import "io"

type Store interface {
	Open(hash string) (io.ReadCloser, error)
	Store(hash string, r io.Reader) (int64, error)
	Exists(hash string) (bool, error)
}

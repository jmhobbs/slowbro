package metadata

type Store interface {
	Get(hash string) (*Artifact, error)
	Store(hash, tag string, duration, size int64) error
}

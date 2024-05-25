package filesystem

type FileSystem interface {
	Get(path string) ([]byte, error)
	Save(path string, content []byte) error
}

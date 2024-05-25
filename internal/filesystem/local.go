package filesystem

import (
	"os"
)

type Local struct {
}

func NewLocal() *Local {
	return &Local{}
}

func (l *Local) Get(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (l *Local) Save(path string, content []byte) error {
	return os.WriteFile(path, content, 0755)
}

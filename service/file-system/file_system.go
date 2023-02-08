package filesystem

import (
	"io"
	"os"
)

type IO struct {
}

func (fsIO IO) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (fsIO IO) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

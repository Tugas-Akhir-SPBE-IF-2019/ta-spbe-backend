package filesystem

import (
	"io"
	"os"
)

type Client interface {
	Create(name string) (*os.File, error)
	Copy(dst io.Writer, src io.Reader) (int64, error)
}

type simpleFSIO struct {
}

func NewSimpleFSIO() (Client, error) {
	return &simpleFSIO{}, nil
}

func (fsIO simpleFSIO) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (fsIO simpleFSIO) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

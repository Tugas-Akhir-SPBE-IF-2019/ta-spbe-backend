package service

import (
	"io"
	"os"
)

type FileSystem interface {
	Create(name string) (*os.File, error)
	Copy(dst io.Writer, src io.Reader) (int64, error)
}

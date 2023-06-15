package model

import "io"

type FileHandle interface {
	Reader() io.ReadCloser
	ContentType() string
	Size() int64
	ETag() string
}

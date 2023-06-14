package renderer

import "io"

type Image interface {
	Reader() io.ReadCloser
	ContentType() string
	Size() int64
	ETag() string
	Bytes() []byte
}

package renderer

import "contrib.rocks/libs/go/model"

type Image interface {
	model.FileHandle
	Bytes() []byte
}

package renderer

import "contrib.rocks/apps/api/go/model"

type Image interface {
	model.FileHandle
	Bytes() []byte
}

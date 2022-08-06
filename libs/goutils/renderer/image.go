package renderer

import "contrib.rocks/libs/goutils/model"

type Image interface {
	model.FileHandle
	Bytes() []byte
}

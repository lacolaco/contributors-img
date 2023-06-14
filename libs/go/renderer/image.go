package renderer

type Image interface {
	ContentType() string
	Bytes() []byte
}

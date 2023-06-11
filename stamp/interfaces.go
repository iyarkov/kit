package stamp

import "errors"

var ErrInvalidBufferSize = errors.New("invalid buffer size")
var ErrInvalidFormat = errors.New("invalid format")

type Stamp interface {
	CompareTo(Stamp) int
}

type Stamper interface {
	Next(Stamp) Stamp
}

type Stamped interface {
	GetStamp() Stamp
	SetStamp(Stamp)
}

type Writer interface {
	Write(Stamp, []byte) error
	BufferSize() int
}

type Reader interface {
	Read(Stamp, []byte) error
}

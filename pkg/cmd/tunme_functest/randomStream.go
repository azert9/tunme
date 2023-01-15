package tunme_functest

import (
	"crypto/sha1"
	"github.com/azert9/tunme/internal/utils"
	"io"
)

type _randomStream struct {
	Buff   []byte
	Offset int
}

func newRandomStream(seed string) io.Reader {

	buff := sha1.Sum([]byte(seed))

	return &_randomStream{
		Buff: buff[:],
	}
}

func (stream *_randomStream) Read(p []byte) (int, error) {

	dst := p

	for len(dst) > 0 {

		if stream.Offset == len(stream.Buff) {
			newBuff := sha1.Sum(stream.Buff)
			stream.Buff = newBuff[:]
			stream.Offset = 0
		}

		copyLen := utils.Min(len(dst), len(stream.Buff)-stream.Offset)

		copy(dst, stream.Buff[stream.Offset:])
		stream.Offset += copyLen
		dst = dst[copyLen:]
	}

	return len(p), nil
}

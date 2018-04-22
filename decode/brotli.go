package decode

import (
	"io"

	"github.com/dsnet/compress/brotli"
)

func Brotli(r io.Reader) io.Reader {
	br, err := brotli.NewReader(r, &brotli.ReaderConfig{})
	if err != nil {
		panic(err)
	}
	r, w := io.Pipe()

	go io.CopyBuffer(w, br, nil)

	return r
}

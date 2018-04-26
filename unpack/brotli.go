package unpack

import (
	"bytes"
	"io"

	"gopkg.in/kothar/brotli-go.v0/enc"
)

func Compress(input io.Reader) io.Reader {
	buffer := bytes.NewBuffer(nil)

	brotliWriter := enc.NewBrotliWriter(nil, buffer)
	defer brotliWriter.Close()

	io.Copy(brotliWriter, input)

	return buffer
}

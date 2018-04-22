package decode

import "io"

type Decoder func(io.Reader) io.Reader

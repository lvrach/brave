package repository

import (
	"io"
)

type Fetcher interface {
	GetResource(resource string) (io.Reader, error)
}

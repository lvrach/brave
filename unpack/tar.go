package unpack

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Tar(dst string, r io.Reader) {

	err := os.MkdirAll(dst, 0755)
	if err != nil {
		panic(err)
	}

	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			panic(err)
		}
		target := filepath.Join(dst, header.Name)
		mode := os.FileMode(header.Mode)

		switch header.Typeflag {

		case tar.TypeDir:
			if err := os.MkdirAll(target, mode); err != nil {
				panic(err)
			}

		// file
		case tar.TypeReg:
			//log.Printf("creating %s, %s\n", header.Name, humanize.Bytes(uint64(header.Size)))
			f, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, mode)
			if err != nil {
				log.Fatal(err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				panic(err)
			}
		}
	}
}

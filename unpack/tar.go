package unpack

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func UnTar(dst string, r io.Reader) {

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
		default:
			panic(fmt.Errorf("not supported type file %s", header.Typeflag))
		}
	}
}

func Tar(workingDir string, sources ...string) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	tw := tar.NewWriter(buf)

	for _, name := range sources {

		file, err := os.Open(filepath.Join(workingDir, name))

		stat, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}

		size := stat.Size()

		hdr := &tar.Header{
			Name:     name,
			Mode:     int64(stat.Mode()),
			Size:     int64(size),
			Typeflag: tar.TypeReg,
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(tw, file)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}

	return buf, nil
}

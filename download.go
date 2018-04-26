package brave

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/lvrach/brave/decode"
	"github.com/lvrach/brave/repository"
	"github.com/lvrach/brave/unpack"
)

type Downloader struct {
	Repositories []repository.Fetcher
	CacheDir     string
	WorkDir      string
}

func (g Downloader) Package(hash string) error {
	start := time.Now()
	pkg, err := g.searchRepositories(hash)

	err = os.MkdirAll(g.CacheDir+"/package", 0755)
	if err != nil {
		return err
	}

	cache, err := os.OpenFile(g.CacheDir+"/package/"+hash, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	checkSum := sha256.New()
	pkgReader, pkgWriter := io.Pipe()

	dst := filepath.Join(g.WorkDir, hash)

	brReader := decode.Brotli(pkgReader)
	go unpack.UnTar(dst, brReader)

	write := io.MultiWriter(cache, checkSum, pkgWriter)

	size, err := io.CopyBuffer(write, pkg, nil)
	if err != nil {
		return err
	}

	log.Printf("checksum: %x\n", checkSum.Sum(nil))
	if fmt.Sprintf("%x", checkSum.Sum(nil)) != hash {
		return fmt.Errorf("checksum miss-match")
	}

	log.Printf("downloaded: %s in %s\n", humanize.Bytes(uint64(size)), time.Since(start))

	return nil
}

func (g Downloader) searchRepositories(hash string) (io.Reader, error) {
	for _, repository := range g.Repositories {
		pkg, err := repository.GetResource("package/" + hash)
		if err != nil {
			log.Println(err)
			continue
		}

		return pkg, nil
	}

	// TODO create an error struct that will include error from all repositories
	return nil, fmt.Errorf("packaged not found in any repository")
}

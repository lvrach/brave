package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/lvrach/brave/repository"

	"github.com/lvrach/brave"
	"github.com/lvrach/brave/repository/bonjour"
	"github.com/lvrach/brave/repository/http"
)

const (
	cacheDir = ".tmp/Library/Brave/cache"
	stageDir = ".tmp/Library/Brave/stage"
)

func main() {

	var useP2P bool
	flag.BoolVar(&useP2P, "p2p", false, "Try to get files from peers on your local network.")
	flag.Parse()

	cmd := flag.Arg(0)

	switch cmd {
	case "install":
		Command{}.Install(flag.Arg(1), useP2P)
	case "share":
		Command{}.Share()
	case "help":
		fmt.Print("Commands:\n\n")
		fmt.Println("install <package_name>  - install the package")
		fmt.Println("share                   - share your cache with others in your network")

		fmt.Println("\nArguments:")
		fmt.Println(" -p2p : try to get files from peers on your local network")
	case "":
		fmt.Println("A command must be specified, try help.")
	}

}

type Command struct{}

func (i Command) Install(name string, useP2P bool) {

	wg := sync.WaitGroup{}

	fetchers := []repository.Fetcher{http.DefaultHTTPMirror}

	if useP2P {
		log.Println("using p2p local discovery")
		wg.Add(1)
		go func() {
			defer wg.Done()
			fetchers = append(bonjour.Discover(), fetchers...)
		}()
	}

	index, err := brave.GetIndex(nil)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	log.Println(fetchers)

	downloader := brave.Downloader{
		fetchers,
		cacheDir,
		stageDir,
	}

	release, ok := index.Find(name)
	if !ok {
		log.Fatalln(name, "not found")
	}
	log.Println("downloading", name, release.PackageHash)

	err = downloader.Package(release.PackageHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
}

func (i Command) Share() {
	srv := bonjour.NewServer(cacheDir)
	err := srv.Serve()
	if err != nil {
		log.Fatal(err)
	}
}

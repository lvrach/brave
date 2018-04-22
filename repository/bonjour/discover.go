package bonjour

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lvrach/brave/repository"
	"github.com/lvrach/brave/repository/http"

	"github.com/grandcat/zeroconf"
)

func Discover() []repository.Fetcher {

	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	fetchers := make([]repository.Fetcher, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println(entry)
			f := http.New(fmt.Sprintf("http://%s:%d/", entry.AddrIPv4[0], entry.Port))
			fetchers = append(fetchers, f)
			cancel()
		}
	}(entries)

	err = resolver.Browse(ctx, "_workstation._tcp", "local.", entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}
	<-ctx.Done()

	return fetchers
}

package bonjour

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grandcat/zeroconf"
)

type Server struct {
	cacheDir string
}

func NewServer(dir string) Server {
	return Server{dir}
}

func (s Server) Serve() error {
	name := "BraveShareLocal"

	port, err := findFreePort(42424)
	if err != nil {
		log.Println(err)
	}

	server, err := zeroconf.Register(name, "_workstation._tcp", "local.", port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		return err
	}
	defer server.Shutdown()

	fs := http.FileServer(http.Dir(s.cacheDir))
	http.Handle("/", fs)

	log.Printf("sharing local cache at 0.0.0.0:%d \n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func findFreePort(fallback int) (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return fallback, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return fallback, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

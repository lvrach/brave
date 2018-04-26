package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/lvrach/brave"
	"gopkg.in/yaml.v2"
)

func main() {

	flag.Parse()
	pd := brave.PackageDefinition{}
	data, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal([]byte(data), &pd)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	p := brave.NewPrepare(pd)
	p.Run()
	defer p.Cleanup()

	dir, err := filepath.Abs(filepath.Join(flag.Arg(1)))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	index := brave.Index{}

	data, err = ioutil.ReadFile(filepath.Join(dir, "index.json"))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = json.Unmarshal([]byte(data), &index)

	index.Put(p.Release())

	data, err = json.Marshal(index)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	ioutil.WriteFile(filepath.Join(dir, "index.json"), data, os.ModePerm)

}

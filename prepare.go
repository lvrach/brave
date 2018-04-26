package brave

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/lvrach/brave/unpack"
)

type PackageDefinition struct {
	Name    string
	Version string
	Prepare []string `yaml:",flow"`
	Pack    string
	Install []string `json:"install"`
}

type Prepare struct {
	PackageDefinition
	dist      string
	dir       string
	hashBuild string
	hashHow   string
	envs      []string
	deferCmd  []*exec.Cmd
}

func NewPrepare(pd PackageDefinition) Prepare {
	return Prepare{
		PackageDefinition: pd,
		dir:               ".",
	}
}

func (p *Prepare) Run() {
	for _, cmd := range p.Prepare {
		p.execCMD(cmd)
	}

	p.packBuild()
}

func (p *Prepare) packBuild() {
	tar, err := unpack.Tar(p.dir, p.Pack)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(filepath.Join(p.dist, "package"), 0755)
	if err != nil {
		log.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile("", "brave_")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	br := unpack.Compress(tar)

	checkSum := sha256.New()
	multi := io.MultiWriter(checkSum, tmpfile)

	size, err := io.Copy(multi, br)
	if err != nil {
		log.Fatal(err)
	}
	p.hashBuild = fmt.Sprintf("%x", checkSum.Sum(nil))

	dist := filepath.Join(p.dist, "package", p.hashBuild)
	os.Rename(tmpfile.Name(), dist)

	log.Printf("package created: %s, build hash: %s\n", humanize.Bytes(uint64(size)), p.hashBuild)
}

func (p *Prepare) Release() Release {
	return Release{
		Name:        p.Name,
		Version:     p.Version,
		PackageHash: p.hashBuild,
		Install:     p.Install,
	}
}

func (p *Prepare) Cleanup() {
	p.execDefer()
}

func (p *Prepare) execDefer() {
	for _, cmd := range p.deferCmd {
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Println(cmd.Dir, cmd)
			log.Fatal(string(out), err)
		}
	}
	p.deferCmd = []*exec.Cmd{}
}

func (p *Prepare) execCMD(rawCMD string) {
	tokens := strings.Split(rawCMD, " ")
	replaceEnvs(tokens)

	switch tokens[0] {
	case "cd":
		log.Printf("cd %s", tokens[1])
		p.dir = tokens[1]
	case "exec":
		cmd := exec.Command(tokens[1], tokens[2:len(tokens)]...)
		cmd.Dir = p.dir
		cmd.Env = os.Environ()
		log.Printf("%s", cmd.Dir)
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Fatal(string(out), err)
		}
	case "defer":
		cmd := exec.Command(tokens[1], tokens[2:len(tokens)]...)
		cmd.Dir = p.dir
		cmd.Env = os.Environ()
		p.deferCmd = append(p.deferCmd, cmd)
	default:
		log.Fatalf("command %s not found\n", rawCMD)
	}
}

func replaceEnvs(tokens []string) {
	for _, t := range tokens {
		os.Environ()
		strings.Replace(t, "$GOPATH", os.Getenv("GOPATH"), 0)
	}
}

package brave

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Installation struct {
	Steps   []string
	WorkDir string
	BinDir  string
}

func (in Installation) Run() {
	for _, step := range in.Steps {
		tokens := strings.Split(step, " ")
		switch tokens[0] {
		case "export":
			from := filepath.Join(in.WorkDir, tokens[1])
			to := filepath.Join(in.BinDir, filepath.Base(tokens[1]))
			err := os.Symlink(from, to)
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("command is not supported")
		}
	}
}

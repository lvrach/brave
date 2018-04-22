package brave

import (
	"encoding/json"

	"github.com/lvrach/brave/repository"

	"github.com/lvrach/brave/repository/http"
)

type Index struct {
	Releases []Release
}

//Release defines a software release
type Release struct {
	Name         string
	Aliases      []string
	Tags         []string
	Version      string
	Description  string
	PackageHash  string `json:"package_hash"`
	InstractHash string
}

func GetIndex(fetcher repository.Fetcher) (Index, error) {
	if fetcher == nil {
		fetcher = http.DefaultHTTPMirror
	}

	index := Index{}
	r, err := fetcher.GetResource("index.json")
	if err != nil {
		return index, err
	}

	err = json.NewDecoder(r).Decode(&index)
	if err != nil {
		return index, err
	}

	return index, nil
}

func (i Index) Find(name string) (Release, bool) {
	for _, r := range i.Releases {
		if r.Name == name {
			return r, true
		}
	}

	return Release{}, false
}

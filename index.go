package brave

import (
	"encoding/json"

	"github.com/lvrach/brave/repository"

	"github.com/lvrach/brave/repository/http"
)

type Index struct {
	Releases []Release `json:"releases"`
}

//Release defines a software release
type Release struct {
	Name         string   `json:"name"`
	Aliases      []string `json:"aliases"`
	Tags         []string `json:"tags"`
	Version      string   `json:"version"`
	Description  string   `json:"description"`
	PackageHash  string   `json:"package_hash"`
	InstructHash string   `json:"instruct_hash"`
	Install      []string `json:"install"`
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
	if ir, exist := i.findIndex(name); exist {
		return i.Releases[ir], true
	}
	return Release{}, false
}

func (i *Index) Put(r Release) {
	ri, exist := i.findIndex(r.Name)
	if !exist {
		i.Releases = append(i.Releases, r)
	} else {
		i.Releases[ri] = r
	}
}

func (i *Index) findIndex(name string) (int, bool) {
	for i, r := range i.Releases {
		if r.Name == name {
			return i, true
		}
	}

	return -1, false
}

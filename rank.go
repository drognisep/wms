package main

import (
	"fmt"
	"github.com/drognisep/wms/data"
	"path/filepath"
	"sort"
)

type rank struct {
	name string
	size data.DiskSpace
}

func (r *rank) String() string {
	return fmt.Sprintf("%s: %s", r.name, r.size)
}

var _ sort.Interface = (ranking)(nil)

type ranking []*rank

func (r ranking) Len() int {
	return len(r)
}

func (r ranking) Less(i, j int) bool {
	return r[i].size < r[j].size
}

func (r ranking) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r ranking) AddFiles(files ...*data.File) ranking {
	if len(files) == 0 {
		return r
	}
	newRanks := make([]*rank, len(files))
	for i, f := range files {
		newRanks[i] = &rank{
			name: f.Name,
			size: f.Size,
		}
	}
	return r.AddRanks(newRanks...)
}

func (r ranking) AddDirectories(dirs ...*data.Directory) ranking {
	if len(dirs) == 0 {
		return r
	}
	newRanks := make([]*rank, len(dirs))
	for i, d := range dirs {
		newRanks[i] = &rank{
			name: "/" + filepath.Base(d.Path),
			size: d.Size(),
		}
	}
	return r.AddRanks(newRanks...)
}

func (r ranking) AddRanks(ranks ...*rank) ranking {
	rlen := len(r)
	newLen := len(ranks)

	newr := make(ranking, rlen+newLen)
	copy(newr, r)
	for i, rank := range ranks {
		newr[i+rlen] = rank
	}
	sort.Sort(sort.Reverse(newr))
	return newr
}

func (r ranking) PrintList() {
	max := 0
	for _, r := range r {
		if len(r.name) > max {
			max = len(r.name)
		}
	}
	format := fmt.Sprintf("%%-%ds : %%s\n", max)
	for _, rank := range r {
		fmt.Printf(format, rank.name, rank.size)
	}
}

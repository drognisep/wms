package data

import (
	"encoding/json"
	"fmt"
)

const (
	KB DiskSpace = 1024
	MB DiskSpace = KB * 1024
	GB DiskSpace = MB * 1024
)

type DiskSpace uint64

func (d DiskSpace) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint64(d))
}

func (d DiskSpace) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &d)
}

func (d DiskSpace) String() string {
	switch {
	case d >= GB:
		return szfmt(d, GB) + "GB"
	case d >= MB:
		return szfmt(d, MB) + "MB"
	case d >= KB:
		return szfmt(d, KB) + "KB"
	default:
		return fmt.Sprintf("%dB", d)
	}
}

func szfmt(val DiskSpace, c DiskSpace) string {
	sz, rem := val/c, val%c
	frem := float64(rem) / float64(c) * 100
	return fmt.Sprintf("%d.%0.0f", sz, frem)
}

type File struct {
	Name string    `json:"name"`
	Path string    `json:"path"`
	Size DiskSpace `json:"size"`
}

type Directory struct {
	Path        string       `json:"path"`
	Files       []*File      `json:"files"`
	Directories []*Directory `json:"directories"`
	cachedSize  DiskSpace
}

func (d *Directory) Size() DiskSpace {
	if d.cachedSize > 0 {
		return d.cachedSize
	}
	var size DiskSpace
	for _, f := range d.Files {
		size += f.Size
	}
	for _, d := range d.Directories {
		size += d.Size()
	}
	d.cachedSize = size
	return size
}

func (d *Directory) AddFile(f *File) {
	d.cachedSize = 0
	d.Files = append(d.Files, f)
}

func (d *Directory) AddDir(dir *Directory) {
	d.cachedSize = 0
	d.Directories = append(d.Directories, dir)
}

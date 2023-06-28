package main

import (
	"errors"
	"fmt"
	"github.com/drognisep/wms/data"
	"io/fs"
	"os"
	"path/filepath"
)

func walkDir(debug bool, directory string) (*data.Directory, error) {
	var err error
	directory, err = filepath.Abs(directory)
	if err != nil {
		return nil, err
	}
	root := &data.Directory{
		Path: directory,
	}
	parentCache := map[string]*data.Directory{
		directory: root,
	}
	err = filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			switch {
			case errors.Is(err, fs.ErrPermission):
				fallthrough
			case errors.Is(err, os.ErrPermission):
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			default:
				return err
			}
		}
		if d.IsDir() && path == directory {
			return nil
		}
		if debug {
			fmt.Println(path)
		}
		pathDir := filepath.Dir(path)
		parent, ok := parentCache[pathDir]
		if !ok {
			parent = &data.Directory{
				Path: pathDir,
			}
			parentCache[pathDir] = parent
		}
		if d.IsDir() {
			dir := &data.Directory{
				Path: path,
			}
			parentCache[path] = dir
			parent.AddDir(dir)
			return nil
		}
		f := &data.File{
			Name: d.Name(),
			Path: path,
		}
		info, err := d.Info()
		if cleanseError(err) != nil {
			return err
		}
		f.Size = data.DiskSpace(info.Size())
		parent.AddFile(f)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return root, nil
}

func cleanseError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, fs.ErrPermission):
		fallthrough
	case errors.Is(err, os.ErrPermission):
		return nil
	default:
		return err
	}
}

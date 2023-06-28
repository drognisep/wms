package data

import (
	"fmt"
	"math"
	"path/filepath"
)

func (d *Directory) MinorOutliers(debug bool) ([]*File, []*Directory) {
	return d.outliers(debug, 0.5)
}

func (d *Directory) Outliers(debug bool) ([]*File, []*Directory) {
	return d.outliers(debug, 1.0)
}

func (d *Directory) LargeOutliers(debug bool) ([]*File, []*Directory) {
	return d.outliers(debug, 2.0)
}

func (d *Directory) outliers(debug bool, sdCoef float64) ([]*File, []*Directory) {
	set := d.entries()
	average, stdDev := d.sizeMetrics(set)
	if average == math.Inf(1) || stdDev == math.Inf(1) {
		panic("Likely floating point overflow")
	}

	var (
		lenFiles  = len(d.Files)
		threshold = DiskSpace(math.Floor(sdCoef*stdDev + average))
		files     []*File
		dirs      []*Directory
	)

	for i := 0; i < len(set.size); i++ {
		if set.size[i] > threshold {
			if i >= lenFiles {
				dir := d.Directories[i-lenFiles]
				if debug {
					fmt.Println("Directory outlier:", dir.Path)
				}
				dirs = append(dirs, dir)
				continue
			}
			if debug {
				fmt.Println("File outlier:", d.Files[i].Path)
			}
			files = append(files, d.Files[i])
		}
	}
	return files, dirs
}

func (d *Directory) sizeMetrics(set entrySet) (average float64, stdDev float64) {
	var (
		total    DiskSpace
		count    DiskSpace
		sumOSq   DiskSpace
		variance DiskSpace
	)

	for _, sz := range set.size {
		total += sz
		count++
	}
	_average := total / count
	average = float64(_average)

	for _, sz := range set.size {
		sz -= _average
		sumOSq += sz * sz
	}
	variance = sumOSq / count
	stdDev = math.Sqrt(float64(variance))
	return average, stdDev
}

type entrySet struct {
	size   []DiskSpace
	string []string
}

func (d *Directory) entries() entrySet {
	lenFiles := len(d.Files)
	l := lenFiles + len(d.Directories)
	set := entrySet{
		size:   make([]DiskSpace, l),
		string: make([]string, l),
	}

	for i, file := range d.Files {
		set.size[i] = file.Size
		set.string[i] = file.Name
	}
	for i, dir := range d.Directories {
		set.size[i+lenFiles] = dir.Size()
		set.string[i+lenFiles] = "/" + filepath.Base(dir.Path)
	}
	return set
}

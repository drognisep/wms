package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirectory_Outliers_Files(t *testing.T) {
	data := &Directory{
		Path: "/",
		Files: []*File{
			{
				Name: "a",
				Path: "/a",
				Size: 8,
			},
			{
				Name: "b",
				Path: "/b",
				Size: 8,
			},
			{
				Name: "c",
				Path: "/c",
				Size: 10,
			},
		},
	}

	files1, dirs1 := data.Outliers()
	assert.Len(t, dirs1, 0)
	assert.Len(t, files1, 1)

	files2, dirs2 := data.LargeOutliers()
	assert.Len(t, dirs2, 0)
	assert.Len(t, files2, 0)
}

func TestDirectory_Outliers_Dirs(t *testing.T) {
	data := &Directory{
		Path: "/",
		Directories: []*Directory{
			{
				Path: "/a",
				Files: []*File{
					{
						Name: "a",
						Path: "/a/a",
						Size: 8,
					},
				},
			},
			{
				Path: "/b",
				Files: []*File{
					{
						Name: "b",
						Path: "/b/b",
						Size: 8,
					},
				},
			},
			{
				Path: "/c",
				Files: []*File{
					{
						Name: "c",
						Path: "/c/c",
						Size: 10,
					},
				},
			},
		},
	}

	files1, dirs1 := data.Outliers()
	assert.Len(t, dirs1, 1)
	assert.Len(t, files1, 0)

	files2, dirs2 := data.LargeOutliers()
	assert.Len(t, dirs2, 0)
	assert.Len(t, files2, 0)
}

//go:build windows

package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSpaceFree(t *testing.T) {
	total, free, err := GetSpaceFree()
	assert.NoError(t, err)
	assert.Greater(t, total, free)
	t.Logf("Total: %s, free: %s\n", total.String(), free.String())
}

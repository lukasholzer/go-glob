package fixtures

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateFixture(t *testing.T) {
	folder := CreateNew(t, map[string]string{
		"a/b/c": "",
		"d":     "",
	})

	require.DirExists(t, filepath.Join(folder, "a"))
	require.DirExists(t, filepath.Join(folder, "a/b"))
	require.FileExists(t, filepath.Join(folder, "a/b/c"))
	require.FileExists(t, filepath.Join(folder, "d"))
}

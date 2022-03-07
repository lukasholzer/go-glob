package glob

import (
	"path/filepath"
	"testing"

	"github.com/lukasholzer/glob/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobPattern(t *testing.T) {
	fileMap := map[string]string{
		"a/b/c/d.txt":           "",
		"c":                     "",
		"some-longer-file":      "",
		"some-longer-file.d.ts": "",
		"d/f/g":                 "",
	}

	dir := fixtures.CreateNew(t, fileMap)

	files, err := Glob(&globOptions{
		Patterns:       []string{"**/*"},
		CWD:            dir,
		IgnorePatterns: []string{"a", "c"},
	})
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{
		"some-longer-file",
		"some-longer-file.d.ts",
		filepath.Join("d/f/g"),
	}, files)
}

func TestGlobPatternAbsolute(t *testing.T) {
	fileMap := map[string]string{
		"a/b/c/d.txt":           "",
		"c":                     "",
		"some-longer-file":      "",
		"some-longer-file.d.ts": "",
		"d/f/g":                 "",
	}

	dir := fixtures.CreateNew(t, fileMap)

	files, err := Glob(&globOptions{
		Patterns:       []string{"**/*"},
		CWD:            dir,
		IgnorePatterns: []string{"a", "c"},
		AbsolutePaths:  true,
	})
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{
		filepath.Join(dir, "some-longer-file"),
		filepath.Join(dir, "some-longer-file.d.ts"),
		filepath.Join(dir, "d/f/g"),
	}, files)
}

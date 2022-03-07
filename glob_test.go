package glob

import (
	"path/filepath"
	"testing"

	"github.com/lukasholzer/go-glob/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobApiShouldThrowWithoutPattern(t *testing.T) {
	files, err := Glob(CWD("test"))
	assert.Error(t, err)
	assert.Len(t, files, 0)
}

func TestGlobApiShouldUseSpreadOptions(t *testing.T) {
	outdir := fixtures.CreateNew(t, map[string]string{
		"a/b/c":   "",
		"d/e":     "",
		"f/g":     "",
		".ignore": "d\n.ignore",
	})
	files, err := Glob(Pattern("**/*"), CWD(outdir), IgnorePattern("a"), IgnoreFile(".ignore"))
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{filepath.Join("f/g")}, files)
}

func TestGlobPattern(t *testing.T) {
	fileMap := map[string]string{
		"a/b/c/d.txt":           "",
		"c":                     "",
		"some-longer-file":      "",
		"some-longer-file.d.ts": "",
		"d/f/g":                 "",
	}

	dir := fixtures.CreateNew(t, fileMap)

	files, err := Glob(&Options{
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

	files, err := Glob(&Options{
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

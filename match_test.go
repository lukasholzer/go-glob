package glob

import (
	"path/filepath"
	"testing"

	"github.com/lukasholzer/go-glob/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestMatchPatternDirect(t *testing.T) {
	outdir := fixtures.CreateNew(t, map[string]string{
		"a/b/c-test.txt": "",
	})

	p, err := Parse("a/b/c-test.txt")
	require.NoError(t, err)

	match, err := matchPattern(outdir, p, nil, false, zap.NewNop())
	require.NoError(t, err)
	assert.Len(t, match, 1)
	assert.Equal(t, filepath.Join("a/b/c-test.txt"), match[0])

	absMatch, err := matchPattern(outdir, p, nil, true, zap.NewNop())
	require.NoError(t, err)
	assert.Len(t, absMatch, 1)
	assert.Equal(t, filepath.Join(outdir, "a/b/c-test.txt"), absMatch[0])
}

func TestStarGlobAndGlobstarMatching(t *testing.T) {
	tests := map[string]interface{}{
		"**/*": map[string]bool{
			"a.txt":     true,
			"a/x/y.txt": true,
			"a/x/y/z":   true,
		},
		"a/**/*.txt": map[string]bool{
			"a.txt":     false,
			"a/x/y.txt": true,
			"a/x/y/z":   false,
		},
		"a/*.txt": map[string]bool{
			"a.txt":     false,
			"a/b.txt":   true,
			"a/x/y.txt": false,
			"a/x/y/z":   false,
		},
		"a*.txt": map[string]bool{
			"a.txt":     true,
			"a/b.txt":   false,
			"a/x/y.txt": false,
			"a/x/y/z":   false,
		},
		"*.txt": map[string]bool{
			"a.txt":     true,
			"a/b.txt":   false,
			"a/x/y.txt": false,
			"a/x/y/z":   false,
		},
		".*": map[string]bool{
			".ignore":            true,
			".git":               true,
			"asdf":               false,
			"other/path/.ignore": false,
		},
		"?at": map[string]bool{
			"cat": true,
			"fat": true,
			"at":  false,
			"ats": false,
		},
	}

	for pattern, tc := range tests {
		p, err := Parse(pattern)
		require.NoError(t, err)
		for path, want := range tc.(map[string]bool) {
			t.Run(pattern, func(t *testing.T) {
				assert.Equal(t, want, Match(p, path), "Pattern %s matching %s path failed with RegExp: %v -> %v", pattern, path, p.String(), want)
			})
		}

	}
}

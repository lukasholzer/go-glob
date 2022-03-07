package glob

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	}

	for pattern, tc := range tests {
		r, err := Parse(pattern)
		require.NoError(t, err)
		for path, want := range tc.(map[string]bool) {
			t.Run(pattern, func(t *testing.T) {
				assert.Equal(t, want, Match(r, path), "Pattern %s matching %s path failed with RegExp: %v -> %v", pattern, path, r.String(), want)
			})
		}

	}
}

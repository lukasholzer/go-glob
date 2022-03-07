package glob

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldParseGlob(t *testing.T) {
	tests := map[string]string{
		"some-file.ext*": `^some-file\.ext([^/]*)$`,
		"**/*":           `^((?:[^/]*(?:\/|$))*)([^/]*)$`,           // https://regex101.com/r/HBDUUM/2
		"src/**/*.txt":   `^src\/((?:[^/]*(?:\/|$))*)([^/]*)\.txt$`, // https://regex101.com/r/fy8Ini/1
	}

	for pattern, want := range tests {
		t.Run(pattern, func(t *testing.T) {
			r, err := Parse(pattern)
			require.NoError(t, err)
			assert.Equal(t, want, r.String(), "wanted %s for glob pattern %s", r.String(), pattern)
		})
	}

}

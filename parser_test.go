package glob

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsedPatternGlobster(t *testing.T) {
	p, err := Parse("src/some/folder/**/*.js")
	require.NoError(t, err)

	assert.Equal(t, "src/some/folder/", p.Base)
	assert.True(t, p.IsGlobstar)
}

func TestShouldParseGlob(t *testing.T) {
	tests := map[string]string{
		"some-file.ext*": `^some-file\.ext([^/]*)$`,
		"**/*":           `^((?:[^/]*(?:\/|$))*)([^/]*)$`,           // https://regex101.com/r/HBDUUM/2
		"src/**/*.txt":   `^src\/((?:[^/]*(?:\/|$))*)([^/]*)\.txt$`, // https://regex101.com/r/fy8Ini/1
		"**/.*":          `^((?:[^/]*(?:\/|$))*)\.([^/]*)$`,         // https://regex101.com/r/3sOACA/1
		"**":             `^((?:[^/]*(?:\/|$))*)$`,                  // https://regex101.com/r/ZHASlR/1
	}

	for pattern, want := range tests {
		t.Run(pattern, func(t *testing.T) {
			p, err := Parse(pattern)
			require.NoError(t, err)
			assert.Equal(t, want, p.String(), "wanted %s for glob pattern %s", p.String(), pattern)
		})
	}
}

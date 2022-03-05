package glob

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldParseARegularGlobsterToRegex(t *testing.T) {
	r, err := Parse("src/**/*.txt")
	require.NoError(t, err)

	assert.Equal(t, `^src\/([^/]*)\/([^/]*)\.txt$`, r.String())
}
func TestShouldParseAWinGlobsterToRegex(t *testing.T) {
	r, err := Parse(`src\**\*.txt`)
	require.NoError(t, err)

	assert.Equal(t, `^src\/([^/]*)\/([^/]*)\.txt$`, r.String())
}

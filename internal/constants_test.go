package internal

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStarGlob(t *testing.T) {
	c, err := regexp.Compile(fmt.Sprintf(`^%s$`, STAR))
	require.NoError(t, err)
	assert.True(t, c.MatchString("abcd"))
	assert.True(t, c.MatchString("abcd.txt"))
	assert.True(t, c.MatchString("abcd.d.ts"))
	assert.True(t, c.MatchString("abcd-some-_other"))
	assert.False(t, c.MatchString("abcd/asdf"))
	assert.False(t, c.MatchString("/abcd"))
}

func TestGlobStarGlob(t *testing.T) {
	c, err := regexp.Compile(fmt.Sprintf(`^%s$`, GLOBSTER))
	require.NoError(t, err)
	assert.True(t, c.MatchString("abcd"))
	assert.True(t, c.MatchString("abcd.txt"))
	assert.True(t, c.MatchString("abcd.d.ts"))
	assert.True(t, c.MatchString("abcd-some-_other"))
	assert.True(t, c.MatchString("abcd/asdf"))
	assert.True(t, c.MatchString("abcd/asdf/segment.txt"))
	assert.True(t, c.MatchString("/abcd"))
	assert.True(t, c.MatchString("/abcd/some/more-segments.txt"))
}

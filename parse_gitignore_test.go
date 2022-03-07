package glob

import (
	"path/filepath"
	"testing"

	"github.com/lukasholzer/go-glob/fixtures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadGitignoreFile(t *testing.T) {

	folder := fixtures.CreateNew(t, map[string]string{
		".gitignore": `# some comment
node_modules

.github/styles`,
	})
	patterns, err := ParseGitignore(filepath.Join(folder, ".gitignore"))
	require.NoError(t, err)

	compareParsedGlob(t, map[string]string{
		"node_modules":   `^node_modules$`,
		".github/styles": `^\.github\/styles$`,
	}, patterns)
}

func TestShouldParseGitignroeSyntax(t *testing.T) {
	patterns, err := ParseGitignoreContent(`

  # some comment
      !ignore
some/path/here
yarn-debug.log*
  `)
	require.NoError(t, err)

	compareParsedGlob(t, map[string]string{
		"!ignore":         `^\!ignore$`,
		"some/path/here":  `^some\/path\/here$`,
		"yarn-debug.log*": `^yarn-debug\.log([^/]*)$`,
	}, patterns)
}

func TestShouldParseVariousLines(t *testing.T) {
	tests := map[string]interface{}{
		"logs":     map[string]string{"logs": "^logs$"},
		"   logs ": map[string]string{"logs": "^logs$"},
		// "!ignore":  map[string]string{"!ignore": `^\!ignore$`}, // wrong should be a not regex
		"\tfolder name": map[string]string{"folder name": "^folder name$"},
		"/*.c":          map[string]string{"/*.c": `^\/([^/]*)\.c$`},
		"# comment":     map[string]string{},
		"":              map[string]string{},
		"  ":            map[string]string{},
		"  \r\n":        map[string]string{},
		"  \n\n\n":      map[string]string{},
	}

	for content, want := range tests {
		t.Run(content, func(t *testing.T) {
			got, err := ParseGitignoreContent(content)
			require.NoError(t, err)
			compareParsedGlob(t, want.(map[string]string), got)
		})
	}
}

func compareParsedGlob(t *testing.T, want map[string]string, got map[string]*ParsedPattern) {
	assert.Len(t, got, len(want))
	for glob, reg := range got {
		w, ok := want[glob]
		require.True(t, ok, "did not receive a pattern %s", glob)
		assert.Equal(t, w, reg.RegExp.String())
	}
}

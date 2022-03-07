package glob

import (
	"os"
	"strings"
)

// ParseGitignoreContent parses a file according to: http://git-scm.com/docs/gitignore
func ParseGitignoreContent(content string) (GlobPatterns, error) {
	pattern := make(GlobPatterns)

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 || strings.HasPrefix(trimmed, "#") {
			continue
		}

		reg, err := Parse(trimmed)
		if err != nil {
			return nil, err
		}
		pattern[trimmed] = reg
	}

	return pattern, nil
}

func ParseGitignore(file string) (GlobPatterns, error) {
	if _, err := os.Stat(file); err != nil {
		return make(GlobPatterns), nil
	}

	bs, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	return ParseGitignoreContent(string(bs))
}

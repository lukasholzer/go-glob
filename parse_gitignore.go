package glob

import (
	"os"
	"regexp"
	"strings"
)

// ParseGitignoreContent parses a file according to: http://git-scm.com/docs/gitignore
func ParseGitignoreContent(content string) (map[string]*regexp.Regexp, error) {
	pattern := make(map[string]*regexp.Regexp)

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

func ParseGitignore(file string) (map[string]*regexp.Regexp, error) {
	if _, err := os.Stat(file); err != nil {
		return make(map[string]*regexp.Regexp), nil
	}

	bs, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	return ParseGitignoreContent(string(bs))
}

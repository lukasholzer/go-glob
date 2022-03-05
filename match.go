package glob

import (
	"regexp"
)

func Match(r *regexp.Regexp, path string) bool {
	return r.MatchString(path)
}

package glob

import (
	"regexp"
	"strings"

	"github.com/lukasholzer/go-glob/internal"
)

func cleanPattern(input string) string {
	// windows path segments to posix characters
	input = strings.ReplaceAll(input, "\\", "/")
	input = strings.Replace(input, "***", "*", 1)
	input = strings.Replace(input, "**/**", "**", 1)
	input = strings.Replace(input, "**/**/**", "**", 1)

	return input
}

type ParsedPattern struct {
	// regExp is the regular expression as string that falls out of the parsedPattern
	stringPattern string
	isBaseSet     bool

	// input is the original glob pattern
	Input      string
	RegExp     *regexp.Regexp
	IsGlobstar bool
	// base is the base folder that can be used for matching a glob.
	// For example if a glob starts with `src/**/*.ts` we don't need to crawl all
	// folders in the current working directory as we see the `src` as base folder
	Base string
}

func (p *ParsedPattern) String() string {
	return p.RegExp.String()
}

func (p *ParsedPattern) setBase(str string) {
	if !p.isBaseSet {
		p.Base = str
		p.isBaseSet = true
	}
}

func (p *ParsedPattern) Compile() (*ParsedPattern, error) {
	re, err := regexp.Compile(`^` + p.stringPattern + `$`)
	if err != nil {
		return nil, err
	}
	p.RegExp = re
	return p, nil
}

func (p *ParsedPattern) add(char string) {
	p.stringPattern += char
}

func Parse(input string) (*ParsedPattern, error) {
	// remove duplicate or redundant parts in the pattern
	// and transform to posix separators
	input = cleanPattern(input)

	parsed := ParsedPattern{
		Input:      input,
		IsGlobstar: false,
	}

	for i := 0; i < len(input); i++ {
		cur := string(input[i])
		// store the previous character if we have some
		// var prevChar string
		// if i > 0 {
		// 	prevChar = string(input[i-1])
		// }

		// store the next character if we have some
		var nextChar string
		if i < len(input)-1 {
			nextChar = string(input[i+1])
		}

		switch cur {
		case "/":
			fallthrough
		case "$":
			fallthrough
		case "^":
			fallthrough
		case "+":
			fallthrough
		case ".":
			fallthrough
		case "(":
			fallthrough
		case ")":
			fallthrough
		case "=":
			fallthrough
		case "!":
			fallthrough
		case "|":
			// escape the following characters as they have a special meaning in regexp
			parsed.add(`\` + cur)

		// special characters are starting now
		case "?":
			parsed.add(`.`)

		case "*":
			// consider everything until now as base
			parsed.setBase(input[0:i])
			// count consecutive stars to determine if it is a globstar `**` or single star
			starCount := 1

			if nextChar == "*" {
				starCount++
				i++
			}

			// for string(input[i+1]) == "*" {
			// 	starCount++
			// 	i++
			// }

			// isStartOfSegment := prevChar == `/` || len(prevChar) == 0
			// isEndOfSegment := nextChar == `/` || len(nextChar) == 0

			isGlobstar := starCount > 1 //  && isStartOfSegment && isEndOfSegment

			if !parsed.IsGlobstar && isGlobstar {
				parsed.IsGlobstar = true
			}

			if isGlobstar {
				// it's a globstar, so match zero or more path segments
				parsed.add(internal.GLOBSTER)
				i++ // move over the "/"
			} else {
				// it's not a globstar, so only match one path segment
				parsed.add(internal.STAR)
			}

		default:
			parsed.add(cur)
		}

	}

	return parsed.Compile()
}

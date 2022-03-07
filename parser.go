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

func Parse(input string) (GlobPattern, error) {
	// remove duplicate or redundant parts in the pattern
	// and transform to posix separators
	input = cleanPattern(input)

	var str string

	for i := 0; i < len(input); i++ {
		cur := string(input[i])
		// store the previous character if we have some
		// var prevChar string
		// store the next character if we have some
		var nextChar string

		// if i > 0 {
		// 	prevChar = string(input[i-1])
		// }

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
			str += `\` + cur
		case "?":
			str += "."

		case "*":
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

			if isGlobstar {
				// it's a globstar, so match zero or more path segments
				str += internal.GLOBSTER
				i++ // move over the "/"
			} else {
				// it's not a globstar, so only match one path segment
				str += internal.STAR
			}

		default:
			str += cur
		}

	}

	return regexp.Compile(`^` + str + `$`)
}

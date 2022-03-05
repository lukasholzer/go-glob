package glob

import (
	"regexp"
	"strings"
)

func cleanPattern(input string) string {
	// windows path segments to posix characters
	input = strings.ReplaceAll(input, "\\", "/")
	input = strings.Replace(input, "***", "*", 1)
	input = strings.Replace(input, "**/**", "**", 1)
	input = strings.Replace(input, "**/**/**", "**", 1)

	return input
}

func Parse(input string) (*regexp.Regexp, error) {
	// remove duplicate or redundant parts in the pattern
	// and transform to posix separators
	input = cleanPattern(input)

	var str string

	for i := 0; i < len(input); i++ {
		cur := string(input[i])
		// store the previous character if we have some
		var prevChar string
		// store the next character if we have some
		var nextChar string

		if i > 0 {
			prevChar = string(input[i-1])
		}

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
			str += `\` + cur
		case "?":
			str += "."

		case "*":
			// count consecutive stars to determine if it is a globstar `**` or single star
			starCount := 1

			for string(input[i+1]) == "*" {
				starCount++
				i++
			}

			isStartOfSegment := prevChar == `/` || len(prevChar) == 0
			isEndOfSegment := nextChar == `/` || len(nextChar) == 0

			isGlobstar := starCount > 1 && isStartOfSegment && isEndOfSegment

			if string(input[i+1]) == "*" {
				isGlobstar = true
			}

			if isGlobstar {
				// it's a globstar, so match zero or more path segments
				str += `((?:[^/]*(?:\/|$))*)`
				i++ // move over the "/"
			} else {
				// it's not a globstar, so only match one path segment
				str += `([^/]*)`
			}

		default:
			str += cur
		}

	}

	return regexp.Compile(`^` + str + `$`)
}

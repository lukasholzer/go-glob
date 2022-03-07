package glob

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type GlobPatterns = map[string]*ParsedPattern

type Options struct {
	IgnorePatterns []string
	IgnoreFiles    []string
	CWD            string
	Patterns       []string
	AbsolutePaths  bool
	Debug          bool
}

func Pattern(pattern string) *Options {
	return &Options{
		Patterns: []string{pattern},
	}
}

// CWD sets a cwd where the glob pattern is executed in by default the process current working directory is used
func CWD(cwd string) *Options {
	return &Options{
		CWD: cwd,
	}
}

// IgnorePatterns sets a list of patterns that is used for ignoring specific files or directories
func IgnorePattern(pattern string) *Options {
	return &Options{
		IgnorePatterns: []string{pattern},
	}
}

// IgnoreFiles provides a list on file paths relative to the current working directory that follow the .gitignore syntax and are used to provide ignore patterns
func IgnoreFile(files string) *Options {
	return &Options{
		IgnoreFiles: []string{files},
	}
}

// Glob is the main glob function that is used to get a list of files in a specific directory matching the patterns and respecting the ignores
func Glob(options ...*Options) ([]string, error) {
	// use a map to avoid duplicates
	ignores := make(GlobPatterns)
	patterns := make(GlobPatterns)
	absolutePaths := false
	logger := zap.NewNop()
	var cwd string

	for _, opt := range options {
		if len(opt.CWD) > 0 {
			cwd = opt.CWD
		}

		if opt.AbsolutePaths {
			absolutePaths = true
		}

		if opt.Debug {
			l, err := zap.NewDevelopment()
			if err != nil {
				return nil, err
			}
			logger = l
		}

		for _, p := range opt.IgnorePatterns {
			reg, err := Parse(p)
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse ignorePattern %s", p)
			}
			ignores[p] = reg
		}

		for _, file := range opt.IgnoreFiles {
			gitIgnores, err := ParseGitignore(filepath.Join(cwd, file))
			if err != nil {
				return nil, errors.Wrapf(err, "could not parse ignore file %s", file)
			}

			for key, val := range gitIgnores {
				ignores[key] = val
			}
		}

		for _, p := range opt.Patterns {
			reg, err := Parse(p)

			if err != nil {
				return nil, errors.Wrapf(err, "could not parse provided pattern %s", p)
			}

			patterns[p] = reg
		}
	}

	if len(patterns) < 1 {
		return nil, errors.New("No patterns provided! Please provide a valid glob pattern as parameter")
	}

	if len(cwd) < 1 {
		wd, err := os.Getwd()
		if err != nil {
			return nil, errors.Wrap(err, "could not determine current working directory, please provide it as argument")
		}
		cwd = wd
	}

	return matchFiles(cwd, patterns, ignores, absolutePaths, logger)
}

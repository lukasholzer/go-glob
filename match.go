package glob

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func matchPattern(dir string, pattern *ParsedPattern, ignore GlobPatterns, absolutePaths bool, log *zap.Logger) ([]string, error) {
	inputPath := filepath.Join(dir, pattern.Input)
	if _, err := os.Stat(inputPath); err == nil {
		log.Sugar().Debugf("Pattern was an exact match %s", pattern.Input)
		if absolutePaths {
			return []string{inputPath}, nil
		}
		return []string{filepath.Join(pattern.Input)}, nil
	}

	fmt.Println()

	startDir := dir
	if len(pattern.Base) > 0 {
		startDir = filepath.Join(dir, pattern.Base)
	}

	logger := log.With(zap.String("base", startDir)).Sugar()

	var files []string
	err := filepath.WalkDir(startDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		for _, igr := range ignore {
			if Match(igr, rel) {
				logger.Debugf("Ignore Path %s", rel)

				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		if !d.IsDir() {
			if Match(pattern, rel) {
				logger.Debugf("Matching Path %s", rel)
				if absolutePaths {
					files = append(files, path)
				} else {
					files = append(files, rel)
				}
			}

		}
		return nil
	})
	return files, err
}

func matchFiles(dir string, globs GlobPatterns, ignore GlobPatterns, absolutePaths bool, logger *zap.Logger) ([]string, error) {
	var files []string
	var mu sync.Mutex
	var g errgroup.Group

	for _, pattern := range globs {
		pattern := pattern
		g.Go(func() error {
			matched, err := matchPattern(dir, pattern, ignore, absolutePaths, logger)
			mu.Lock()
			files = append(files, matched...)
			mu.Unlock()

			return err
		})
	}

	return files, g.Wait()
}

func Match(p *ParsedPattern, path string) bool {
	return p.RegExp.MatchString(path)
}

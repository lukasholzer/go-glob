package glob

import (
	"io/fs"
	"path/filepath"
)

func matchFiles(dir string, globs GlobPatterns, ignore GlobPatterns, absolutePaths bool) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		for _, igr := range ignore {
			if Match(igr, rel) {
				// fmt.Printf("IGNORE PATH: %s\n", rel)

				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		if !d.IsDir() {
			for _, fg := range globs {
				if Match(fg, rel) {
					// fmt.Printf("MATCHING PATH: %s\n", rel)
					if absolutePaths {
						files = append(files, path)
					} else {
						files = append(files, rel)
					}
				}
			}

		}
		return nil
	})
	return files, err
}

func Match(p GlobPattern, path string) bool {
	return p.MatchString(path)
}

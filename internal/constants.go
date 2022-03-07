package internal

const (
	STAR     = `([^/]*)`              // should match one path segment https://regex101.com/r/974GPG/1
	GLOBSTER = `((?:[^/]*(?:\/|$))*)` // should match multiple path segments https://regex101.com/r/7wy4wm/1
)

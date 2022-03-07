# glob

[![Main](https://github.com/lukasholzer/go-glob/actions/workflows/main.yml/badge.svg)](https://github.com/lukasholzer/go-glob/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/lukasholzer/go-glob/branch/main/graph/badge.svg?token=1RS0He2qXE)](https://codecov.io/gh/lukasholzer/go-glob)

## Usage

The most simple form would be by just providing a glob pattern

```golang
files, err := glob.Glob(glob.Pattern("**/*.txt"))
```

You can pass additionally to a pattern multiple options like `IgnoreFiles`. The Ignore files will parse the provided files according to the [Git Ignore Spec](http://git-scm.com/docs/gitignore) and will ignore the directories or files.

```golang
files, err := glob.Glob(glob.Pattern("**/*"), glob.IgnoreFile(".gitignore"))
```

The Ign

Advanced configuration use case:

```golang
files, err := glob.Glob(&glob.Options{
  Patterns:       []string{"**/*.txt"},
  CWD:            "/some/abs/path",
  IgnorePatterns: []string{"node_modules", "other/*.file"},
})
```

## Installation

```bash
go get -u github.com/lukasholzer/go-glob
```

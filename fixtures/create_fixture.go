package fixtures

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateNew(t *testing.T, fileMap map[string]string) string {
	tmpDir := t.TempDir()
	require.NoError(t, os.MkdirAll(tmpDir, 0777))

	for f, data := range fileMap {
		file := filepath.Join(tmpDir, f)
		dir := filepath.Dir(file)

		if dir != tmpDir {
			require.NoError(t, os.MkdirAll(dir, 0777))
		}
		require.NoError(t, os.WriteFile(file, []byte(data), 0644))
	}

	return tmpDir
}

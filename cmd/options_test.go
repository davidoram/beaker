package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetOptions(t *testing.T) {
	// Create a temporary directory scoped to this test (cleaned up automatically)
	dir := t.TempDir()

	// Create a temporary file within that directory
	tmpFile, err := os.CreateTemp(dir, "my-test-*.txt")
	require.NoError(t, err, "Failed to create temporary file")
	defer tmpFile.Close() //nolint:errcheck

	// Get the file path
	path := tmpFile.Name()
	t.Logf("path %s", path)
	// Test with direct arguments instead of modifying os.Args
	args := []string{"-credentials", path, "-postgres", "postgres://user:pass@localhost:5432/testdb", "-schema", filepath.Join("..", "schemas")}
	opts, err := ParseOptions(args)
	require.NoError(t, err)
	require.Equal(t, path, opts.CredentialsFile)
	require.Equal(t, "postgres://user:pass@localhost:5432/testdb", opts.PostgresURL)
}

func TestGetOptionsError(t *testing.T) {

	// Create a temporary directory scoped to this test (cleaned up automatically)
	dir := t.TempDir()

	// Create a temporary file within that directory
	tmpFile, err := os.CreateTemp(dir, "my-test-*.txt")
	require.NoError(t, err, "Failed to create temporary file")
	defer tmpFile.Close() //nolint:errcheck

	// Get the file path
	path := tmpFile.Name()

	scenarios := []struct {
		name        string
		args        []string
		expectedErr error
	}{
		{
			name:        "Empty credentials file",
			args:        []string{"-credentials", "", "-postgres", "postgres://user:pass@localhost:5432/testdb"},
			expectedErr: ErrBadCredentialsFile,
		},
		{
			name:        "Invalid credentials file",
			args:        []string{"-credentials", "/invalid/path/to/credentials.json", "-postgres", "postgres://user:pass@localhost:5432/testdb"},
			expectedErr: fs.ErrNotExist,
		},
		{
			name:        "Empty Postgres URL",
			args:        []string{"-credentials", path, "-postgres", ""},
			expectedErr: ErrBadPostgresURL,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			_, err = ParseOptions(scenario.args)
			require.ErrorIs(t, err, scenario.expectedErr)
		})
	}
}

package utility

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLoader_HappyPath(t *testing.T) {
	schemaDir := filepath.Join(".", "..", "..", "schemas")
	mappings := map[string]string{
		"http://github.com/davidoram/beaker/schemas/": schemaDir,
	}

	loader, err := NewLoader(mappings)
	require.NoError(t, err)
	require.NotNil(t, loader)

	files := []string{
		"http://github.com/davidoram/beaker/schemas/stock-add.request.json",
		"http://github.com/davidoram/beaker/schemas/product-sku.json",
		"https://json-schema.org/draft/2020-12/schema",
	}

	for _, file := range files {
		schema, err := loader.Load(file)
		require.NoError(t, err)
		require.NotNil(t, schema)
	}
}

func TestNewLoader_SchemaNotFound(t *testing.T) {
	schemaDir := filepath.Join(".", "..", "..", "schemas")
	mappings := map[string]string{
		"http://github.com/davidoram/beaker/schemas/": schemaDir,
	}

	loader, err := NewLoader(mappings)
	require.NoError(t, err)
	require.NotNil(t, loader)

	files := []string{
		"http://github.com/davidoram/beaker/schemas/does-not-exist.json",
		"https://json-schema.org/does-not-exist",
	}
	for _, file := range files {
		_, err = loader.Load(file)
		require.Error(t, err)
	}

}

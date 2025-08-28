package utility

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
)

func Test_NewJSONSchemaCompiler(t *testing.T) {
	compiler, err := NewJSONSchemaCompiler(t.Context(), filepath.Join(".", "..", "..", "schemas"))
	require.NotNil(t, compiler)
	require.NoError(t, err)

	// Check it can compile a schema
	schema, err := compiler.Compile("http://github.com/davidoram/beaker/schemas/stock-add.request.json")
	require.NoError(t, err)

	// Check the schema validates ok
	validJSON := `{"product-sku": "abc-123", "quantity": 5}`
	data, err := jsonschema.UnmarshalJSON(bytes.NewReader([]byte(validJSON)))
	require.NoError(t, err)
	require.NoError(t, schema.Validate(data))

	// Check the schema detects invalid JSON
	invalidJSON := `{"product-sku": "abc-123", "quantity": "five"}` // Note 'quantity' not a number
	data, err = jsonschema.UnmarshalJSON(bytes.NewReader([]byte(invalidJSON)))
	require.NoError(t, err)
	require.Error(t, schema.Validate(data))
}

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

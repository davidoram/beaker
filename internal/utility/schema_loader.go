package utility

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// NewJSONSchemaCompiler creates a new JSON schema compiler with the given schema directory.
func NewJSONSchemaCompiler(ctx context.Context, schemaDir string) (*jsonschema.Compiler, error) {
	loader, err := NewLoader(map[string]string{
		"http://github.com/davidoram/beaker/schemas/": schemaDir,
	})
	if err != nil {
		return nil, err
	}
	compiler := jsonschema.NewCompiler()
	compiler.UseLoader(loader)
	compiler.AssertContent()
	compiler.AssertFormat()
	compiler.DefaultDraft(jsonschema.Draft2020)
	return compiler, nil
}

// NewLoader creates a custom JSON schema loader that maps URL prefixes to local directories.
// It supports both HTTP and file-based loading.
// The mappings parameter is a map where keys are URL prefixes and values are local directory paths.
// It returns a jsonschema.URLLoader that can be used to load schemas.
// If a URL does not match any prefix, it falls back to the provided HTTP loader.
func NewLoader(mappings map[string]string) (jsonschema.URLLoader, error) {
	httpLoader := HTTPLoader(http.Client{
		Timeout: 15 * time.Second,
	})

	return &JVLoader{
		mappings: mappings,
		fallback: jsonschema.SchemeURLLoader{
			"file":  FileLoader{},
			"http":  &httpLoader,
			"https": &httpLoader,
		}}, nil
}

type JVLoader struct {
	mappings map[string]string
	fallback jsonschema.URLLoader
}

func (l *JVLoader) Load(url string) (any, error) {
	for prefix, dir := range l.mappings {
		if suffix, ok := strings.CutPrefix(url, prefix); ok {
			return loadFile(filepath.Join(dir, suffix))
		}
	}
	return l.fallback.Load(url)
}

func loadFile(path string) (any, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return jsonschema.UnmarshalJSON(f)
}

type FileLoader struct{}

func (l FileLoader) Load(url string) (any, error) {
	path, err := jsonschema.FileLoader{}.ToFile(url)
	if err != nil {
		return nil, err
	}
	return loadFile(path)
}

type HTTPLoader http.Client

func (l *HTTPLoader) Load(url string) (any, error) {
	client := (*http.Client)(l)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("%s returned status code %d", url, resp.StatusCode)
	}
	defer resp.Body.Close()
	return jsonschema.UnmarshalJSON(resp.Body)
}

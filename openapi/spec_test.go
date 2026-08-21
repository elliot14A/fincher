package openapi_test

import (
	"encoding/json"
	"testing"

	"github.com/elliot14A/fincher/openapi"
)

type swaggerDoc struct {
	Paths       map[string]any `json:"paths"`
	Definitions map[string]any `json:"definitions"`
}

func TestOpenAPISpec_ContentNotEmpty(t *testing.T) {
	if len(openapi.SpecJSON) == 0 {
		t.Fatal("openapi.SpecJSON is empty")
	}

	var doc swaggerDoc
	if err := json.Unmarshal(openapi.SpecJSON, &doc); err != nil {
		t.Fatalf("failed to parse openapi.SpecJSON: %v", err)
	}

	if len(doc.Paths) == 0 {
		t.Fatalf("expected non-empty paths in OpenAPI spec, got 0 paths")
	}
	if len(doc.Definitions) == 0 {
		t.Fatalf("expected non-empty definitions in OpenAPI spec, got 0 definitions")
	}

	// Verify all 6 core route groups exist
	expectedPaths := []string{
		"/titles",
		"/titles/{id}",
		"/masters",
		"/masters/{id}",
		"/vendors",
		"/vendors/{id}",
		"/packages",
		"/packages/{id}",
		"/deliveries",
		"/deliveries/{id}",
		"/dependencies",
		"/dependencies/{id}",
		"/dependencies/graph/{title_id}",
	}

	for _, p := range expectedPaths {
		if _, ok := doc.Paths[p]; !ok {
			t.Errorf("expected path %q in OpenAPI spec, but was not found", p)
		}
	}
}

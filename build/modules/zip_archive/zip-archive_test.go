package zip_archive_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/burbokop/design-practice-1/build/modules/zip_archive"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

func TestSimpleBinFactory(t *testing.T) {
	ctx := blueprint.NewContext()

	ctx.MockFileSystem(map[string][]byte{
		"Blueprints": []byte(`
			zip_archive {
				name: "test-zip",
			  	files: ["test-file1.txt", "test-file2.txt"]
			}
		`),
		"test-file1.txt": nil,
		"test-file2.txt": nil,
	})

	ctx.RegisterModuleType("zip_archive", zip_archive.SimpleBinFactory)

	cfg := bood.NewConfig()

	_, errs := ctx.ParseBlueprintsFiles(".", cfg)
	if len(errs) != 0 {
		t.Fatalf("Syntax errors in the test blueprint file: %s", errs)
	}

	_, errs = ctx.PrepareBuildActions(cfg)
	if len(errs) != 0 {
		t.Errorf("Unexpected errors while preparing build actions: %s", errs)
	}
	buffer := new(bytes.Buffer)
	if err := ctx.WriteBuildFile(buffer); err != nil {
		t.Errorf("Error writing ninja file: %s", err)
	} else {
		text := buffer.String()
		t.Logf("Gennerated ninja build file:\n%s", text)
		if !strings.Contains(text, "test-file1.txt") {
			t.Errorf("Generated ninja file does not have source: test-file1.txt")
		}
		if !strings.Contains(text, "test-file2.txt") {
			t.Errorf("Generated ninja file does not have source: test-file2.txt")
		}
	}
}

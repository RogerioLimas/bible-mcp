package mcpserver

import "testing"

func TestVerseInputSchemaEnums(t *testing.T) {
	schema := verseInputSchema([]string{"Almeida Revista e Corrigida"})

	book := schema.Properties["book"]
	if len(book.Enum) != 66 {
		t.Fatalf("book enum tem %d valores, want 66", len(book.Enum))
	}

	version := schema.Properties["version"]
	if len(version.Enum) != 1 || version.Enum[0] != "Almeida Revista e Corrigida" {
		t.Fatalf("version enum = %v, want [\"Almeida Revista e Corrigida\"]", version.Enum)
	}

	want := []string{"book", "chapter", "verse"}
	if len(schema.Required) != len(want) {
		t.Fatalf("Required = %v, want %v", schema.Required, want)
	}
	for i, name := range want {
		if schema.Required[i] != name {
			t.Fatalf("Required[%d] = %q, want %q", i, schema.Required[i], name)
		}
	}
}

func TestPassageInputSchemaRequired(t *testing.T) {
	schema := passageInputSchema([]string{"Almeida Revista e Corrigida"})
	want := []string{"book", "chapter", "verse_start", "verse_end"}
	if len(schema.Required) != len(want) {
		t.Fatalf("Required = %v, want %v", schema.Required, want)
	}
	for i, name := range want {
		if schema.Required[i] != name {
			t.Fatalf("Required[%d] = %q, want %q", i, schema.Required[i], name)
		}
	}
}

func TestChapterPropertyRejeitaZero(t *testing.T) {
	prop := chapterProperty()
	if prop.Minimum == nil || *prop.Minimum != 1 {
		t.Fatalf("chapterProperty().Minimum = %v, want 1", prop.Minimum)
	}
}

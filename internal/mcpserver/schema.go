package mcpserver

import (
	"github.com/google/jsonschema-go/jsonschema"

	"bible-mcp/internal/bible"
)

func stringEnum(values []string) []any {
	enum := make([]any, len(values))
	for i, v := range values {
		enum[i] = v
	}
	return enum
}

func minimum(n float64) *float64 {
	return &n
}

func bookProperty() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Enum:        stringEnum(bible.Names()),
		Description: "Nome canônico do livro bíblico",
	}
}

func versionProperty(versionNames []string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Enum:        stringEnum(versionNames),
		Description: "Versão da Bíblia; se omitido, usa " + bible.DefaultVersionName,
	}
}

func chapterProperty() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "integer",
		Minimum:     minimum(1),
		Description: "Número do capítulo",
	}
}

func verseProperty(description string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "integer",
		Minimum:     minimum(1),
		Description: description,
	}
}

func verseInputSchema(versionNames []string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"book":    bookProperty(),
			"chapter": chapterProperty(),
			"verse":   verseProperty("Número do versículo"),
			"version": versionProperty(versionNames),
		},
		Required: []string{"book", "chapter", "verse"},
	}
}

func passageInputSchema(versionNames []string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"book":        bookProperty(),
			"chapter":     chapterProperty(),
			"verse_start": verseProperty("Primeiro versículo do intervalo"),
			"verse_end":   verseProperty("Último versículo do intervalo"),
			"version":     versionProperty(versionNames),
		},
		Required: []string{"book", "chapter", "verse_start", "verse_end"},
	}
}

func chapterInputSchema(versionNames []string) *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"book":    bookProperty(),
			"chapter": chapterProperty(),
			"version": versionProperty(versionNames),
		},
		Required: []string{"book", "chapter"},
	}
}

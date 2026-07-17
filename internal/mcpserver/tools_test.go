package mcpserver

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"bible-mcp/internal/bible"
)

func openTestHandlers(t *testing.T) *Handlers {
	t.Helper()
	store, err := bible.Open("../../data")
	if err != nil {
		t.Fatalf("bible.Open(...) = %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return &Handlers{Store: store}
}

func TestGetVerse(t *testing.T) {
	h := openTestHandlers(t)
	res, _, err := h.GetVerse(context.Background(), nil, VerseInput{Book: "Gênesis", Chapter: 1, Verse: 1})
	if err != nil {
		t.Fatalf("GetVerse(...) = %v", err)
	}
	got := res.Content[0].(*mcp.TextContent).Text
	want := "1 No princípio, criou Deus os céus e a terra. \n- Gênesis 1:1"
	if got != want {
		t.Fatalf("GetVerse(...) texto = %q, want %q", got, want)
	}
}

func TestGetVerseLivroDesconhecido(t *testing.T) {
	h := openTestHandlers(t)
	if _, _, err := h.GetVerse(context.Background(), nil, VerseInput{Book: "Livro Inexistente", Chapter: 1, Verse: 1}); err == nil {
		t.Fatal("GetVerse(livro inexistente) não retornou erro")
	}
}

func TestGetPassage(t *testing.T) {
	h := openTestHandlers(t)
	res, _, err := h.GetPassage(context.Background(), nil, PassageInput{Book: "Gênesis", Chapter: 1, VerseStart: 1, VerseEnd: 2})
	if err != nil {
		t.Fatalf("GetPassage(...) = %v", err)
	}
	got := res.Content[0].(*mcp.TextContent).Text
	if !strings.HasSuffix(got, "- Gênesis 1:1-2") {
		t.Fatalf("GetPassage(...) texto = %q, want sufixo %q", got, "- Gênesis 1:1-2")
	}
}

func TestGetChapter(t *testing.T) {
	h := openTestHandlers(t)
	res, _, err := h.GetChapter(context.Background(), nil, ChapterInput{Book: "Gênesis", Chapter: 1})
	if err != nil {
		t.Fatalf("GetChapter(...) = %v", err)
	}
	got := res.Content[0].(*mcp.TextContent).Text
	if !strings.HasSuffix(got, "- Gênesis 1:1-31") {
		t.Fatalf("GetChapter(...) texto = %q, want sufixo %q", got, "- Gênesis 1:1-31")
	}
}

func TestGetVerseUsaVersaoExplicita(t *testing.T) {
	h := openTestHandlers(t)
	res, _, err := h.GetVerse(context.Background(), nil, VerseInput{
		Book: "João", Chapter: 3, Verse: 16, Version: "Almeida Revista e Atualizada",
	})
	if err != nil {
		t.Fatalf("GetVerse(...) = %v", err)
	}
	got := res.Content[0].(*mcp.TextContent).Text
	if !strings.HasPrefix(got, "16 ") {
		t.Fatalf("GetVerse(...) texto = %q, want prefixo %q", got, "16 ")
	}
}

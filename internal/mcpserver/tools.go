package mcpserver

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"bible-mcp/internal/bible"
)

// Handlers conecta os handlers das ferramentas MCP a um bible.Store.
type Handlers struct {
	Store *bible.Store
}

func (h *Handlers) resolveVersion(version string) string {
	if version == "" {
		return bible.DefaultVersionName
	}
	return version
}

func (h *Handlers) resolveBook(name string) (bible.Book, error) {
	book, ok := bible.BookByName(name)
	if !ok {
		return bible.Book{}, fmt.Errorf("livro canônico desconhecido: %q", name)
	}
	return book, nil
}

func quoteResult(bookName string, chapter int, verses []bible.VerseText) *mcp.CallToolResult {
	text := bible.FormatQuote(bookName, chapter, verses)
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: text}}}
}

// VerseInput é a entrada da ferramenta get_verse.
type VerseInput struct {
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Version string `json:"version,omitempty"`
}

// GetVerse resolve um único Versículo.
func (h *Handlers) GetVerse(_ context.Context, _ *mcp.CallToolRequest, in VerseInput) (*mcp.CallToolResult, any, error) {
	book, err := h.resolveBook(in.Book)
	if err != nil {
		return nil, nil, err
	}
	verses, err := h.Store.Verses(h.resolveVersion(in.Version), book, in.Chapter, in.Verse, in.Verse)
	if err != nil {
		return nil, nil, err
	}
	return quoteResult(book.Name, in.Chapter, verses), nil, nil
}

// PassageInput é a entrada da ferramenta get_passage.
type PassageInput struct {
	Book       string `json:"book"`
	Chapter    int    `json:"chapter"`
	VerseStart int    `json:"verse_start"`
	VerseEnd   int    `json:"verse_end"`
	Version    string `json:"version,omitempty"`
}

// GetPassage resolve uma Passagem: um intervalo de versículos dentro de um
// único capítulo.
func (h *Handlers) GetPassage(_ context.Context, _ *mcp.CallToolRequest, in PassageInput) (*mcp.CallToolResult, any, error) {
	book, err := h.resolveBook(in.Book)
	if err != nil {
		return nil, nil, err
	}
	verses, err := h.Store.Verses(h.resolveVersion(in.Version), book, in.Chapter, in.VerseStart, in.VerseEnd)
	if err != nil {
		return nil, nil, err
	}
	return quoteResult(book.Name, in.Chapter, verses), nil, nil
}

// ChapterInput é a entrada da ferramenta get_chapter.
type ChapterInput struct {
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Version string `json:"version,omitempty"`
}

// GetChapter resolve todos os Versículos de um Capítulo.
func (h *Handlers) GetChapter(_ context.Context, _ *mcp.CallToolRequest, in ChapterInput) (*mcp.CallToolResult, any, error) {
	book, err := h.resolveBook(in.Book)
	if err != nil {
		return nil, nil, err
	}
	verses, err := h.Store.Verses(h.resolveVersion(in.Version), book, in.Chapter, 0, 0)
	if err != nil {
		return nil, nil, err
	}
	return quoteResult(book.Name, in.Chapter, verses), nil, nil
}

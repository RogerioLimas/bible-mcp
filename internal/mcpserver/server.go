package mcpserver

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"bible-mcp/internal/bible"
)

// Register adiciona get_verse, get_passage e get_chapter a server, usando
// store como fonte de dados. Os schemas das ferramentas declaram um enum
// version construído a partir das Versões descobertas por store, e um enum
// book construído a partir de bible.Canon (docs/adr/0001,
// docs/adr/0002-external-data-directory.md).
func Register(server *mcp.Server, store *bible.Store) {
	h := &Handlers{Store: store}
	versionNames := store.VersionNames()

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_verse",
		Description: "Retorna um único versículo bíblico.",
		InputSchema: verseInputSchema(versionNames),
	}, h.GetVerse)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_passage",
		Description: "Retorna um intervalo de versículos dentro de um mesmo capítulo.",
		InputSchema: passageInputSchema(versionNames),
	}, h.GetPassage)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_chapter",
		Description: "Retorna todos os versículos de um capítulo.",
		InputSchema: chapterInputSchema(versionNames),
	}, h.GetChapter)
}

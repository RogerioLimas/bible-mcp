package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"bible-mcp/internal/bible"
	"bible-mcp/internal/mcpserver"
)

func main() {
	dir := os.Getenv("BIBLE_MCP_DATA_DIR")
	if dir == "" {
		log.Fatal("BIBLE_MCP_DATA_DIR não está configurada")
	}

	store, err := bible.Open(dir)
	if err != nil {
		log.Fatalf("falha ao carregar Versões de %s: %v", dir, err)
	}
	defer store.Close()

	server := mcp.NewServer(&mcp.Implementation{Name: "bible-mcp", Version: "0.1.0"}, nil)
	mcpserver.Register(server, store)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("servidor encerrado com erro: %v", err)
	}
}

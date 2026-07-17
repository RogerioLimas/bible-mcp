package mcpserver

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"bible-mcp/internal/bible"
)

func connectTestSession(t *testing.T, server *mcp.Server) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, serverTransport, nil); err != nil {
		t.Fatalf("server.Connect(...) = %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "test"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client.Connect(...) = %v", err)
	}
	t.Cleanup(func() { session.Close() })
	return session
}

func TestRegisterGetVerseEndToEnd(t *testing.T) {
	store, err := bible.Open("../../data")
	if err != nil {
		t.Fatalf("bible.Open(...) = %v", err)
	}
	t.Cleanup(func() { store.Close() })

	server := mcp.NewServer(&mcp.Implementation{Name: "bible-mcp-test", Version: "test"}, nil)
	Register(server, store)

	session := connectTestSession(t, server)
	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "get_verse",
		Arguments: map[string]any{"book": "Gênesis", "chapter": 1, "verse": 1},
	})
	if err != nil {
		t.Fatalf("CallTool(get_verse) = %v", err)
	}
	if res.IsError {
		t.Fatalf("CallTool(get_verse) retornou IsError=true, content=%v", res.Content)
	}
	got := res.Content[0].(*mcp.TextContent).Text
	want := "1 No princípio, criou Deus os céus e a terra. \n- Gênesis 1:1"
	if got != want {
		t.Fatalf("CallTool(get_verse) texto = %q, want %q", got, want)
	}
}

func TestRegisterRejeitaLivroForaDoEnum(t *testing.T) {
	store, err := bible.Open("../../data")
	if err != nil {
		t.Fatalf("bible.Open(...) = %v", err)
	}
	t.Cleanup(func() { store.Close() })

	server := mcp.NewServer(&mcp.Implementation{Name: "bible-mcp-test", Version: "test"}, nil)
	Register(server, store)

	session := connectTestSession(t, server)
	res, err := session.CallTool(context.Background(), &mcp.CallToolParams{
		Name:      "get_verse",
		Arguments: map[string]any{"book": "Not A Real Book", "chapter": 1, "verse": 1},
	})
	if err != nil {
		t.Fatalf("CallTool(get_verse) erro de transporte = %v", err)
	}
	if !res.IsError {
		t.Fatal("CallTool(get_verse) com livro fora do enum não retornou IsError=true")
	}
}

func TestRegisterExpoeTresFerramentas(t *testing.T) {
	store, err := bible.Open("../../data")
	if err != nil {
		t.Fatalf("bible.Open(...) = %v", err)
	}
	t.Cleanup(func() { store.Close() })

	server := mcp.NewServer(&mcp.Implementation{Name: "bible-mcp-test", Version: "test"}, nil)
	Register(server, store)

	session := connectTestSession(t, server)
	res, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListTools(...) = %v", err)
	}
	if len(res.Tools) != 3 {
		t.Fatalf("ListTools() retornou %d ferramentas, want 3", len(res.Tools))
	}
}

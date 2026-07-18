# Graph Report - .  (2026-07-18)

## Corpus Check
- 38 files · ~33,991 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 147 nodes · 325 edges · 13 communities (11 shown, 2 thin omitted)
- Extraction: 85% EXTRACTED · 15% INFERRED · 0% AMBIGUOUS · INFERRED: 49 edges (avg confidence: 0.84)
- Token cost: 198,373 input · 0 output

## Community Hubs (Navigation)
- Canon & Query Tests
- Domain Docs & ADRs
- Version Discovery (Store)
- Implementation Plan Docs
- Input Schemas
- MCP Tool Handlers
- Server Bootstrap & Registration
- Quote Formatting
- SPDD Command Specs
- Handler Tests
- Phase Gate Hook
- Go Module

## God Nodes (most connected - your core abstractions)
1. `Open()` - 13 edges
2. `bible-mcp implementation plan (superpowers TDD plan)` - 13 edges
3. `CLAUDE.md project guidance` - 12 edges
4. `BookByName()` - 10 edges
5. `Register()` - 10 edges
6. `SPDD REASONS Canvas prompt: get_verse/get_passage/get_chapter MCP tools` - 10 edges
7. `openTestStore()` - 9 edges
8. `openTestHandlers()` - 9 edges
9. `CONTEXT.md ubiquitous language document` - 9 edges
10. `verseInputSchema()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `Open()`  [INFERRED]
  main.go → internal/bible/store.go
- `bible-mcp implementation plan (superpowers TDD plan)` --references--> `superpowers:test-driven-development skill (only after SPDD prompt exists, on explicit request)`  [EXTRACTED]
  docs/superpowers/plans/2026-07-15-bible-mcp.md → CLAUDE.md
- `bible-mcp implementation plan (superpowers TDD plan)` --references--> `superpowers:executing-plans skill (only after SPDD prompt exists, on explicit request)`  [EXTRACTED]
  docs/superpowers/plans/2026-07-15-bible-mcp.md → CLAUDE.md
- `main()` --calls--> `Register()`  [INFERRED]
  main.go → internal/mcpserver/server.go
- `CLAUDE.md project guidance` --references--> `/spdd-analysis command (Claude)`  [EXTRACTED]
  CLAUDE.md → .claude/commands/spdd-analysis.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **SPDD Workflow Pipeline (Analysis to Reasons Canvas to Generate to Sync)** — claude_commands_spdd_analysis_command, claude_commands_spdd_reasons_canvas_command, claude_commands_spdd_generate_command, claude_commands_spdd_sync_command [EXTRACTED 1.00]
- **Bible domain ubiquitous language (Livro, Capitulo, Versiculo, Passagem, Versao)** — context_livro, context_capitulo, context_versiculo, context_passagem, context_versao [EXTRACTED 1.00]
- **Architectural decision set governing bible-mcp domain model (ADR-0001, ADR-0002, CONTEXT.md)** — docs_adr_0001_canonical_book_identity, docs_adr_0002_external_data_directory, context_doc [INFERRED 0.85]

## Communities (13 total, 2 thin omitted)

### Community 0 - "Canon & Query Tests"
Cohesion: 0.22
Nodes (16): Book, Testament, BookByName(), Names(), T, TestBookByName(), TestBookByNameDesconhecido(), TestCanonTem66Livros() (+8 more)

### Community 1 - "Domain Docs & ADRs"
Cohesion: 0.22
Nodes (19): CLAUDE.md project guidance, superpowers:brainstorming skill (excluded during SPDD planning phase per project policy), superpowers:executing-plans skill (only after SPDD prompt exists, on explicit request), superpowers:test-driven-development skill (only after SPDD prompt exists, on explicit request), superpowers:writing-plans skill (excluded during SPDD planning phase per project policy), Capítulo (numbered subdivision of a Livro), CONTEXT.md ubiquitous language document, Livro (canonical Bible book concept) (+11 more)

### Community 2 - "Version Discovery (Store)"
Cohesion: 0.27
Nodes (12): version, DB, Store, metadataName(), Open(), openVersion(), copyFile(), T (+4 more)

### Community 3 - "Implementation Plan Docs"
Cohesion: 0.34
Nodes (16): Task 1: Go module and canonical Book catalog (Canon/BookByName), bible-mcp implementation plan (superpowers TDD plan), Task 4: Quote formatting (FormatQuote/Reference), Task 6: MCP tool handlers (GetVerse/GetPassage/GetChapter), Task 7: MCP server registration and binary (Register/main), Task 5: Dynamic tool input schemas, Task 2: Version loading and validation (Store/Open), Task 3: Verse query (Store.Verses) (+8 more)

### Community 4 - "Input Schemas"
Cohesion: 0.33
Nodes (14): bookProperty(), chapterInputSchema(), chapterProperty(), minimum(), passageInputSchema(), stringEnum(), T, TestChapterPropertyRejeitaZero() (+6 more)

### Community 5 - "MCP Tool Handlers"
Cohesion: 0.31
Nodes (9): CallToolRequest, CallToolResult, Context, Store, quoteResult(), ChapterInput, Handlers, PassageInput (+1 more)

### Community 6 - "Server Bootstrap & Registration"
Cohesion: 0.23
Nodes (11): ClientSession, Server, Store, Register(), connectTestSession(), Server, T, TestRegisterExpoeTresFerramentas() (+3 more)

### Community 7 - "Quote Formatting"
Cohesion: 0.29
Nodes (7): VerseText, FormatQuote(), Reference(), T, TestFormatQuotePassagem(), TestFormatQuoteVersiculoUnico(), Store

### Community 8 - "SPDD Command Specs"
Cohesion: 0.47
Nodes (10): /spdd-analysis command (Claude), /spdd-generate command (Claude), /spdd-prompt-update command (Claude), /spdd-reasons-canvas command (Claude), /spdd-sync command (Claude), /spdd-analysis command (OpenCode), /spdd-generate command (OpenCode), /spdd-prompt-update command (OpenCode) (+2 more)

### Community 9 - "Handler Tests"
Cohesion: 0.61
Nodes (7): T, openTestHandlers(), TestGetChapter(), TestGetPassage(), TestGetVerse(), TestGetVerseLivroDesconhecido(), TestGetVerseUsaVersaoExplicita()

## Knowledge Gaps
- **5 isolated node(s):** `phase-gate.sh script`, `bible-mcp`, `Store`, `superpowers:brainstorming skill (excluded during SPDD planning phase per project policy)`, `superpowers:writing-plans skill (excluded during SPDD planning phase per project policy)`
  These have ≤1 connection - possible missing edges or undocumented components.
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Open()` connect `Version Discovery (Store)` to `Canon & Query Tests`, `Handler Tests`, `Server Bootstrap & Registration`?**
  _High betweenness centrality (0.223) - this node is a cross-community bridge._
- **Why does `Register()` connect `Server Bootstrap & Registration` to `Input Schemas`?**
  _High betweenness centrality (0.114) - this node is a cross-community bridge._
- **Why does `openTestHandlers()` connect `Handler Tests` to `Version Discovery (Store)`, `MCP Tool Handlers`?**
  _High betweenness centrality (0.104) - this node is a cross-community bridge._
- **Are the 9 inferred relationships involving `Open()` (e.g. with `openTestStore()` and `TestOpenDescobreTodasAsVersoes()`) actually correct?**
  _`Open()` has 9 INFERRED edges - model-reasoned connections that need verification._
- **Are the 8 inferred relationships involving `BookByName()` (e.g. with `TestBookByName()` and `TestBookByNameDesconhecido()`) actually correct?**
  _`BookByName()` has 8 INFERRED edges - model-reasoned connections that need verification._
- **Are the 7 inferred relationships involving `Register()` (e.g. with `chapterInputSchema()` and `passageInputSchema()`) actually correct?**
  _`Register()` has 7 INFERRED edges - model-reasoned connections that need verification._
- **What connects `phase-gate.sh script`, `bible-mcp`, `Store` to the rest of the system?**
  _5 weakly-connected nodes found - possible documentation gaps or missing edges._
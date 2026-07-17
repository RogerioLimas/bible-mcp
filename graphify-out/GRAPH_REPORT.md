# Graph Report - .  (2026-07-15)

## Corpus Check
- Corpus is ~457 words - fits in a single context window. You may not need a graph.

## Summary
- 11 nodes · 22 edges · 2 communities
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 30,000 input · 6,051 output

## Community Hubs (Navigation)
- Esquema de Dados Biblicos (SQLite)
- Arquitetura do Servidor MCP

## God Nodes (most connected - your core abstractions)
1. `data/ARA.sqlite (Almeida Revista e Atualizada)` - 9 edges
2. `data/ARC.sqlite (Almeida Revista e Corrigida)` - 9 edges
3. `MCP Server (to be implemented)` - 5 edges
4. `book table (66 books, testament/book reference ids)` - 4 edges
5. `bible-mcp Project` - 3 edges
6. `verse table (individual verses with chapter/verse/text)` - 3 edges
7. `Data access layer must parametrize which .sqlite version is used` - 3 edges
8. `book.id / book_reference_id as stable cross-version reference key` - 3 edges
9. `metadata table (key/value per database)` - 2 edges
10. `Sociedade Bíblica do Brasil (2009 copyright holder)` - 2 edges

## Surprising Connections (you probably didn't know these)
- `bible-mcp Project` --conceptually_related_to--> `data/ARA.sqlite (Almeida Revista e Atualizada)`  [EXTRACTED]
  CLAUDE.md → CLAUDE.md  _Bridges community 1 → community 0_

## Hyperedges (group relationships)
- **ARA and ARC databases share identical 3-table schema (metadata, book, verse)** — claude_ara_sqlite, claude_arc_sqlite, claude_metadata_table, claude_book_table, claude_verse_table [INFERRED 0.85]
- **MCP server exposes ARA/ARC bible text as MCP tools/resources** — claude_mcp_server, claude_model_context_protocol, claude_ara_sqlite, claude_arc_sqlite [EXTRACTED 1.00]

## Communities (2 total, 0 thin omitted)

### Community 0 - "Esquema de Dados Biblicos (SQLite)"
Cohesion: 0.62
Nodes (7): data/ARA.sqlite (Almeida Revista e Atualizada), data/ARC.sqlite (Almeida Revista e Corrigida), book.id / book_reference_id as stable cross-version reference key, book table (66 books, testament/book reference ids), metadata table (key/value per database), Sociedade Bíblica do Brasil (2009 copyright holder), verse table (individual verses with chapter/verse/text)

### Community 1 - "Arquitetura do Servidor MCP"
Cohesion: 0.50
Nodes (4): bible-mcp Project, Data access layer must parametrize which .sqlite version is used, MCP Server (to be implemented), Model Context Protocol (MCP)

## Knowledge Gaps
- **1 isolated node(s):** `Model Context Protocol (MCP)`
  These have ≤1 connection - possible missing edges or undocumented components.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `data/ARA.sqlite (Almeida Revista e Atualizada)` connect `Esquema de Dados Biblicos (SQLite)` to `Arquitetura do Servidor MCP`?**
  _High betweenness centrality (0.315) - this node is a cross-community bridge._
- **Why does `data/ARC.sqlite (Almeida Revista e Corrigida)` connect `Esquema de Dados Biblicos (SQLite)` to `Arquitetura do Servidor MCP`?**
  _High betweenness centrality (0.315) - this node is a cross-community bridge._
- **Why does `MCP Server (to be implemented)` connect `Arquitetura do Servidor MCP` to `Esquema de Dados Biblicos (SQLite)`?**
  _High betweenness centrality (0.207) - this node is a cross-community bridge._
- **What connects `Model Context Protocol (MCP)` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Estado atual do projeto

A v1 do servidor MCP está implementada em Go. Stack:

- Go 1.25+ (ver `go.mod`).
- `github.com/modelcontextprotocol/go-sdk/mcp` — SDK oficial de MCP, transporte stdio (`mcp.StdioTransport`).
- `github.com/google/jsonschema-go/jsonschema` — schemas de entrada das ferramentas, com enum dinâmico de `book`/`version`.
- `modernc.org/sqlite` — driver SQLite puro Go (sem cgo), acesso somente leitura (`mode=ro`).

Estrutura:

- `main.go` — binário: lê `BIBLE_MCP_DATA_DIR`, abre o `Store`, registra as ferramentas, roda o transporte stdio.
- `internal/bible` — camada de domínio: `Canon`/`BookByName` (catálogo canônico dos 66 livros), `Store`/`Open` (descoberta e validação de Versões `.sqlite`), `Verses` (consulta unificada de Versículo/Passagem/Capítulo), `FormatQuote` (formatação de saída).
- `internal/mcpserver` — camada de protocolo: `Handlers` (`GetVerse`, `GetPassage`, `GetChapter`), schemas de entrada, `Register` (registro das 3 ferramentas no `*mcp.Server`).
- `data/ARA.sqlite`, `data/ARC.sqlite`, `data/NVI.sqlite` — as três Versões (ARA/ARC/NVI) usadas como diretório de dados padrão.

Ferramentas MCP expostas: `get_verse`, `get_passage`, `get_chapter` (ver `spdd/prompt/GGQPA-XXX-202607162106-[Feat]-mcp-verse-passage-chapter-tools.md` para a especificação completa e `README.md` para exemplos de uso).

Configuração obrigatória: variável de ambiente `BIBLE_MCP_DATA_DIR` apontando para um diretório com arquivos `.sqlite` no esquema descrito abaixo — o servidor falha ao iniciar (`log.Fatal`) se a variável estiver vazia, o diretório não tiver nenhum `.sqlite` válido, ou qualquer arquivo reprovar a validação de alinhamento com o cânone (ver `docs/adr/0001-canonical-book-identity.md` e `docs/adr/0002-external-data-directory.md`).

Comandos:

```bash
go build -o bible-mcp .   # build do binário
go test ./...             # suíte de testes (internal/bible, internal/mcpserver)
```

## Dados: esquema dos bancos SQLite

Os três arquivos `.sqlite` em `data/` compartilham exatamente o mesmo esquema (3 tabelas):

```sql
CREATE TABLE metadata (
    "key" VARCHAR(255) NOT NULL,
    value VARCHAR(255),
    PRIMARY KEY ("key")
)

CREATE TABLE book (
    id INTEGER NOT NULL,
    book_reference_id INTEGER,
    testament_reference_id INTEGER,
    name VARCHAR(50),
    PRIMARY KEY (id)
)

CREATE TABLE verse (
    id INTEGER NOT NULL,
    book_id INTEGER,
    chapter INTEGER,
    verse INTEGER,
    text TEXT,
    PRIMARY KEY (id),
    FOREIGN KEY(book_id) REFERENCES book (id)
)
```

Detalhes relevantes:

- `book.testament_reference_id`: `1` = Antigo Testamento, `2` = Novo Testamento (66 livros no total, nas três versões).
- `book.book_reference_id`: identificador estável do livro, independente da versão (útil para cruzar referências entre ARA, ARC e NVI — os `book.id` e `book_reference_id`/`testament_reference_id` coincidem posição a posição nos três bancos para o mesmo livro; a única divergência observada é a grafia de `book.name`, ver `docs/adr/0001-canonical-book-identity.md`).
- `verse`: cada linha é um versículo individual, referenciando `book_id`, `chapter` e `verse`. O texto em `text` inclui espaços/pontuação originais (não normalizado).
- `metadata`: tabela chave/valor por banco, contendo `language_id`, `version`, `name` (nome completo da versão) e `copyright`. Decisão de produto: as ferramentas do MCP retornam apenas o texto puro do versículo/passagem/capítulo — nenhuma nota de copyright ou texto adicional é incluída nas respostas.

Todas as três versões têm 66 livros; ARA tem 31.104 versículos, ARC tem 31.105 e NVI tem 31.087 (pequena diferença de versificação entre as edições — não assumir que os `verse.id` são idênticos entre os bancos).

## Ao implementar o servidor MCP

- As três versões (ARA/ARC/NVI) devem provavelmente ser tratadas como bancos de dados intercambiáveis por um mesmo conjunto de queries — qualquer camada de acesso a dados deve descobrir/parametrizar qual arquivo `.sqlite` usar em vez de fixar um dos três (ver `docs/adr/0002-external-data-directory.md`).
- Como os `book.id`/`book_reference_id` coincidem entre as três versões para o mesmo livro, é seguro usar esse id como chave de referência ao permitir comparar/alternar entre versões de uma mesma passagem.

## Agent skills

### Issue tracker

Issues are tracked in GitHub Issues. See `docs/agents/issue-tracker.md`.

### Triage labels

Five canonical roles: `needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context layout: ubiquitous language in `CONTEXT.md`, architectural decisions in `docs/adr/`. See `docs/agents/domain.md`.

### Fluxo de planejamento

Fase de planejamento usa exclusivamente grill-with-docs + open-spdd (`/spdd-analysis`, `/spdd-reasons-canvas`). NÃO invocar `superpowers:brainstorming`, `domain-modeling` ou `writing-plans` nessa fase.
`superpowers:test-driven-development`, `executing-plans`, `subagent-driven-development` só entram depois que o prompt SPDD existir e o usuário pedir explicitamente para implementar.

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Estado atual do projeto

Este repositório ainda não contém código-fonte, build tooling, testes ou um `package.json`/manifest de qualquer runtime. O único conteúdo existente é:

- `data/ARA.sqlite` — texto bíblico da versão "Almeida Revista e Atualizada"
- `data/ARC.sqlite` — texto bíblico da versão "Almeida Revista e Corrigida"
- `data/NVI.sqlite` — texto bíblico da versão "Nova Versão Internacional"

O nome do projeto (`bible-mcp`) e o conteúdo de `data/` indicam que o objetivo é implementar um servidor MCP (Model Context Protocol) que exponha esse texto bíblico como ferramentas/recursos para um LLM. Nenhuma implementação de servidor MCP existe ainda — isso ainda está por ser construído.

Não há comandos de build, lint ou teste a documentar até que o código seja criado. Ao iniciar a implementação, atualize esta seção com o stack escolhido (linguagem/runtime MCP SDK) e os comandos correspondentes.

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

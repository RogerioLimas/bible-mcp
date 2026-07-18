# bible-mcp

Servidor [MCP](https://modelcontextprotocol.io) somente-leitura, escrito em Go, que expõe o texto bíblico em português (versões **ARA**, **ARC** e **NVI**) como ferramentas para agentes de IA — busca de um versículo único, de uma passagem (intervalo de versículos) ou de um capítulo inteiro.

## Sumário

- [Sobre](#sobre)
- [Ferramentas MCP](#ferramentas-mcp)
- [Arquitetura](#arquitetura)
- [Requisitos](#requisitos)
- [Instalação](#instalação)
- [Configuração](#configuração)
- [Uso](#uso)
- [Desenvolvimento](#desenvolvimento)
- [Estrutura do projeto](#estrutura-do-projeto)
- [Domínio e decisões de arquitetura](#domínio-e-decisões-de-arquitetura)
- [Roadmap](#roadmap)

## Sobre

`bible-mcp` conecta um LLM a três versões do texto bíblico em português — Almeida Revista e Atualizada (ARA), Almeida Revista e Corrigida (ARC) e Nova Versão Internacional (NVI) — através do Model Context Protocol. O servidor:

- Resolve livros por um **nome canônico fixo** em português (enum de 66 valores), independente de como cada Versão grafa o nome internamente (ver [ADR-0001](docs/adr/0001-canonical-book-identity.md)).
- Descobre as Versões disponíveis **em tempo de execução**, varrendo um diretório de dados configurável — nenhum arquivo `.sqlite` é embutido no binário (ver [ADR-0002](docs/adr/0002-external-data-directory.md)).
- Nunca inclui metadata de copyright ou qualquer texto além do versículo/passagem/capítulo pedido nas respostas das ferramentas.
- Falha rápido na inicialização se qualquer Versão descoberta não corresponder exatamente ao cânone de 66 livros — nunca inicia em modo parcial/degradado.

## Ferramentas MCP


| Ferramenta    | Descrição                                                    | Parâmetros obrigatórios                     | Parâmetros opcionais |
| ------------- | -------------------------------------------------------------- | --------------------------------------------- | --------------------- |
| `get_verse`   | Retorna um único versículo.                                  | `book`, `chapter`, `verse`                    | `version`             |
| `get_passage` | Retorna um intervalo de versículos dentro do mesmo capítulo. | `book`, `chapter`, `verse_start`, `verse_end` | `version`             |
| `get_chapter` | Retorna todos os versículos de um capítulo.                  | `book`, `chapter`                             | `version`             |

- `book`: enum fixo com os 66 nomes canônicos (ex: `"Gênesis"`, `"1Coríntios"`).
- `version`: enum dinâmico, construído a partir das Versões descobertas em `BIBLE_MCP_DATA_DIR`. Se omitido, usa `"Almeida Revista e Corrigida"`.

Toda resposta de sucesso vem no formato "quote": uma linha por versículo (numerado), seguida de uma linha de atribuição.

```
$ get_verse(book="João", chapter=3, verse=16)

16 Porque Deus amou o mundo de tal maneira que deu o seu Filho unigênito, para que todo aquele que nele crê não pereça, mas tenha a vida eterna.
- João 3:16
```

Qualquer referência inválida (livro fora do enum, `verse_start > verse_end`, capítulo/versículo inexistente na Versão) retorna um erro de ferramenta MCP — nunca uma resposta de sucesso vazia.

## Arquitetura

```
main.go               → lê BIBLE_MCP_DATA_DIR, abre o Store, registra as ferramentas, roda stdio
internal/mcpserver/    → camada de protocolo: schemas de entrada, Handlers, Register
internal/bible/        → camada de domínio: Canon, Store (descoberta/validação), Verses, FormatQuote
<versão>.sqlite         → camada de dados: acesso somente leitura via database/sql + modernc.org/sqlite
```

Nenhuma string de entrada do agente é concatenada em SQL (todos os parâmetros passam por `?`); toda conexão é aberta em `mode=ro`. Decisões de design completas em [Domínio e decisões de arquitetura](#domínio-e-decisões-de-arquitetura).

## Requisitos

- Go 1.25 ou superior.
- Um ou mais arquivos `.sqlite` no esquema `metadata`/`book`/`verse` descrito em [CLAUDE.md](CLAUDE.md#dados-esquema-dos-bancos-sqlite) — o repositório já inclui `data/ARA.sqlite`, `data/ARC.sqlite` e `data/NVI.sqlite`.

## Instalação

```bash
git clone git@github.com:RogerioLimas/bible-mcp.git
cd bible-mcp
go build -o bible-mcp .
```

## Configuração

O servidor é configurado por uma única variável de ambiente:


| Variável            | Obrigatória | Descrição                                                     |
| -------------------- | ------------ | --------------------------------------------------------------- |
| `BIBLE_MCP_DATA_DIR` | Sim          | Diretório contendo os arquivos`*.sqlite` das Versões a expor. |

O identificador de cada Versão (valor aceito pelo parâmetro `version`) vem do campo `metadata.name` de dentro do próprio arquivo, não do nome do arquivo. Na inicialização, cada arquivo descoberto é validado contra o cânone de 66 livros; se a variável estiver vazia, o diretório não tiver nenhum `.sqlite` válido, ou qualquer arquivo reprovar a validação (ou duas Versões tiverem o mesmo `metadata.name`), o processo termina imediatamente com erro — não há modo de inicialização parcial.

## Uso

Executando o binário diretamente, usando os dados já incluídos no repositório:

```bash
BIBLE_MCP_DATA_DIR=./data ./bible-mcp
```

O servidor conversa via stdio, seguindo o Model Context Protocol — para uso interativo, registre-o em um cliente MCP (ex: Claude Desktop, `claude mcp add`, ou qualquer cliente compatível). Exemplo de configuração:

```json
{
  "mcpServers": {
    "bible": {
      "command": "/caminho/absoluto/para/bible-mcp",
      "env": {
        "BIBLE_MCP_DATA_DIR": "/caminho/absoluto/para/bible-mcp/data"
      }
    }
  }
}
```

## Desenvolvimento

```bash
go test ./...
```

Convenções do projeto (nomenclatura, tratamento de erro, validação): funções de domínio retornam `error` como último valor (nunca `panic` para entrada inválida); handlers de ferramenta MCP nunca traduzem um erro de domínio em sucesso — sempre `(nil, nil, err)`, que o SDK converte em `CallToolResult{IsError: true}`. Especificação completa em [spdd/prompt/GGQPA-XXX-202607162106-\[Feat\]-mcp-verse-passage-chapter-tools.md](spdd/prompt/GGQPA-XXX-202607162106-%5BFeat%5D-mcp-verse-passage-chapter-tools.md).

## Estrutura do projeto

```
.
├── main.go                      # binário: bootstrap e transporte stdio
├── internal/
│   ├── bible/                   # domínio: cânone, Store, consulta, formatação
│   │   ├── canon.go
│   │   ├── store.go
│   │   ├── query.go
│   │   └── format.go
│   └── mcpserver/                # protocolo: schemas, handlers, registro
│       ├── schema.go
│       ├── tools.go
│       └── server.go
├── data/                          # Versões .sqlite padrão (ARA/ARC/NVI)
├── docs/
│   ├── adr/                       # Architecture Decision Records
│   └── BACKLOG.md                 # features adiadas para versões futuras
├── spdd/                          # especificação estruturada (REASONS Canvas)
└── CONTEXT.md                     # linguagem ubíqua do domínio
```

## Domínio e decisões de arquitetura

- [`CONTEXT.md`](CONTEXT.md) — linguagem ubíqua (Livro, Nome Canônico, Capítulo, Versículo, Passagem, Versão).
- [`docs/adr/0001-canonical-book-identity.md`](docs/adr/0001-canonical-book-identity.md) — por que a identidade de livro é um enum canônico fixo, e não `book.name` de cada Versão.
- [`docs/adr/0002-external-data-directory.md`](docs/adr/0002-external-data-directory.md) — por que as Versões são um diretório externo descoberto em runtime, com falha rápida na validação.

## Roadmap

Itens deliberadamente adiados para uma versão futura (ver [`docs/BACKLOG.md`](docs/BACKLOG.md)):

- Busca por texto (FTS5) — encontrar versículos por conteúdo, não apenas por referência.
- Formato de resposta configurável — hoje fixo em "quote"; uma versão futura deve permitir formatos alternativos (ex: JSON estruturado).

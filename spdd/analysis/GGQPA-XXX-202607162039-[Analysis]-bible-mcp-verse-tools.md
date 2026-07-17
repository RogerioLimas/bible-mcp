# SPDD Analysis: Servidor MCP de texto bíblico (get_verse/get_passage/get_chapter sobre Versões intercambiáveis)

## Original Business Requirement

> Conteúdo consolidado, verbatim, dos arquivos referenciados: `CONTEXT.md`, `docs/adr/0001-canonical-book-identity.md`, `docs/adr/0002-external-data-directory.md`, `docs/BACKLOG.md`.

### CONTEXT.md

```markdown
# bible-mcp

Servidor MCP que expõe o texto bíblico (versões ARA/ARC/NVI) como ferramentas para agentes de IA buscarem versículos, passagens ou capítulos inteiros.

## Language

**Livro**:
Uma das 66 unidades bíblicas canônicas, composta por um ou mais Capítulos e pertencente a um Testamento (Antigo ou Novo). Identificado externamente por um Nome Canônico fixo, resolvido internamente para o registro correspondente (via `book_reference_id`) em cada Versão.
_Avoid_: Livro bíblico (redundante)

**Nome Canônico**:
O nome fixo em português usado para identificar um Livro nas ferramentas do MCP, independente da grafia usada pela `book.name` de uma Versão específica (ex: sempre "1 Coríntios", mesmo que uma Versão grave "1Corintios" internamente).
_Avoid_: Nome do livro (ambíguo entre canônico e específico da Versão), título

**Capítulo**:
Subdivisão numerada de um Livro, contendo um ou mais Versículos.

**Versículo**:
Unidade atômica de texto bíblico, endereçada pela combinação (Livro, Capítulo, número do versículo).
_Avoid_: Verse

**Passagem**:
Um intervalo contíguo de Versículos dentro de um mesmo Capítulo, delimitado por um versículo inicial e um versículo final. Um Versículo único é uma Passagem degenerada (início == fim); um Capítulo inteiro é uma Passagem sem limites definidos (início e fim ausentes).
_Avoid_: Trecho, range

**Versão**:
A edição do texto bíblico (ex: ARA, Almeida Revista e Atualizada; ARC, Almeida Revista e Corrigida). Mesmo conjunto de Livros e Capítulos entre Versões, mas texto e versificação (contagem/limites de Versículos) podem diferir levemente.
_Avoid_: Tradução, edição
```

### docs/adr/0001-canonical-book-identity.md

```markdown
# Identidade de Livro canônica e independente de Versão

As ferramentas do MCP precisam identificar um Livro para qualquer Versão (ARA, ARC, NVI, e futuras). Casar diretamente com a coluna `book.name` de cada base foi rejeitado: validação nos dados reais mostrou que o livro `book_reference_id=22` já tem 3 grafias diferentes entre as Versões atuais ("Cântico", "Cantares", "Cântico dos Cânticos"), e isso só cresceria com novas Versões/idiomas. Em vez disso, o parâmetro `book` exposto ao agente é um enum fixo de Nomes Canônicos em português (definidos pelo servidor, não lidos de nenhuma base), resolvido internamente para o registro de cada Versão via `book_reference_id` — que a validação confirmou estar perfeitamente alinhado entre as 3 bases atuais. Qualquer nova Versão passa por uma checagem de alinhamento de `book_reference_id`/`testament_reference_id` antes de ser habilitada.

## Considered Options

- Casar por nome direto contra `book.name` de cada Versão, com fuzzy matching entre grafias — rejeitado por não escalar entre idiomas/Versões e por já falhar com dados reais.
- Códigos padronizados OSIS/USFM (`GEN`, `1CO`) — rejeitado por ser menos natural para um agente compor a partir de linguagem natural em português.
```

### docs/adr/0002-external-data-directory.md

```markdown
# Diretório de dados externo, descoberto em tempo de execução, com falha rápida na validação

Os arquivos `.sqlite` de cada Versão não são embutidos no binário: o usuário configura um diretório externo, e o servidor faz uma varredura desse diretório na inicialização para descobrir quais Versões estão disponíveis — permitindo adicionar, atualizar ou remover Versões sem recompilar. O identificador de cada Versão (o valor aceito pelo parâmetro `version` das ferramentas) vem do campo `metadata.name` de dentro do próprio arquivo (ex: `"Almeida Revista e Corrigida"`), não do nome do arquivo — o nome do arquivo é um detalhe do sistema de arquivos, enquanto `metadata.name` é o dado autoritativo sobre qual Versão aquele arquivo representa. Cada arquivo descoberto passa pela validação de alinhamento estrutural definida na ADR 0001; se qualquer arquivo falhar essa validação, o servidor **falha ao iniciar** (não apenas exclui aquele arquivo) — um dado corrompido ou incompatível deve ser detectado no deploy, não silenciosamente reduzir as Versões disponíveis em produção.

## Considered Options

- Embutir os arquivos no binário via `go:embed` — rejeitado porque o usuário quer trocar/atualizar Versões sem recompilar.
- Identificador de Versão a partir do nome do arquivo — rejeitado em favor de `metadata.name`, que é o dado autoritativo dentro do próprio arquivo.
- Excluir apenas o arquivo inválido e continuar com os demais (degradação graciosa) — rejeitado em favor de falha rápida no startup.
```

### docs/BACKLOG.md

```markdown
# Backlog — features futuras (fora do escopo da primeira implementação)

Itens identificados durante a sessão de design inicial, deliberadamente adiados.

- **Busca por texto (FTS5)**: ferramenta para buscar versículos por conteúdo textual (ex: "encontre versículos sobre 'amor'"), usando a extensão FTS5 do SQLite. Ver `docs/adr/` quando essa decisão for tomada — adiada porque não há requisito de busca textual na v1; adicionar depois é aditivo (nova tabela virtual, não exige mudar o schema de `verse`).
- **Formato de resposta configurável**: a v1 fixa o formato de saída como "quote" (texto entre aspas + linha de atribuição `- Livro Capítulo:Versículo`). Uma versão futura deve permitir configurar esse formato (ex: estruturado/JSON puro, texto plano sem atribuição, etc.) via parâmetro da ferramenta ou configuração do servidor.
```

## Domain Concept Identification

#### Existing Concepts (from codebase)

- **Versão** (schema `metadata`/`book`/`verse` em `data/*.sqlite`): edição do texto bíblico. Hoje há 3 arquivos na pasta `data/` — `ARA.sqlite` (`metadata.name = "Almeida Revista e Atualizada"`, 66 livros, 31.104 versículos), `ARC.sqlite` (`"Almeida Revista e Corrigida"`, 66 livros, 31.105 versículos) e `NVI.sqlite` (`"Nova Versão Internacional"`, 66 livros, 31.087 versículos) — confirmado por inspeção direta dos 3 arquivos. Todas têm exatamente 66 registros em `book`, com `book_reference_id`/`testament_reference_id` idênticos posição a posição entre as 3 (0 divergências verificadas), única variação sendo a grafia de `book.name` (ex: `book_reference_id=22` → "Cântico", "Cantares", "Cântico dos Cânticos" nas 3 bases, exatamente como a ADR-0001 descreve).
- **Livro**: uma das 66 unidades canônicas, com Testamento e `book_reference_id` estável entre Versões — dado estrutural já existente e validado no schema `book`.
- **Versículo**: unidade atômica em `verse` (book_id, chapter, verse, text) — texto armazenado com espaçamento/pontuação originais, não normalizado.
- **Nome Canônico**: conceito ainda não persistido em nenhum lugar do repositório (não há enum/lista de 66 nomes em código) — é uma decisão de design já fechada na ADR-0001, mas sem artefato de implementação.

#### New Concepts Required

- **Catálogo canônico de Livros**: uma fonte única e fixa dos 66 Nomes Canônicos em português, com seu `book_reference_id`/Testamento, independente de qualquer arquivo `.sqlite` — hoje inexistente; toda a estratégia de identidade de Livro (ADR-0001) depende dela existir.
- **Descoberta e validação de Versões**: mecanismo que varre um diretório configurável, abre cada `.sqlite`, lê `metadata.name` e valida alinhamento contra o Catálogo canônico antes de disponibilizar aquela Versão — hoje inexistente (nenhum código de acesso a dados existe; ver "Estado atual do projeto" em `CLAUDE.md`).
- **Passagem** como abstração de consulta: unificar Versículo único, intervalo e Capítulo inteiro sob uma única operação de consulta (conforme a definição em `CONTEXT.md`, os três são casos do mesmo conceito) — ainda não implementado.
- **Camada de exposição ao agente (ferramentas MCP)**: superfície de ferramentas que aceita `book` (enum fixo) e `version` (enum dinâmico) e devolve o texto puro (sem metadata de copyright) — hoje inexistente; nem a linguagem/runtime/SDK do servidor MCP foram escolhidos ainda.

#### Key Business Rules

- **Identidade de Livro é fixa e não-fuzzy**: o parâmetro `book` só aceita um dos 66 Nomes Canônicos definidos pelo servidor; nunca casamento textual contra `book.name` de uma Versão. Governa: Livro, Nome Canônico, Catálogo canônico.
- **Identidade de Versão vem de `metadata.name`, não do nome do arquivo**: o enum `version` é construído em tempo de execução a partir do conteúdo de cada `.sqlite` descoberto. Governa: Versão, Descoberta e validação de Versões.
- **Falha rápida e total na validação de alinhamento**: se qualquer arquivo descoberto não corresponder ao Catálogo canônico (66 livros, mesma ordem de `book_reference_id`/`testament_reference_id`), nenhuma Versão é servida — não há degradação parcial. Governa: Descoberta e validação de Versões.
- **Passagem nunca cruza Capítulo**: um intervalo de Versículos é sempre relativo a um único (Livro, Capítulo). Governa: Passagem.
- **Extensibilidade sem recompilação**: adicionar/atualizar/remover uma Versão é uma operação de arquivo no diretório de dados, não uma mudança de código. Governa: Versão, Descoberta e validação de Versões.
- **Escopo fechado da v1** (implícita, mas explicitada no Backlog): sem busca textual (FTS5) e sem formato de resposta configurável — a saída é sempre o formato "quote" fixo. Governa: Passagem, camada de exposição ao agente.

## Strategic Approach

#### Solution Direction

O servidor MCP trata cada arquivo `.sqlite` descoberto como uma Versão intercambiável atrás de um modelo canônico único de Livro/Passagem: na inicialização, varre um diretório configurável, valida cada arquivo contra um Catálogo canônico fixo de 66 Livros e expõe as Versões válidas (identificadas por `metadata.name`) como um enum dinâmico. As ferramentas oferecidas ao agente resolvem sempre (Nome Canônico, Versão, Capítulo, intervalo de Versículo) para uma única operação de consulta de Passagem contra o arquivo correspondente, devolvendo apenas o texto do(s) Versículo(s) — nenhuma metadata de copyright é repassada ao agente (constraint já fixada no `CLAUDE.md` do projeto, fora do escopo dos documentos analisados aqui, mas relevante para a Reasons Canvas). Fluxo geral: descoberta/validação de Versões (startup) → resolução de Livro via Catálogo canônico + resolução de Versão via enum descoberto → consulta de Passagem no arquivo da Versão escolhida → formatação de saída → resposta à ferramenta MCP.

#### Key Design Decisions

- **Identidade de Livro**: enum fixo de Nomes Canônicos definido pelo servidor, resolvido via `book_reference_id`, vs. casamento direto/fuzzy contra `book.name` → **recomendação: enum fixo** (já decidido na ADR-0001) — trade-off é exigir manutenção manual do Catálogo canônico ao suportar novos idiomas/Versões, mas elimina ambiguidade de grafia hoje já observada nos dados reais.
- **Descoberta de Versão**: diretório externo varrido em runtime, identificado por `metadata.name`, vs. embutir os arquivos no binário → **recomendação: diretório externo** (já decidido na ADR-0002) — trade-off é exigir uma etapa de validação de alinhamento em toda inicialização (custo de startup), em troca de trocar/atualizar Versões sem recompilar.
- **Política de falha na validação**: falha total do servidor se qualquer arquivo descoberto for inválido, vs. excluir apenas o arquivo inválido e continuar com os demais → **recomendação: falha total** (já decidido na ADR-0002) — trade-off é reduzir disponibilidade (nenhuma Versão é servida mesmo havendo Versões válidas) em troca de nunca mascarar um dado corrompido em produção.
- **Modelo de consulta**: tratar Versículo único e Capítulo inteiro como casos degenerados de Passagem (um único conceito de consulta com limites opcionais) vs. três operações de consulta independentes → **recomendação: unificar sob Passagem** conforme `CONTEXT.md` já define o conceito dessa forma — simplifica a camada de domínio; a divisão em ferramentas MCP distintas (se houver) é uma decisão de superfície de API, a ser resolvida na Reasons Canvas.

#### Alternatives Considered

- Casar Livro por nome direto contra `book.name`, com fuzzy matching entre grafias — rejeitada por não escalar entre idiomas/Versões e por já falhar com os dados reais (3 grafias distintas para o mesmo `book_reference_id=22`).
- Identificar Livro por códigos padronizados OSIS/USFM (`GEN`, `1CO`) — rejeitada por ser menos natural para um agente compor a partir de linguagem natural em português.
- Embutir os arquivos `.sqlite` no binário via `go:embed` (ou equivalente) — rejeitada porque o requisito é trocar/atualizar Versões sem recompilar.
- Identificar Versão pelo nome do arquivo em vez de `metadata.name` — rejeitada porque o nome do arquivo é um detalhe de sistema de arquivos, não um dado autoritativo.
- Excluir apenas o arquivo inválido na validação e continuar com os demais (degradação graciosa) — rejeitada em favor de falha rápida e total no startup.

## Risk & Gap Analysis

#### Requirement Ambiguities

- **Versão padrão não especificada**: nenhum dos 4 documentos analisados define qual Versão é usada quando o parâmetro `version` é omitido pelo agente. É uma decisão pendente para a Reasons Canvas.
- **Colisão de `metadata.name`**: nenhum documento trata o caso de dois arquivos `.sqlite` descobertos no mesmo diretório reportando o mesmo `metadata.name` — o enum de Versão ficaria ambíguo. Precisa de uma regra explícita (ex: falhar a inicialização, como já ocorre para desalinhamento estrutural).
- **Superfície exata de ferramentas MCP não definida aqui**: `CONTEXT.md` descreve o conceito unificado de Passagem (Versículo/Passagem/Capítulo como casos de um só conceito), mas nenhum dos 4 documentos define quantas ferramentas MCP existem, seus nomes ou parâmetros exatos — correto, pois isso é explicitamente escopo da Reasons Canvas, não desta análise.
- **Nenhum documento de negócio com Acceptance Criteria formais existe** — os 4 arquivos são artefatos de design (linguagem ubíqua, ADRs, backlog), não uma história de usuário com ACs explícitos. A tabela de cobertura abaixo usa regras extraídas dos próprios documentos como ACs inferidos (ver nota na seção seguinte).
- **`CLAUDE.md` do projeto está desatualizado em relação ao diretório `data/`**: a seção "Estado atual do projeto" do `CLAUDE.md` afirma que só existem `data/ARA.sqlite` e `data/ARC.sqlite`, mas a inspeção direta do diretório mostra que `data/NVI.sqlite` também já existe e está estruturalmente alinhado (0 mismatches de `book_reference_id`/`testament_reference_id` contra as outras duas bases) — `CONTEXT.md`, por outro lado, já cita as 3 Versões (ARA/ARC/NVI) corretamente. Recomenda-se atualizar `CLAUDE.md` para refletir as 3 Versões atualmente disponíveis.

#### Edge Cases

- Novo arquivo `.sqlite` com 66 registros em `book`, mas em ordem/valores de `book_reference_id`/`testament_reference_id` diferente do Catálogo canônico — deve reprovar a validação e impedir a inicialização do servidor por completo (ADR-0002), não apenas invalidar aquele arquivo.
- Consulta de Passagem com início de intervalo maior que o fim, ou com Capítulo/Versículo inexistente naquela Versão — precisa de uma resposta de erro explícita ao agente, nunca uma resposta de sucesso "vazia" ou parcial disfarçada de sucesso.
- Consulta de Passagem cujo intervalo, se interpretado ingenuamente, cruzaria dois Capítulos — deve ser rejeitada, pois `CONTEXT.md` define Passagem como sempre relativa a um único Capítulo.
- Diretório de dados vazio ou inacessível na inicialização — comportamento não especificado nos 4 documentos (a ADR-0002 só cobre o caso de arquivo inválido, não a ausência total de arquivos); precisa de decisão explícita na Reasons Canvas.
- Diferença de versificação entre Versões (ARA 31.104 vs. ARC 31.105 vs. NVI 31.087 versículos, confirmado por inspeção direta) significa que a mesma referência (Livro, Capítulo, Versículo) pode existir em uma Versão e não em outra — qualquer camada de comparação entre Versões precisa tratar isso como esperado, não como erro de dados.

#### Technical Risks

- **Nenhuma implementação existe ainda**: não há `go.mod`, `package.json` ou qualquer manifesto de build no repositório — linguagem, runtime e SDK de MCP ainda não foram escolhidos. Isso é esperado nesta fase (conforme o próprio `CLAUDE.md`), mas é a maior fonte de risco de estimativa até a Reasons Canvas fixar a Structure/Approach.
- **Validação de alinhamento estrutural deve ser recorrente, não assumida**: hoje as 3 Versões existentes passam a validação (confirmado por inspeção: 0 mismatches de `book_reference_id`/`testament_reference_id`), mas a ADR-0002 exige que essa checagem ocorra a cada novo arquivo descoberto — a lógica de validação não pode ser um teste único, precisa rodar em toda inicialização do servidor.
- **Não normalizar o texto do Versículo**: os documentos analisados não repetem essa regra explicitamente (ela está em `CLAUDE.md`, fora do escopo desta análise), mas qualquer camada de consulta/formatação desenhada na Reasons Canvas precisa preservar espaços/pontuação originais de `verse.text` sem `trim`.
- **Acentuação/diacríticos nos Nomes Canônicos**: nomes como "Gênesis", "Êxodo", "Jó" contêm caracteres acentuados — qualquer mecanismo de enum exposto ao agente (seja qual for o SDK/linguagem escolhido) precisa preservar UTF-8 sem normalização, já que a ADR-0001 rejeita justamente a variação de grafia como fonte de ambiguidade.
- **Acesso concorrente aos arquivos `.sqlite`**: se o servidor atender múltiplas chamadas de ferramenta simultâneas, o acesso (idealmente somente leitura) aos mesmos arquivos por múltiplas Versões precisa ser seguro para concorrência — não tratado nos 4 documentos, a decidir na Reasons Canvas.

#### Acceptance Criteria Coverage

> Nenhum dos 4 documentos contém uma lista formal de Acceptance Criteria (são artefatos de design, não uma história de usuário). Os itens abaixo foram **inferidos** das regras explícitas nos documentos, para permitir uma avaliação de cobertura; devem ser revisados/formalizados na Reasons Canvas.

| AC# | Description | Addressable? | Gaps/Notes |
|-----|-------------|--------------|------------|
| 1 | `book` é restrito a um enum fixo dos 66 Nomes Canônicos, nunca por fuzzy matching contra `book.name` de uma Versão (ADR-0001) | Yes | Requer que o Catálogo canônico canônico seja implementado antes de qualquer ferramenta MCP |
| 2 | `version` é um enum construído dinamicamente a partir de `metadata.name` de cada `.sqlite` válido descoberto (ADR-0002) | Yes | Regra de colisão de `metadata.name` duplicado não definida — gap |
| 3 | Falha de alinhamento estrutural em qualquer arquivo descoberto impede a inicialização do servidor por completo (ADR-0002) | Yes | Comportamento para diretório vazio/inacessível não coberto — gap |
| 4 | Uma Passagem nunca cruza Capítulos; Versículo único e Capítulo inteiro são casos degenerados de Passagem (CONTEXT.md) | Yes | Nomes/quantidade de ferramentas MCP que expõem esses casos ainda não definidos — correto, fica para a Reasons Canvas |
| 5 | Adicionar/atualizar/remover uma Versão não exige recompilar o servidor (ADR-0002) | Yes | Nenhum gap identificado nos documentos analisados |
| 6 | Busca por texto (FTS5) está fora do escopo da v1 (BACKLOG.md) | Yes | Nenhum gap — item explicitamente adiado |
| 7 | Formato de resposta é fixo ("quote") na v1; configurabilidade é adiada (BACKLOG.md) | Partial | O formato exato da string "quote" não é especificado nestes 4 documentos — detalhe de formatação cabe à Reasons Canvas |

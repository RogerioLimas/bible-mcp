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

# Identidade de Livro canônica e independente de Versão

As ferramentas do MCP precisam identificar um Livro para qualquer Versão (ARA, ARC, NVI, e futuras). Casar diretamente com a coluna `book.name` de cada base foi rejeitado: validação nos dados reais mostrou que o livro `book_reference_id=22` já tem 3 grafias diferentes entre as Versões atuais ("Cântico", "Cantares", "Cântico dos Cânticos"), e isso só cresceria com novas Versões/idiomas. Em vez disso, o parâmetro `book` exposto ao agente é um enum fixo de Nomes Canônicos em português (definidos pelo servidor, não lidos de nenhuma base), resolvido internamente para o registro de cada Versão via `book_reference_id` — que a validação confirmou estar perfeitamente alinhado entre as 3 bases atuais. Qualquer nova Versão passa por uma checagem de alinhamento de `book_reference_id`/`testament_reference_id` antes de ser habilitada.

## Considered Options

- Casar por nome direto contra `book.name` de cada Versão, com fuzzy matching entre grafias — rejeitado por não escalar entre idiomas/Versões e por já falhar com dados reais.
- Códigos padronizados OSIS/USFM (`GEN`, `1CO`) — rejeitado por ser menos natural para um agente compor a partir de linguagem natural em português.

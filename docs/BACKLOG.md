# Backlog — features futuras (fora do escopo da primeira implementação)

Itens identificados durante a sessão de design inicial, deliberadamente adiados.

- **Busca por texto (FTS5)**: ferramenta para buscar versículos por conteúdo textual (ex: "encontre versículos sobre 'amor'"), usando a extensão FTS5 do SQLite. Ver `docs/adr/` quando essa decisão for tomada — adiada porque não há requisito de busca textual na v1; adicionar depois é aditivo (nova tabela virtual, não exige mudar o schema de `verse`).
- **Formato de resposta configurável**: a v1 fixa o formato de saída como "quote" (texto entre aspas + linha de atribuição `- Livro Capítulo:Versículo`). Uma versão futura deve permitir configurar esse formato (ex: estruturado/JSON puro, texto plano sem atribuição, etc.) via parâmetro da ferramenta ou configuração do servidor.

# Diretório de dados externo, descoberto em tempo de execução, com falha rápida na validação

Os arquivos `.sqlite` de cada Versão não são embutidos no binário: o usuário configura um diretório externo, e o servidor faz uma varredura desse diretório na inicialização para descobrir quais Versões estão disponíveis — permitindo adicionar, atualizar ou remover Versões sem recompilar. O identificador de cada Versão (o valor aceito pelo parâmetro `version` das ferramentas) vem do campo `metadata.name` de dentro do próprio arquivo (ex: `"Almeida Revista e Corrigida"`), não do nome do arquivo — o nome do arquivo é um detalhe do sistema de arquivos, enquanto `metadata.name` é o dado autoritativo sobre qual Versão aquele arquivo representa. Cada arquivo descoberto passa pela validação de alinhamento estrutural definida na ADR 0001; se qualquer arquivo falhar essa validação, o servidor **falha ao iniciar** (não apenas exclui aquele arquivo) — um dado corrompido ou incompatível deve ser detectado no deploy, não silenciosamente reduzir as Versões disponíveis em produção.

## Considered Options

- Embutir os arquivos no binário via `go:embed` — rejeitado porque o usuário quer trocar/atualizar Versões sem recompilar.
- Identificador de Versão a partir do nome do arquivo — rejeitado em favor de `metadata.name`, que é o dado autoritativo dentro do próprio arquivo.
- Excluir apenas o arquivo inválido e continuar com os demais (degradação graciosa) — rejeitado em favor de falha rápida no startup.

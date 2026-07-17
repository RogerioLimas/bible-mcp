package bible

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatQuote renderiza os versículos como um bloco de citação: uma linha
// numerada por versículo, seguida de uma linha de atribuição com a
// referência livro/capítulo/versículo.
func FormatQuote(bookName string, chapter int, verses []VerseText) string {
	lines := make([]string, 0, len(verses)+1)
	for _, v := range verses {
		lines = append(lines, strconv.Itoa(v.Number)+" "+v.Text)
	}
	lines = append(lines, "- "+Reference(bookName, chapter, verses))
	return strings.Join(lines, "\n")
}

// Reference renderiza a linha de atribuição para um conjunto de
// versículos, ex: "Gênesis 1:1" para um único versículo ou
// "Gênesis 1:1-2" para um intervalo. Assume len(verses) >= 1 — Store.Verses
// nunca retorna uma lista vazia sem erro.
func Reference(bookName string, chapter int, verses []VerseText) string {
	first, last := verses[0].Number, verses[len(verses)-1].Number
	if first == last {
		return fmt.Sprintf("%s %d:%d", bookName, chapter, first)
	}
	return fmt.Sprintf("%s %d:%d-%d", bookName, chapter, first, last)
}

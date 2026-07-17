package bible

import "testing"

func TestFormatQuoteVersiculoUnico(t *testing.T) {
	got := FormatQuote("Gênesis", 1, []VerseText{
		{Number: 1, Text: "No princípio, criou Deus os céus e a terra. "},
	})
	want := "1 No princípio, criou Deus os céus e a terra. \n- Gênesis 1:1"
	if got != want {
		t.Fatalf("FormatQuote() = %q, want %q", got, want)
	}
}

func TestFormatQuotePassagem(t *testing.T) {
	got := FormatQuote("Gênesis", 1, []VerseText{
		{Number: 1, Text: "No princípio, criou Deus os céus e a terra. "},
		{Number: 2, Text: "E a terra era sem forma e vazia; e havia trevas sobre a face do abismo; e o Espírito de Deus se movia sobre a face das águas."},
	})
	want := "1 No princípio, criou Deus os céus e a terra. \n2 E a terra era sem forma e vazia; e havia trevas sobre a face do abismo; e o Espírito de Deus se movia sobre a face das águas.\n- Gênesis 1:1-2"
	if got != want {
		t.Fatalf("FormatQuote() = %q, want %q", got, want)
	}
}

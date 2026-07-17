package bible

import (
	"reflect"
	"testing"
)

func openTestStore(t *testing.T) *Store {
	t.Helper()
	store, err := Open(testDataDir)
	if err != nil {
		t.Fatalf("Open(%q) = %v", testDataDir, err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestVersesVersiculoUnico(t *testing.T) {
	store := openTestStore(t)
	genesis, _ := BookByName("Gênesis")

	got, err := store.Verses(DefaultVersionName, genesis, 1, 1, 1)
	if err != nil {
		t.Fatalf("Verses(...) = %v", err)
	}
	want := []VerseText{{Number: 1, Text: "No princípio, criou Deus os céus e a terra. "}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Verses(Gênesis 1:1) = %+v, want %+v", got, want)
	}
}

func TestVersesPassagem(t *testing.T) {
	store := openTestStore(t)
	genesis, _ := BookByName("Gênesis")

	got, err := store.Verses(DefaultVersionName, genesis, 1, 1, 2)
	if err != nil {
		t.Fatalf("Verses(...) = %v", err)
	}
	want := []VerseText{
		{Number: 1, Text: "No princípio, criou Deus os céus e a terra. "},
		{Number: 2, Text: "E a terra era sem forma e vazia; e havia trevas sobre a face do abismo; e o Espírito de Deus se movia sobre a face das águas."},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Verses(Gênesis 1:1-2) = %+v, want %+v", got, want)
	}
}

func TestVersesCapituloInteiro(t *testing.T) {
	store := openTestStore(t)
	genesis, _ := BookByName("Gênesis")

	got, err := store.Verses(DefaultVersionName, genesis, 1, 0, 0)
	if err != nil {
		t.Fatalf("Verses(...) = %v", err)
	}
	if len(got) != 31 {
		t.Fatalf("len(Verses(Gênesis 1)) = %d, want 31", len(got))
	}
	if got[0].Number != 1 || got[len(got)-1].Number != 31 {
		t.Fatalf("Verses(Gênesis 1) primeiro/último número = %d/%d, want 1/31", got[0].Number, got[len(got)-1].Number)
	}
}

func TestVersesIntervaloInvalido(t *testing.T) {
	store := openTestStore(t)
	genesis, _ := BookByName("Gênesis")

	if _, err := store.Verses(DefaultVersionName, genesis, 1, 5, 2); err == nil {
		t.Fatal("Verses(verse_start > verse_end) não retornou erro")
	}
}

func TestVersesReferenciaInexistente(t *testing.T) {
	store := openTestStore(t)
	genesis, _ := BookByName("Gênesis")

	if _, err := store.Verses(DefaultVersionName, genesis, 999, 1, 1); err == nil {
		t.Fatal("Verses(capítulo inexistente) não retornou erro")
	}
}

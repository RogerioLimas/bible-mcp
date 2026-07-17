package bible

import "testing"

func TestBookByName(t *testing.T) {
	book, ok := BookByName("Gênesis")
	if !ok {
		t.Fatalf("BookByName(%q) reportou não encontrado", "Gênesis")
	}
	if book.ReferenceID != 1 || book.Testament != OldTestament {
		t.Fatalf("BookByName(%q) = %+v, want ReferenceID=1 Testament=OldTestament", "Gênesis", book)
	}
}

func TestBookByNameDesconhecido(t *testing.T) {
	if _, ok := BookByName("Livro Inexistente"); ok {
		t.Fatal("BookByName(livro inexistente) reportou encontrado")
	}
}

func TestCanonTem66Livros(t *testing.T) {
	if len(Canon) != 66 {
		t.Fatalf("len(Canon) = %d, want 66", len(Canon))
	}
	if len(Names()) != 66 {
		t.Fatalf("len(Names()) = %d, want 66", len(Names()))
	}
	last := Canon[len(Canon)-1]
	if last.Name != "Apocalipse" || last.Testament != NewTestament {
		t.Fatalf("último livro = %+v, want Apocalipse/NewTestament", last)
	}
}

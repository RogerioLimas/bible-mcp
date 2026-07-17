package bible

import (
	"os"
	"path/filepath"
	"testing"
)

const testDataDir = "../../data"

func TestOpenDescobreTodasAsVersoes(t *testing.T) {
	store, err := Open(testDataDir)
	if err != nil {
		t.Fatalf("Open(%q) = %v", testDataDir, err)
	}
	defer store.Close()

	got := store.VersionNames()
	want := []string{"Almeida Revista e Atualizada", "Almeida Revista e Corrigida", "Nova Versão Internacional"}
	if len(got) != len(want) {
		t.Fatalf("VersionNames() = %v, want %v", got, want)
	}
	for i, name := range want {
		if got[i] != name {
			t.Fatalf("VersionNames()[%d] = %q, want %q", i, got[i], name)
		}
	}
}

func TestOpenRejeitaDiretorioSemArquivosSqlite(t *testing.T) {
	if _, err := Open(t.TempDir()); err == nil {
		t.Fatal("Open(diretório vazio) não retornou erro")
	}
}

func copyFile(t *testing.T, dst, src string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("lendo %s: %v", src, err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		t.Fatalf("escrevendo %s: %v", dst, err)
	}
}

func TestOpenRejeitaMetadataNameDuplicado(t *testing.T) {
	dir := t.TempDir()
	// Duas cópias do mesmo arquivo real => mesmo metadata.name duplicado,
	// sem inventar dados sintéticos.
	copyFile(t, filepath.Join(dir, "a.sqlite"), filepath.Join(testDataDir, "ARC.sqlite"))
	copyFile(t, filepath.Join(dir, "b.sqlite"), filepath.Join(testDataDir, "ARC.sqlite"))

	if _, err := Open(dir); err == nil {
		t.Fatal("Open(duas Versões com metadata.name duplicado) não retornou erro")
	}
}

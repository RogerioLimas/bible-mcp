package bible

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sort"

	_ "modernc.org/sqlite"
)

// DefaultVersionName é a Versão usada quando uma chamada de ferramenta omite
// o parâmetro version. Ver docs/adr/0002-external-data-directory.md.
const DefaultVersionName = "Almeida Revista e Corrigida"

// VerseText é o número e o texto de um versículo, como armazenado por uma
// Versão.
type VerseText struct {
	Number int
	Text   string
}

type version struct {
	name        string
	db          *sql.DB
	bookIDByRef map[int]int // ReferenceID -> book.id, nesta base específica
}

// Store mantém cada Versão descoberta em um diretório de dados, indexada
// pelo seu metadata.name.
type Store struct {
	versions map[string]*version
	names    []string // em ordem alfabética, para uma listagem de enum estável
}

// Open varre dir por arquivos *.sqlite, abre cada um em modo somente
// leitura e valida se a tabela book corresponde exatamente ao Canon. Se
// qualquer arquivo falhar na validação, ou se duas Versões reportarem o
// mesmo metadata.name, Open falha e nenhum Store é retornado
// (docs/adr/0002-external-data-directory.md).
func Open(dir string) (*Store, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.sqlite"))
	if err != nil {
		return nil, fmt.Errorf("listando %s: %w", dir, err)
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("nenhum arquivo .sqlite encontrado em %s", dir)
	}

	store := &Store{versions: make(map[string]*version)}
	for _, path := range paths {
		v, err := openVersion(path)
		if err != nil {
			store.Close()
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if _, exists := store.versions[v.name]; exists {
			store.Close()
			return nil, fmt.Errorf("%s: Versão %q duplicada", path, v.name)
		}
		store.versions[v.name] = v
		store.names = append(store.names, v.name)
	}
	sort.Strings(store.names)
	return store, nil
}

func openVersion(path string) (*version, error) {
	dsn := "file:" + path + "?mode=ro"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrindo: %w", err)
	}
	name, err := metadataName(db)
	if err != nil {
		db.Close()
		return nil, err
	}
	bookIDByRef, err := validateBooks(db)
	if err != nil {
		db.Close()
		return nil, err
	}
	return &version{name: name, db: db, bookIDByRef: bookIDByRef}, nil
}

func metadataName(db *sql.DB) (string, error) {
	var name string
	err := db.QueryRow(`SELECT value FROM metadata WHERE key = 'name'`).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("lendo metadata.name: %w", err)
	}
	return name, nil
}

func validateBooks(db *sql.DB) (map[int]int, error) {
	rows, err := db.Query(`SELECT id, book_reference_id, testament_reference_id FROM book ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("lendo book: %w", err)
	}
	defer rows.Close()

	bookIDByRef := make(map[int]int, len(Canon))
	count := 0
	for rows.Next() {
		var id, refID, testamentID int
		if err := rows.Scan(&id, &refID, &testamentID); err != nil {
			return nil, fmt.Errorf("lendo linha de book: %w", err)
		}
		if count >= len(Canon) {
			return nil, fmt.Errorf("mais de %d livros encontrados", len(Canon))
		}
		want := Canon[count]
		if refID != want.ReferenceID || Testament(testamentID) != want.Testament {
			return nil, fmt.Errorf(
				"livro na posição %d (id=%d) não corresponde ao cânone: got book_reference_id=%d testament_reference_id=%d, want book_reference_id=%d testament_reference_id=%d",
				count, id, refID, testamentID, want.ReferenceID, want.Testament)
		}
		bookIDByRef[want.ReferenceID] = id
		count++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if count != len(Canon) {
		return nil, fmt.Errorf("%d livros encontrados, esperado %d", count, len(Canon))
	}
	return bookIDByRef, nil
}

// Close fecha todas as bases de Versão abertas.
func (s *Store) Close() error {
	var firstErr error
	for _, v := range s.versions {
		if err := v.db.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// VersionNames retorna o metadata.name de cada Versão descoberta, em ordem
// alfabética.
func (s *Store) VersionNames() []string {
	return s.names
}

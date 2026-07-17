package bible

import "fmt"

// Verses retorna os versículos de book/chapter em uma Versão. Se
// verseStart e verseEnd forem ambos zero, todos os versículos do capítulo
// são retornados (um Capítulo). Caso contrário, apenas os versículos em
// [verseStart, verseEnd] são retornados (uma Passagem, ou um único
// Versículo quando verseStart == verseEnd). A consulta é sempre restrita a
// um único (book_id, chapter) — nunca cruza Capítulos.
func (s *Store) Verses(versionName string, book Book, chapter, verseStart, verseEnd int) ([]VerseText, error) {
	if verseStart != 0 && verseStart > verseEnd {
		return nil, fmt.Errorf("verse_start (%d) maior que verse_end (%d)", verseStart, verseEnd)
	}

	v, ok := s.versions[versionName]
	if !ok {
		return nil, fmt.Errorf("versão desconhecida: %q", versionName)
	}
	bookID, ok := v.bookIDByRef[book.ReferenceID]
	if !ok {
		return nil, fmt.Errorf("livro %q não encontrado na versão %q", book.Name, versionName)
	}

	query := `SELECT verse, text FROM verse WHERE book_id = ? AND chapter = ?`
	args := []any{bookID, chapter}
	if verseStart != 0 {
		query += ` AND verse BETWEEN ? AND ?`
		args = append(args, verseStart, verseEnd)
	}
	query += ` ORDER BY verse`

	rows, err := v.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("consultando versículos: %w", err)
	}
	defer rows.Close()

	var verses []VerseText
	for rows.Next() {
		var vt VerseText
		if err := rows.Scan(&vt.Number, &vt.Text); err != nil {
			return nil, fmt.Errorf("lendo versículo: %w", err)
		}
		verses = append(verses, vt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(verses) == 0 {
		return nil, fmt.Errorf("nenhum versículo encontrado para %s %d:%d-%d na versão %q", book.Name, chapter, verseStart, verseEnd, versionName)
	}
	return verses, nil
}

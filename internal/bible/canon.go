package bible

// Testament identifica a qual Testamento um Livro pertence.
type Testament int

const (
	OldTestament Testament = 1
	NewTestament Testament = 2
)

// Book é um dos 66 livros canônicos da Bíblia, identificado
// independentemente de qualquer Versão específica. Ver CONTEXT.md para as
// definições de Livro e Nome Canônico.
type Book struct {
	ReferenceID int
	Name        string
	Testament   Testament
}

// Canon é a lista fixa e ordenada dos 66 livros canônicos.
var Canon = []Book{
	{1, "Gênesis", OldTestament},
	{2, "Êxodo", OldTestament},
	{3, "Levítico", OldTestament},
	{4, "Números", OldTestament},
	{5, "Deuteronômio", OldTestament},
	{6, "Josué", OldTestament},
	{7, "Juízes", OldTestament},
	{8, "Rute", OldTestament},
	{9, "1Samuel", OldTestament},
	{10, "2Samuel", OldTestament},
	{11, "1Reis", OldTestament},
	{12, "2Reis", OldTestament},
	{13, "1Crônicas", OldTestament},
	{14, "2Crônicas", OldTestament},
	{15, "Esdras", OldTestament},
	{16, "Neemias", OldTestament},
	{17, "Ester", OldTestament},
	{18, "Jó", OldTestament},
	{19, "Salmos", OldTestament},
	{20, "Provérbios", OldTestament},
	{21, "Eclesiastes", OldTestament},
	{22, "Cântico dos Cânticos", OldTestament},
	{23, "Isaías", OldTestament},
	{24, "Jeremias", OldTestament},
	{25, "Lamentações", OldTestament},
	{26, "Ezequiel", OldTestament},
	{27, "Daniel", OldTestament},
	{28, "Oseias", OldTestament},
	{29, "Joel", OldTestament},
	{30, "Amós", OldTestament},
	{31, "Obadias", OldTestament},
	{32, "Jonas", OldTestament},
	{33, "Miqueias", OldTestament},
	{34, "Naum", OldTestament},
	{35, "Habacuque", OldTestament},
	{36, "Sofonias", OldTestament},
	{37, "Ageu", OldTestament},
	{38, "Zacarias", OldTestament},
	{39, "Malaquias", OldTestament},
	{40, "Mateus", NewTestament},
	{41, "Marcos", NewTestament},
	{42, "Lucas", NewTestament},
	{43, "João", NewTestament},
	{44, "Atos", NewTestament},
	{45, "Romanos", NewTestament},
	{46, "1Coríntios", NewTestament},
	{47, "2Coríntios", NewTestament},
	{48, "Gálatas", NewTestament},
	{49, "Efésios", NewTestament},
	{50, "Filipenses", NewTestament},
	{51, "Colossenses", NewTestament},
	{52, "1Tessalonicenses", NewTestament},
	{53, "2Tessalonicenses", NewTestament},
	{54, "1Timóteo", NewTestament},
	{55, "2Timóteo", NewTestament},
	{56, "Tito", NewTestament},
	{57, "Filemom", NewTestament},
	{58, "Hebreus", NewTestament},
	{59, "Tiago", NewTestament},
	{60, "1Pedro", NewTestament},
	{61, "2Pedro", NewTestament},
	{62, "1João", NewTestament},
	{63, "2João", NewTestament},
	{64, "3João", NewTestament},
	{65, "Judas", NewTestament},
	{66, "Apocalipse", NewTestament},
}

// BookByName resolve um Nome Canônico para o Book correspondente,
// retornando ok=false se name não for um dos 66 nomes canônicos.
func BookByName(name string) (Book, bool) {
	for _, b := range Canon {
		if b.Name == name {
			return b, true
		}
	}
	return Book{}, false
}

// Names retorna os nomes canônicos na ordem do cânone.
func Names() []string {
	names := make([]string, len(Canon))
	for i, b := range Canon {
		names[i] = b.Name
	}
	return names
}

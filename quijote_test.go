package main

import (
	"reflect"
	"strings"
	"testing"
)

var wordList = []string{
	"Hola",
	"adiós",
	"yo",
	"únicos",
}

var wordsExpected = []wordParsed{
	{true, "hola"},
	{true, "adiós"},
	{false, ""},
	{true, "únicos"},
}

var lines = []string{
	"gentes. Nuestra mejor salsa es la hambre; y como ésta no falta a los",
	"-Calla, boba -dijo Sancho-, que todo será usarlo dos o tres años; que.",
	"-Pues guíe vuestra merced -respondió Sancho-: quizá será así; aunque yo lo",
	"-¡Válame Dios! -dijo la sobrina-. ¡Que sepa vuestra merced tanto, señor",
	"emendarme; que yo soy tan fócil...",
	"sin duda alguna. Vale.",
	"outside the United States. U.S. laws alone swamp our small staff.",
	"vuestra fermosura. Y también cuando leía: Dios...los altos cielos que de",
	"le dijo: ''Hermano:",
}

var linesExpected = []lineParsed{
	{false, []string{"nuestra"}},
	{true, nil},
	{false, []string{"quizá"}},
	{false, []string{"que"}},
	{true, nil},
	{true, []string{"vale"}},
	{true, []string{"laws"}},
	{false, []string{"dios", "los"}},
	{true, []string{"hermano"}},
}

var text = `Media noche era por filo, poco más a menos, cuando don Quijote y Sancho
dejaron el monte y entraron en el Toboso. Estaba el pueblo en un sosegado
silencio, porque todos sus vecinos dormían y reposaban a pierna tendida,
como suele decirse. Era la noche entreclara, puesto que quisiera Sancho que
fuera del todo escura, por hallar en su escuridad disculpa de su sandez. No
se oía en todo el lugar sino ladridos de perros, que atronaban los oídos de
don Quijote y turbaban el corazón de Sancho. De cuando en cuando, rebuznaba
un jumento, gruñían puercos, mayaban gatos, cuyas voces, de diferentes
sonidos, se aumentaban con el silencio de la noche, todo lo cual tuvo el
enamorado caballero a mal agüero; pero, con todo esto, dijo a Sancho:

-Sancho, hijo, guía al palacio de Dulcinea: quizá podrá ser que la hallemos
despierta.`

var textExpected = []wordLine{
	{2, []string{"estaba"}},
	{4, []string{"era"}},
	{12, []string{"sancho", "quizá"}},
}

var printExpected = "test:4\n" +
	"\ttestFile:1\n" +
	"\ttestFile:2\n" +
	"\ttestFile:3\n" +
	"\ttestFile:4\n"

func TestParseWord(t *testing.T) {
	for i := 0; i < len(wordList); i++ {
		if got := parseWord(wordList[i]); got != wordsExpected[i] {
			t.Errorf("parseWord(%s) = %v, expected %v", wordList[i], got, wordsExpected[i])
		}
	}
}

func testEq(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func testEqLine(a, b lineParsed) bool {

	if a.scanNext != b.scanNext {
		return false
	}

	if !testEq(a.wordList, b.wordList) {
		return false
	}

	return true
}

func TestParseLine(t *testing.T) {
	for i := 0; i < len(lines); i++ {
		if got := parseLine(false, lines[i]); !testEqLine(got, linesExpected[i]) {
			t.Errorf("parseLine(%s) = %v, expected %v", lines[i], got, linesExpected[i])
		}
	}
}

func TestNewWords(t *testing.T) {
	var words, expected words
	words = newWords()
	expected = make(map[string][]word)

	if reflect.TypeOf(words) != reflect.TypeOf(expected) {
		t.Errorf("not equal typeOf -> expected: %s, got: %s", reflect.TypeOf(expected), reflect.TypeOf(words))
	}
	if reflect.TypeOf(words).Kind() != reflect.TypeOf(expected).Kind() {
		t.Errorf("not equal kind -> expected: %s, got: %s", reflect.TypeOf(expected).Kind(), reflect.TypeOf(words).Kind())
	}
}

func TestAddPrint(t *testing.T) {
	words := newWords()

	words.addWord("test", 1, "testFile")
	words.addWord("test", 2, "testFile")
	words.addWord("test", 3, "testFile")
	words.addWord("test", 4, "testFile")

	if got := words.String(); got != printExpected {
		t.Errorf("words.String() =\n%v\n, expected\n%v", got, printExpected)
	}
}

func testEqFile(a, b []wordLine) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].lineNumber != b[i].lineNumber {
			return false
		}

		if !testEq(a[i].wordList, b[i].wordList) {
			return false
		}
	}

	return true
}

func TestParseFile(t *testing.T) {
	reader := strings.NewReader(text)

	if got := parseFile(reader); !testEqFile(got, textExpected) {
		t.Errorf("parseFile failed: expected -> %v, got -> %v", textExpected, got)
	}
}

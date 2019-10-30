package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"
)

type word struct {
	filename   string
	lineNumber int
}

type words map[string][]word

type wordParsed struct {
	isValid    bool
	wordString string
}

type lineParsed struct {
	scanNext bool
	wordList []string
}

type wordLine struct {
	lineNumber int
	wordList   []string
}

func newWords() words {
	return words(make(map[string][]word))
}

func (w words) String() string {
	var s string
	var keys []string

	for k := range w {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		t := fmt.Sprintf("%s:%d\n", k, len(w[k]))
		s += t

		for _, event := range w[k] {
			t := fmt.Sprintf("\t%s:%d\n", event.filename, event.lineNumber)
			s += t
		}
	}

	return s
}

// Adds the word to the dictionary
func (w words) addWord(wordString string, lineNumber int, filename string) {
	event := word{filename, lineNumber}

	w[wordString] = append(w[wordString], event)
}

func parseWord(word string) wordParsed {
	if len(word) < 3 {
		return wordParsed{false, ""}
	}

	return wordParsed{true, strings.ToLower(word)}
}

func splitLine(r rune) bool {
	return r == ':' || r == '.'
}

// Splits a line and convert every word into lowercase
func parseLine(scanNext bool, line string) lineParsed {
	var words []string

	lineSplitted := strings.FieldsFunc(line, splitLine)
	for _, e := range lineSplitted {
		f := func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
		}
		eFields := strings.FieldsFunc(e, f)

		if scanNext && len(eFields) > 0 {
			result := parseWord(eFields[0])
			if result.isValid {
				words = append(words, result.wordString)
			}
		}
		scanNext = true
	}

	if strings.LastIndexFunc(line, splitLine) != len(line)-1 {
		scanNext = false
	}

	return lineParsed{scanNext, words}
}

// Reads the file line-by-line and discards any non spanish alphabetic character
func parseFile(file io.Reader) []wordLine {
	var words []wordLine
	lineNumber := 1
	scanNext := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		lineParsed := parseLine(scanNext, line)
		scanNext = lineParsed.scanNext
		if lineParsed.wordList != nil {
			words = append(words, wordLine{lineNumber, lineParsed.wordList})
		}

		lineNumber++
	}

	// Prints any error encountered by the Scanner
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading file:", err)
	}

	return words
}

func main() {
	args := os.Args

	if len(args) <= 1 {
		fmt.Fprintf(os.Stderr, "%s requires one or more arguments\n", os.Args[0])
		os.Exit(1)
	}

	words := newWords()

	for i := 1; i < len(args); i++ {
		filename := args[i]
		file, err := os.Open(filename)

		if err != nil {
			fmt.Fprintln(os.Stderr, "opening file:", err)
			return
		}

		wordLines := parseFile(file)
		for _, wordLine := range wordLines {
			for _, word := range wordLine.wordList {
				words.addWord(word, wordLine.lineNumber, filename)
			}
		}

		file.Close()
	}

	fmt.Print(words.String())
}

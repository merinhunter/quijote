package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
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

type wordLine struct {
	lineNumber int
	wordList   []string
}

func newWords() words {
	return words(make(map[string][]word))
}

func (w words) printWords() {
	var keys []string

	for k := range w {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s:%d\n", k, len(w[k]))

		for _, event := range w[k] {
			fmt.Printf("\t%s:%d\n", event.filename, event.lineNumber)
		}
	}
}

// Adds the word to the dictionary
func (w words) addWord(wordString string, lineNumber int, filename string) {
	event := word{filename, lineNumber}

	w[wordString] = append(w[wordString], event)
}

// Splits a line and convert every word into lowercase
func parseLine(line string) []string {
	// Makes a regex to discard non spanish alphabetic characters
	reg, err := regexp.Compile("[^a-zA-ZáéíóúïüñÁÉÍÓÚÏÜÑ.:]+")
	if err != nil {
		fmt.Fprintln(os.Stderr, "compiling regex:", err)
	} // TODO: Is necesary?
	line = reg.ReplaceAllString(line, " ") // Applies the regex

	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)

	var words []string
	scanNext := false
	for scanner.Scan() {
		word := scanner.Text()

		if scanNext {
			result := parseWord(word)
			if result.isValid {
				words = append(words, result.wordString)
			}
			scanNext = false
		}

		if string(word[len(word)-1:]) == "." || string(word[len(word)-1:]) == ":" {
			scanNext = true
		}
	}

	// Prints any error encountered by the Scanner
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading line:", err)
	}

	return words
}

func parseWord(word string) wordParsed {
	// Makes a regex to discard non spanish alphabetic characters
	reg, err := regexp.Compile("[^a-zA-ZáéíóúïüñÁÉÍÓÚÏÜÑ]+")
	if err != nil {
		fmt.Fprintln(os.Stderr, "compiling regex:", err)
	} // TODO: Is necesary?
	word = reg.ReplaceAllString(word, "") // Applies the regex

	if len(word) < 3 {
		return wordParsed{false, ""}
	}

	return wordParsed{true, strings.ToLower(word)}
}

// Reads the file line-by-line and discards any non spanish alphabetic character
func parseFile(file io.Reader) []wordLine {
	scanner := bufio.NewScanner(file)

	var words []wordLine
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()

		wordList := parseLine(line)
		if wordList != nil {
			words = append(words, wordLine{lineNumber, wordList})
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

	words.printWords()
}

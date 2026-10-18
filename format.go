package main

import (
	"strings"
	"unicode"
)

func fixArticles(text string) string {
	words := strings.Fields(text)
	for index := 0; index < len(words)-1; index++ {
		article := cleanTextWord(words[index])
		if !strings.EqualFold(article, "a") {
			continue
		}

		nextWord := ""
		for nextIndex := index + 1; nextIndex < len(words); nextIndex++ {
			nextWord = cleanTextWord(words[nextIndex])
			if nextWord != "" {
				break
			}
		}
		if nextWord == "" {
			continue
		}

		firstCharacter := unicode.ToLower([]rune(nextWord)[0])
		if !strings.ContainsRune("aeiouh", firstCharacter) {
			continue
		}

		replacement := "an"
		if article == "A" {
			replacement = "An"
		}
		position := strings.Index(words[index], article)
		if position >= 0 {
			words[index] = words[index][:position] + replacement + words[index][position+len(article):]
		}
	}
	return strings.Join(words, " ")
}
func cleanTextWord(word string) string {
	return strings.Trim(word, ".,!?:;'\"()[]{}")
}
func fixPunctuation(text string) string {
	characters := []rune(text)
	var output []rune

	for index := 0; index < len(characters); index++ {
		character := characters[index]
		if !strings.ContainsRune(".,!?:;", character) {
			output = append(output, character)
			continue
		}

		if (character == '.' || character == ',') && index > 0 && index+1 < len(characters) && unicode.IsDigit(characters[index-1]) && unicode.IsDigit(characters[index+1]) {
			output = append(output, character)
			continue
		}

		for len(output) > 0 && output[len(output)-1] == ' ' {
			output = output[:len(output)-1]
		}
		output = append(output, character)
		for index+1 < len(characters) && strings.ContainsRune(".,!?:;", characters[index+1]) {
			index++
			output = append(output, characters[index])
		}
		for index+1 < len(characters) && characters[index+1] == ' ' {
			index++
		}
		if index+1 < len(characters) {
			output = append(output, ' ')
		}
	}
	return strings.Join(strings.Fields(string(output)), " ")
}
func fixQuotes(text string) string {
	characters := []rune(text)
	var output []rune
	insideQuote := false

	for index := 0; index < len(characters); index++ {
		if characters[index] != '\'' {
			output = append(output, characters[index])
			continue
		}

		if index > 0 && index+1 < len(characters) && unicode.IsLetter(characters[index-1]) && unicode.IsLetter(characters[index+1]) {
			output = append(output, '\'')
			continue
		}

		if !insideQuote {
			output = append(output, '\'')
			insideQuote = true
			for index+1 < len(characters) && characters[index+1] == ' ' {
				index++
			}
		} else {
			for len(output) > 0 && output[len(output)-1] == ' ' {
				output = output[:len(output)-1]
			}
			output = append(output, '\'')
			insideQuote = false
		}
	}
	return strings.Join(strings.Fields(string(output)), " ")
}
func fixBrackets(text string) string {
	characters := []rune(text)
	var output []rune

	for index := 0; index < len(characters); index++ {
		if characters[index] == '(' {
			output = append(output, '(')
			for index+1 < len(characters) && characters[index+1] == ' ' {
				index++
			}
			continue
		}
		if characters[index] == ')' {
			for len(output) > 0 && output[len(output)-1] == ' ' {
				output = output[:len(output)-1]
			}
			output = append(output, ')')
			continue
		}
		output = append(output, characters[index])
	}
	return strings.Join(strings.Fields(string(output)), " ")
}
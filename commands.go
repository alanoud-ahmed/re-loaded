package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)

func applyCommands(text string, lineNumber int) (string, string) {
	words := strings.Fields(text)
	var output []string

	for _, word := range words {
		command := strings.ToLower(word)

		if command == "(hex)" || command == "(bin)" {
			base := 16
			numberType := "hexadecimal"
			if command == "(bin)" {
				base = 2
				numberType = "binary"
			}

			targetIndex := lastTarget(output)
			if targetIndex == -1 {
				return "", fmt.Sprintf("Line %d: %q has no value before it", lineNumber, word)
			}

			originalValue := cleanToken(output[targetIndex])
			number, ok := new(big.Int).SetString(originalValue, base)
			if !ok {
				return "", fmt.Sprintf("Line %d: %q is not a valid %s number", lineNumber, originalValue, numberType)
			}

			output[targetIndex] = strings.Replace(output[targetIndex], originalValue, number.String(), 1)
			continue
		}

		if command == "(up)" || command == "(low)" || command == "(cap)" {
			action := strings.Trim(command, "()")
			indexes := lastTargets(output, 1)
			if len(indexes) == 0 {
				return "", fmt.Sprintf("Line %d: %q has no word before it", lineNumber, word)
			}
			applyCase(&output[indexes[0]], action)
			continue
		}

		action, count, matched := countAction(command)
		if matched {
			if count <= 0 {
				return "", fmt.Sprintf("Line %d: invalid number %d in %q", lineNumber, count, word)
			}

			indexes := lastTargets(output, count)
			
			for _, targetIndex := range indexes {
				applyCase(&output[targetIndex], action)
			}
			continue
		}

		output = append(output, word)
	}

	return strings.Join(output, " "), ""
}
func countAction(word string) (string, int, bool) {
	if !strings.HasPrefix(word, "(") || !strings.HasSuffix(word, ")") || !strings.Contains(word, ",") {
		return "", 0, false
	}

	parts := strings.Split(strings.Trim(word, "()"), ",")
	if len(parts) != 2 {
		return "", 0, false
	}

	action := parts[0]
	if action != "up" && action != "low" && action != "cap" {
		return "", 0, false
	}

	count, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, false
	}
	return action, count, true
}
func lastTarget(words []string) int {
	for index := len(words) - 1; index >= 0; index-- {
		if isTarget(words[index]) {
			return index
		}
	}
	return -1
}
func lastTargets(words []string, count int) []int {
	var indexes []int
	for index := len(words) - 1; index >= 0 && len(indexes) < count; index-- {
		if isTarget(words[index]) {
			indexes = append(indexes, index)
		}
	}
	return indexes
}
func isTarget(word string) bool {
	for _, character := range word {
		if unicode.IsLetter(character) || unicode.IsDigit(character) {
			return true
		}
	}
	return false
}
func cleanToken(word string) string {
	return strings.Trim(word, ".,!?:;'\"()[]{}")
}
func applyCase(word *string, action string) {
	switch action {
	case "up":
		*word = strings.ToUpper(*word)
	case "low":
		*word = strings.ToLower(*word)
	case "cap":
		*word = capitalize(*word)
	}
}
func capitalize(word string) string {
	characters := []rune(strings.ToLower(word))
	for index, character := range characters {
		if unicode.IsLetter(character) {
			characters[index] = unicode.ToUpper(character)
			break
		}
	}
	return string(characters)
}
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	contractionPattern     = regexp.MustCompile(`(?i)\b([a-z]+)\s*'\s*(s|t|d|ve|ll|re|m)\b`)
	countCommandPattern    = regexp.MustCompile(`(?i)\(\s*(up|low|cap)\s*,\s*(-?\d+)\s*\)`)
	singleCommandPattern   = regexp.MustCompile(`(?i)\(\s*(up|low|cap|hex|bin)\s*\)`)
	possibleCommandPattern = regexp.MustCompile(`(?i)\(\s*(up|low|cap|hex|bin)\s*(,[^)]*)?\)`)
	validCommandPattern    = regexp.MustCompile(`(?i)^((\(up\)|\(low\)|\(cap\)|\(hex\)|\(bin\))|\((up|low|cap),-?\d+\))$`)
)

func processText(text string) string {
	lines := strings.Split(text, "\n")

	for index, originalLine := range lines {
		line := strings.TrimSuffix(originalLine, "\r")
		if strings.TrimSpace(line) == "" {
			lines[index] = ""
			continue
		}

		line = normalizeInput(line)
		if message := validateLine(line, index+1); message != "" {
			lines[index] = message
			continue
		}

		processedLine, message := applyCommands(line, index+1)
		if message != "" {
			lines[index] = message
			continue
		}

		processedLine = fixArticles(processedLine)
		processedLine = fixPunctuation(processedLine)
		processedLine = fixQuotes(processedLine)
		processedLine = fixBrackets(processedLine)
		lines[index] = processedLine
	}

	return strings.Join(lines, "\n")
}
func normalizeInput(text string) string {
	text = contractionPattern.ReplaceAllString(text, "$1'$2")
	text = countCommandPattern.ReplaceAllStringFunc(text, func(command string) string {
		parts := countCommandPattern.FindStringSubmatch(command)
		count, err := strconv.Atoi(parts[2])
		if err != nil {
			return command
		}
		if count < 0 {
			count = -count
		}
		return " (" + strings.ToLower(parts[1]) + "," + strconv.Itoa(count) + ") "
	})
	text = singleCommandPattern.ReplaceAllString(text, " ($1) ")
	return strings.Join(strings.Fields(text), " ")
}
func validateLine(text string, lineNumber int) string {
	if !quotesAreValid(text) {
		return fmt.Sprintf("Line %d: unmatched single quote", lineNumber)
	}

	commandRanges := possibleCommandPattern.FindAllStringIndex(text, -1)
	for _, commandRange := range commandRanges {
		start, end := commandRange[0], commandRange[1]
		command := text[start:end]
		if commandIsNested(text, start) {
			return fmt.Sprintf("Line %d: invalid command %q", lineNumber, strings.ReplaceAll(command, " ", ""))
		}

		compactCommand := strings.ReplaceAll(command, " ", "")
		if !validCommandPattern.MatchString(compactCommand) {
			return fmt.Sprintf("Line %d: invalid command %q", lineNumber, compactCommand)
		}
	}
	return ""
}
func commandIsNested(text string, commandStart int) bool {
	depth := 0
	for _, character := range text[:commandStart] {
		if character == '(' {
			depth++
		} else if character == ')' {
			depth--
		}
	}
	return depth > 0
}
func quotesAreValid(text string) bool {
	characters := []rune(text)
	quoteCount := 0
	for index, character := range characters {
		if character != '\'' {
			continue
		}
		if index > 0 && index+1 < len(characters) && unicode.IsLetter(characters[index-1]) && unicode.IsLetter(characters[index+1]) {
			continue
		}
		quoteCount++
	}
	return quoteCount%2 == 0
}




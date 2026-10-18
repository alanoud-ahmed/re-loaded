package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Error: use go run . input.txt output.txt")
		return
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	if message := validateFiles(inputPath, outputPath); message != "" {
		fmt.Println(message)
		return
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Error: cannot read", inputPath)
		return
	}

	result := processText(string(data))
	if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
		fmt.Println("Error: cannot write", outputPath)
	}
}
func validateFiles(inputPath, outputPath string) string {
	if !strings.EqualFold(filepath.Ext(inputPath), ".txt") {
		return fmt.Sprintf("Error: input file %q must be a .txt file", inputPath)
	}
	if !strings.EqualFold(filepath.Ext(outputPath), ".txt") {
		return fmt.Sprintf("Error: output file %q must be a .txt file", outputPath)
	}

	absoluteInput, _ := filepath.Abs(inputPath)
	absoluteOutput, _ := filepath.Abs(outputPath)
	if absoluteInput == absoluteOutput {
		return "Error: input and output files cannot be the same"
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Sprintf("Error: input file %q does not exist", inputPath)
	}
	if info.IsDir() {
		return fmt.Sprintf("Error: %q is a folder, not a file", inputPath)
	}
	return ""
}



package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Use: go run . input.txt output.txt")
		return
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	content, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	finalText := Process(string(content))

	err = os.WriteFile(outputPath, []byte(finalText), 0644)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	fmt.Printf("Processing complete. Result written to %s\n", outputPath)
}

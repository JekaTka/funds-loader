package main

import (
	"github.com/JekaTka/funds-loader/fundsprocessor"
	"log"
	"os"
)

func main() {
	inputFile, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create("./final-output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	if err := fundsprocessor.New(inputFile).ProcessTo(outputFile); err != nil {
		log.Fatal(err)
	}
}

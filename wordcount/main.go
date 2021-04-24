package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

const totalDivisions = 5

func processInput(filename string) ([]string, int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	text := ""

	buf := make([]byte, 32*1024)
	reader := bufio.NewReader(file)

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		text += string(buf[:n])
	}

	text = strings.ToLower(text)
	words := strings.Split(text, " ")
	chunkSize := len(words) / 5
	// fmt.Printf("Total size %v and division %v and the remainder %v\n", len(words), len(words)/5, len(words)%5)

	return words, chunkSize
}

func countWords(words []string, output chan map[string]int) {
	counts := make(map[string]int)
	for _, word := range words {
		_, ok := counts[word]
		if ok {
			counts[word] += 1
		} else {
			counts[word] = 1
		}
	}

	output <- counts
}

func main() {
	processInput("ex_input.txt")
}

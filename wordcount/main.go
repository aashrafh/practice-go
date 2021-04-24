package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

const totalDivisions = 5

type shmMap struct {
	mtx    sync.Mutex
	counts map[string]int
}

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

func countWords(words []string, i int, mp *shmMap) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	fmt.Printf("Proc %v\n", i)
	for _, word := range words {
		_, ok := mp.counts[word]
		if ok {
			mp.counts[word] += 1
		} else {
			mp.counts[word] = 1
		}
	}

	fmt.Println(mp.counts)
}

func Reducer(words []string, chunkSize int) {
	mp := shmMap{counts: make(map[string]int)}

	var procGroup sync.WaitGroup
	for i := 0; i < 5; i++ {
		procGroup.Add(1)
		start := i * chunkSize
		var end int
		if end = (i + 1) * chunkSize; i == 4 {
			end = len(words)
		}
		go func(words []string, i int, mp *shmMap) {
			defer procGroup.Done()
			countWords(words, i, mp)
		}(words[start:end], i, &mp)
	}
	procGroup.Wait()

	fmt.Println(mp.counts)
}

func main() {
	words, chunkSize := processInput("ex_input.txt")
	Reducer(words, chunkSize)
}

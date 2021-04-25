package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

const totalDivisions = 5

type shMap struct {
	mtx    sync.Mutex
	counts map[string]int
}
type Pair struct {
	key   string
	value int
}
type PairList []Pair

func (p PairList) Len() int      { return len(p) }
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool {
	if p[i].value == p[j].value {
		return p[i].key < p[j].key
	}
	return p[i].value > p[j].value
}

func processInput(filename string) []string {
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

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = string(line)
		line = strings.ToLower(line)
		wordsSeperated := strings.Split(line, " ")
		words = append(words, wordsSeperated...)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return words
}

func countWords(words []string, i int, mp *shMap) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	for _, word := range words {
		_, ok := mp.counts[word]
		if ok {
			mp.counts[word] += 1
		} else {
			mp.counts[word] = 1
		}
	}
}

func sortMap(mp map[string]int) []string {
	pairs := make(PairList, len(mp))

	i := 0
	for key, value := range mp {
		pairs[i] = Pair{key, value}
		i++
	}

	sort.Sort(pairs)

	var result []string
	for _, k := range pairs {
		result = append(result, k.key)
	}

	return result
}

func writeResult(fileName string, counts map[string]int, sortedKeys []string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, key := range sortedKeys {
		_, err := f.WriteString(fmt.Sprintf("%v : %v \n", key, counts[key]))
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Reducer(words []string) {
	chunkSize := len(words) / 5
	mp := shMap{counts: make(map[string]int)}

	var procGroup sync.WaitGroup
	for i := 0; i < 5; i++ {
		procGroup.Add(1)
		start := i * chunkSize
		var end int
		if end = (i + 1) * chunkSize; i == 4 {
			end = len(words)
		}
		go func(words []string, i int, mp *shMap) {
			defer procGroup.Done()
			countWords(words, i, mp)
		}(words[start:end], i, &mp)
	}
	procGroup.Wait()

	sortedKeys := sortMap(mp.counts)
	writeResult("WordCountOutput.txt", mp.counts, sortedKeys)
}

func main() {
	words := processInput("test.txt")
	Reducer(words)
}

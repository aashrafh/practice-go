package main

import (
	"bufio"
	"io"
	"os"
)

func processInput(filename string) chan string {
	output := make(chan string)

	go func() {
		file, err := os.Open(filename)
		if err != nil {
			return
		}
		defer file.Close()

		reader := bufio.NewReader(file)

		for {
			line, err := reader.ReadString('\n')
			output <- line
			if err == io.EOF {
				break
			}
		}

		close(output)
	}()
	return output
}

func main() {

}

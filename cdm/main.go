package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Download struct {
	URL              string
	Path             string
	TotalConnections int
}

func (download Download) DownloadFile() error {
	fmt.Printf("Starting the connection...\n")

	request, err := download.getHttpRequest("HEAD")
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode > 299 {
		return errors.New(fmt.Sprintf("Failed to process the request, %v error", response.StatusCode))
	} else {
		fmt.Printf("HTTP response status code %v\n", response.StatusCode)
	}

	size, err := strconv.Atoi(response.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	fmt.Printf("Size of the file = %.2f MiB\n", float64(size)/1048576.0)

	connections := make([][2]int, download.TotalConnections)
	connectionSize := size / download.TotalConnections
	for i := range connections {
		if i == 0 {
			connections[i][0] = 0
		} else {
			connections[i][0] = connections[i-1][1] + 1
		}

		if i == download.TotalConnections-1 {
			connections[i][1] = size - 1
		} else {
			connections[i][1] = connections[i][0] + connectionSize
		}
	}

	var procGroup sync.WaitGroup
	for i, connection := range connections {
		procGroup.Add(1)
		go func(i int, connection [2]int) {
			defer procGroup.Done()
			err = download.downloadChunk(i, connection)
			if err != nil {
				panic(err)
			}
		}(i, connection)
	}
	procGroup.Wait()

	return download.mergeChunks(connections)
}

func (download Download) getHttpRequest(method string) (*http.Request, error) {
	request, err := http.NewRequest(
		method,
		download.URL,
		nil,
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Concurrent Download Manager v1.0")
	return request, nil
}

func (download Download) downloadChunk(i int, connection [2]int) error {
	request, err := download.getHttpRequest("GET")
	if err != nil {
		return err
	}

	request.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", connection[0], connection[1]))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode > 299 {
		return errors.New(fmt.Sprintf("Failed to process the request, %v error", response.StatusCode))
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("chunk-%v.tmp", i), bytes, os.ModePerm)
	if err != nil {
		return nil
	}

	return nil
}

func (download Download) mergeChunks(connections [][2]int) error {
	file, err := os.OpenFile(download.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := range connections {
		name := fmt.Sprintf("chunk-%v.tmp", i)
		bytes, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}

		total, err := file.Write(bytes)
		if err != nil {
			return err
		}

		err = os.Remove(name)
		if err != nil {
			return err
		}

		fmt.Printf("%v bytes has downloaded and merged\n", total)
	}
	return nil
}

func main() {
	urlPtr := flag.String("url", "", "A link to the file to be downloaded")
	pathPtr := flag.String("path", "download.tmp", "The path to store the downloaded file")
	flag.Parse()

	if *urlPtr == "" {
		log.Printf("URL of the file is empty, please try again with a valid URL\n")
		os.Exit(1)
	}

	startTime := time.Now()
	download := Download{
		URL:              *urlPtr,
		Path:             *pathPtr,
		TotalConnections: 10,
	}
	err := download.DownloadFile()
	if err != nil {
		log.Printf("An error occured while downloading the file: %s\n", err)
	}
	fmt.Printf("Download completed in %v seconds\n", time.Now().Sub(startTime).Seconds())
}

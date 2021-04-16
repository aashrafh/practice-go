package main

import (
	"errors"
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
	fmt.Printf("Size of the file = %v bytes\n", size)

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

	return nil
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
	} else {
		fmt.Printf("Downloaded %v bytes  for connection %v\n", response.Header.Get("Content-Length"), i)
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

func main() {
	startTime := time.Now()
	download := Download{
		URL:              "https://d1os9znak2uyeg.cloudfront.net/content/04846b59-0840-5c1f-81d5-022c53992031/21/04/15/09/04846b59-0840-5c1f-81d5-022c53992031_1_210415T091101027Z.mp4?Expires=1618647615&Signature=ZkOBax~0Lijxxmy~JvrGZlWa9Qh4Q44gYteYG1ui8MrZWjTCZSAnWjPJKU8VYdCvAIHucFPLn6tNFz~S83fL7-n17vq-dsYJaram0ccrdC7GjccRMSBzAWBxZo3oMTouuLvqcJBG6qNCd8gag-Hn3P2ben9uMmVssyAm7W2pgf7vxQUb4KpVJ72CgtNkgo4fKTKIXyu-kSMwwYhtHTD-NslJxnflZTACoZjVr~-H3DxyPT7~yKlz63rhZJ6QDZ83Mh~F5mhC9G~Hkhx6tGFYp7fMovihtVghmNZatOP0pblU7Ef9bpisEFgiYO5p4~yA8e1rlbLw3~n-TDoZ7EJy7g__&Key-Pair-Id=APKAIOBDBIMXUOQOBYVA",
		Path:             "lec4.mp4",
		TotalConnections: 10,
	}
	err := download.DownloadFile()
	if err != nil {
		log.Printf("An error occured while downloading the file: %s\n", err)
	}
	fmt.Printf("Download completed in %v seconds\n", time.Now().Sub(startTime).Seconds())
}

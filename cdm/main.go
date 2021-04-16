package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"errors"
	"strconv"
)

type Download struct {
	URL string
	Path string
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
	if err != nil{
		return err
	}
	fmt.Printf("Size of the file = %v bytes\n", size)

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


func main()  {
	startTime := time.Now()
	download := Download{
		URL: "https://d1os9znak2uyeg.cloudfront.net/content/b49d56d8-0640-5c1f-9614-022c53992031/21/04/15/08/b49d56d8-0640-5c1f-9614-022c53992031_1_210415T084737293Z.mp4?Expires=1618571777&Signature=d3Vyplv36Dvj3Vws8e9MejEwIIjkIoEhlcHezVSoYZuY12in6MfaTFyFtlrftXfhEzpu9Id~v0Ik6kEJBdQZFvEzmZDPJSCzyFQ7OfFzmF-1MIYE~35lLdrmALXF3iAv9tggwZgrDKFyxUEGzN~~PpIiRY6URa6JQiU30gYfgrHYmW6UPfTitQGLckHpX2wxUU6qHk9zB1by8-RTXA0e8wzXdRXz8SAvmJSepi~WN9RjQ1-yUrk40cHvFH2pOyRA8ka8cIRdU1xsg0FqCyKjl51aH5HrWpf~tuJJOFRlQrJmNT8m9k8z9zX82tvkZSXTu7XD8Y5kfyOUcSCEFPlV7w__&Key-Pair-Id=APKAIOBDBIMXUOQOBYVA",
		Path: "lec4.mp4",
		TotalConnections: 10,
	}
	err := download.DownloadFile()
	if err != nil {
		log.Printf("An error occured while downloading the file: %s\n", err)
	}
	fmt.Printf("Download completed in %v seconds\n", time.Now().Sub(startTime).Seconds())
}
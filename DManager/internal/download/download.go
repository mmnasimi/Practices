package internal

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type DownloadRequest struct {
	url       string
	speed     float32
	bandWidth float32
}

func Download(url string, queue string) string {
	if queue == "" {
		queue = "default"
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Test-Header", "testtsets")
	errHandle(err)
	req.Write(os.Stdout)
	fmt.Println("========\nHost:", req.Host)

	resp, err := http.DefaultClient.Do(req)
	errHandle(err)
	defer resp.Body.Close()

	fmt.Println("status code = ", resp.StatusCode)
	for k, vs := range resp.Header {
		fmt.Printf("%s: %d, %+v\n", k, len(vs), vs)
	}

	message := "Download completed"
	// Determine file type from Content-Type header
	contentType := resp.Header.Get("Content-Type")
	var fileExtension string
	switch contentType {
	case "text/html":
		fileExtension = ".html"
	case "application/json":
		fileExtension = ".json"
	case "text/plain":
		fileExtension = ".txt"
	case "audio/mp3", "audio/mpeg":
		fileExtension = ".mp3"
	default:
		message = message + "\nPlease set the correct extension manually"
		fileExtension = ".bin" // Default to binary if type is unknown
	}

	// Create a file with the appropriate extension
	fileName := "downloaded_filedsv" + fileExtension
	file, err := os.Create(fileName)
	errHandle(err)
	defer file.Close()

	// Write response body to the file
	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		_, err := file.WriteString(sc.Text())
		errHandle(err)
	}

	fmt.Printf("File downloaded successfully: %s\n", fileName)

	return message
}

func errHandle(err error) {
	if err != nil {
		logrus.Fatalln(err)
	}
}

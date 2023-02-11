package services

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func getRandomString() string {
	rand.Seed(time.Now().Unix())
	var output strings.Builder
	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	length := 20
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

func downloadSaveFile(url string, extension string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", errors.New("received non 200 response code")
	}
	defer response.Body.Close()

	filename := getRandomString() + "." + extension
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	counter := &WriteCounter{}
	_, err = io.Copy(file, io.TeeReader(response.Body, counter))
	if err != nil {
		return "", err
	}

	return "Success download filename: " + filename, nil
}

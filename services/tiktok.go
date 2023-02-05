package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/dustin/go-humanize"
)

type Tiktok struct {
	URL      string
	VideoURL string
}

func (t *Tiktok) HTMLExtract() error {

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(t.URL),
		chromedp.OuterHTML("html", &res),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	t.extract(res)
	return nil
}

func (t *Tiktok) extract(res string) {
	fmt.Println(res)
	srcRegex := regexp.MustCompile(`<video.*?src="(.*?)".*?</video>`)
	src := srcRegex.FindStringSubmatch(res)[1]
	fmt.Println(src)
	t.VideoURL = src
}

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

func (t *Tiktok) DownloadFile() (string, error) {
	t.HTMLExtract()

	// TODO
	// THIS somethimes files, try n times
	response, err := http.Get(t.VideoURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", errors.New("Received non 200 response code")
	}

	filename := getRandomString() + ".mp4"
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

package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/chromedp/chromedp"
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
	srcRegex := regexp.MustCompile(`<video.*?src="(.*?)".*?</video>`)
	src := srcRegex.FindStringSubmatch(res)[1]
	t.VideoURL = src
}

func (t *Tiktok) download() (string, error) {
	t.HTMLExtract()

	response, err := http.Get(t.VideoURL)

	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", errors.New("Received non 200 response code")
	}
	defer response.Body.Close()

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

func (t *Tiktok) DownloadFile() (string, error) {
	retries := 5
	var s string
	var err error

	fmt.Println("Searching file url...")

	for i := 0; i < retries; i++ {
		s, err = t.download()
		if err == nil {
			return s, nil
		}
	}

	return s, err
}

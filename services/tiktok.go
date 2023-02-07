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

func (t *Tiktok) GetHTMLString() (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(t.URL),
		chromedp.OuterHTML("html", &res),
	)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (t *Tiktok) setVideoURL(res string) error {
	srcRegex := regexp.MustCompile(`<video.*?src="(.*?)".*?</video>`)
	src := srcRegex.FindStringSubmatch(res)
	if len(src) < 2 {
		return errors.New("video URL not found")
	}

	t.VideoURL = src[1]
	return nil
}

func (t *Tiktok) downloadSaveFile() (string, error) {
	response, err := http.Get(t.VideoURL)

	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", errors.New("received non 200 response code")
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
	var response string
	var err error

	fmt.Println("Searching file url...")

	for i := 0; i < retries; i++ {
		s, err := t.GetHTMLString()
		if err != nil {
			continue
		}

		err = t.setVideoURL(s)
		if err != nil {
			continue
		}

		response, err := t.downloadSaveFile()
		if err == nil {
			return response, nil
		}
	}

	return response, err
}

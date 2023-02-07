package services

import (
	"errors"
	"io"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

type Instagram struct {
	URL         string
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	SiteName    string `json:"site_name"`
	Video       string `json:"video"`
}

func (i *Instagram) DownloadFile() (string, error) {
	err := i.getHTMLMeta()
	if err != nil {
		return "", err
	}

	var url string
	var filenameExt string

	if i.Video != "" {
		url = i.Video
		filenameExt = "mp4"
	} else {
		url = i.Image
		filenameExt = "jpeg"
	}

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", errors.New("received non 200 response code")
	}

	filename := getRandomString() + "." + filenameExt
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

func (i *Instagram) getHTMLMeta() error {
	req, err := http.NewRequest("GET", i.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Instagram 10.3.2 (iPhone7,2; iPhone OS 9_3_3; en_US; en-US; scale=2.00; 750x1334) AppleWebKit/420+")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	defer resp.Body.Close()

	i.extract(resp.Body)
	return nil
}

func (i *Instagram) extract(resp io.Reader) {
	z := html.NewTokenizer(resp)

	titleFound := false

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "title" {
				titleFound = true
			}
			if t.Data == "meta" {
				desc, ok := extractMetaProperty(t, "description")
				if ok {
					i.Description = desc
				}

				ogTitle, ok := extractMetaProperty(t, "og:title")
				if ok {
					i.Title = ogTitle
				}

				ogDesc, ok := extractMetaProperty(t, "og:description")
				if ok {
					i.Description = ogDesc
				}

				ogImage, ok := extractMetaProperty(t, "og:image")
				if ok {
					i.Image = ogImage
				}

				ogSiteName, ok := extractMetaProperty(t, "og:site_name")
				if ok {
					i.SiteName = ogSiteName
				}

				ogVideo, ok := extractMetaProperty(t, "og:video")
				if ok {
					i.Video = ogVideo
				}
			}
		case html.TextToken:
			if titleFound {
				t := z.Token()
				i.Title = t.Data
				titleFound = false
			}
		}
	}
}

func extractMetaProperty(t html.Token, prop string) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == "property" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}
	return
}

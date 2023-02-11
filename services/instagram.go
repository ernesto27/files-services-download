package services

import (
	"errors"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

type Instagram struct {
	URL         string
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image"`
	SiteName    string `json:"site_name"`
	VideoURL    string `json:"video"`
}

func (i *Instagram) DownloadFile() (string, error) {
	resp, err := i.getBody()
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	err = i.setValues(resp.Body)
	if err != nil {
		return "", err
	}

	var url string
	var filenameExt string

	if i.VideoURL != "" {
		url = i.VideoURL
		filenameExt = "mp4"
	} else {
		url = i.ImageURL
		filenameExt = "jpeg"
	}

	r, err := downloadSaveFile(url, filenameExt)
	if err != nil {
		return "", err
	}

	return r, nil
}

func (i *Instagram) getBody() (*http.Response, error) {
	req, err := http.NewRequest("GET", i.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Instagram 10.3.2 (iPhone7,2; iPhone OS 9_3_3; en_US; en-US; scale=2.00; 750x1334) AppleWebKit/420+")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("received non 200 response code")
	}

	return resp, nil
}

func (i *Instagram) setValues(resp io.Reader) error {
	z := html.NewTokenizer(resp)

	titleFound := false

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return nil
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
					i.ImageURL = ogImage
				}

				ogSiteName, ok := extractMetaProperty(t, "og:site_name")
				if ok {
					i.SiteName = ogSiteName
				}

				ogVideo, ok := extractMetaProperty(t, "og:video")
				if ok {
					i.VideoURL = ogVideo
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

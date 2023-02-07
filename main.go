package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/ernesto27/files-services-download/services"
)

func Builder(urlStr string) (Service, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	switch u.Host {
	case "www.instagram.com":
		return &services.Instagram{
			URL: urlStr,
		}, nil
	case "www.tiktok.com":
		return &services.Tiktok{
			URL: urlStr,
		}, nil
	default:
		return nil, errors.New("not valid url host")
	}

}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage ./files-download URLIMAGE")
		return
	}

	service, err := Builder(os.Args[1])
	if err != nil {
		panic(err)
	}

	resp, err := service.DownloadFile()
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

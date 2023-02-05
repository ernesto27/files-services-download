package main

type Service interface {
	DownloadFile() (string, error)
}

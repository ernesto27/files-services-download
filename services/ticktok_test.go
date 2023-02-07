package services

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGetHTMLString(t *testing.T) {
	t.Run("Should return a string html", func(t *testing.T) {
		tiktok := Tiktok{
			URL: "https://example.com/",
		}

		htmlString, err := tiktok.GetHTMLString()
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if len(htmlString) == 0 {
			t.Errorf("HTML string is empty")
		}
	})
}

func TestSetVideoUrl(t *testing.T) {
	t.Run("Video URL is found", func(t *testing.T) {

		tiktok := Tiktok{}

		htmlString := `<html>
	<head>
		<title>My Video Page</title>
	</head>
	<body>
		<h1>Here's a Video:</h1>
		<video width="320" height="240" controls src="https://www.example.com/video.mp4">Your browser does not support the video tag.</video>
	</body>
	</html>`

		tiktok.setVideoURL(htmlString)

		if tiktok.VideoURL != "https://www.example.com/video.mp4" {
			t.Errorf("VideoURL was incorrect, got: %s, want: %s.", tiktok.VideoURL, "https://www.example.com/video.mp4")
		}
	})

	t.Run("Video URL is not found", func(t *testing.T) {

		tiktok := Tiktok{}

		htmlString := `<html>
	<head>
		<title>My Video Page</title>
	</head>
	<body>
		<h1>No video tag</h1>
	</body>
	</html>`

		err := tiktok.setVideoURL(htmlString)
		if err == nil {
			t.Errorf("Expected error, got: %v", err)
		}
	})
}

func TestDownloadSaveVideo(t *testing.T) {
	t.Run("Should download and save video", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "https://server.com/video.mp4",
			httpmock.NewBytesResponder(200, []byte("test")))

		tiktok := Tiktok{
			VideoURL: "https://server.com/video.mp4",
		}

		_, err := tiktok.downloadSaveFile()

		if err != nil {
			t.Errorf("Error: %v", err)
		}
		removeMP4TestFiles()
	})

	t.Run("Should return error not valid URL", func(t *testing.T) {

		tiktok := Tiktok{
			VideoURL: "https://notvalid.3333",
		}

		_, err := tiktok.downloadSaveFile()

		if err == nil {
			t.Errorf("Expected error, got: %v", err)
		}

	})
}

func removeMP4TestFiles() {
	dir := "./"

	files, err := filepath.Glob(dir + "*.mp4")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("All .mp4 files removed.")
}

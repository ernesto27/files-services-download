package services

import (
	"bytes"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGetBody(t *testing.T) {
	t.Run("Should return a http.Response valid url", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "https://server.com/image",
			httpmock.NewBytesResponder(200, []byte("test")))

		instagram := Instagram{
			URL: "https://server.com/image",
		}

		resp, err := instagram.getBody()
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if resp == nil {
			t.Errorf("Response is nil")
		}
	})

	t.Run("Should return error if not 200 status repsonse", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", "https://server.com/invalidimage",
			httpmock.NewBytesResponder(404, []byte("test")))

		instagram := Instagram{
			URL: "https://server.com/invalidimage",
		}

		_, err := instagram.getBody()
		if err == nil {
			t.Errorf("Expected error, got: %v", err)
		}
	})
}

func TestSetValues(t *testing.T) {
	t.Run("Should set values instagram struct", func(t *testing.T) {
		body := `<html>
  <head>
    <title>Example Page</title>
	<meta property="og:image" content="https://server.com/image.jpg">
  </head>
  <body>
    <h1>Hello, World!</h1>
  </body>
</html>`

		reader := bytes.NewReader([]byte(body))

		instagram := Instagram{}
		err := instagram.setValues(reader)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if instagram.ImageURL != "https://server.com/image.jpg" {
			t.Errorf("Expected: https://server.com/image.jpg, got: %v", instagram.ImageURL)
		}

	})

	t.Run("Image URL is not found", func(t *testing.T) {

		body := `<html>
  <head>
    <title>Example Page</title>
  </head>
  <body>
    <h1>Hello, World!</h1>
  </body>
</html>`

		reader := bytes.NewReader([]byte(body))

		instagram := Instagram{}
		err := instagram.setValues(reader)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if instagram.ImageURL != "" {
			t.Errorf("Expected empty image, got: %v", instagram.ImageURL)
		}
	})
}

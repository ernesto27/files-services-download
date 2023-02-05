package main

import (
	"testing"

	"github.com/ernesto27/files-services-download/services"
)

func TestBuilder(t *testing.T) {

	t.Run("Should get a instragram service instance", func(t *testing.T) {
		b, err := Builder("https://www.instagram.com/p/CF1Z5Z8J8ZU/")
		if err != nil {
			t.Error(err)
		}
		if _, ok := b.(*services.Instagram); !ok {
			t.Error("b is not a pointer to Instagram")
		}
	})

	t.Run("Should get a ticktok service instance", func(t *testing.T) {
		b, err := Builder("https://www.tiktok.com/@mar.iac7/video/7191316200914341166?q=alvarez&t=1675615631355")
		if err != nil {
			t.Error(err)
		}

		if _, ok := b.(*services.Tiktok); !ok {
			t.Error("b is not a pointer to Ticktok")
		}
	})

	t.Run("Should get error on invalid url", func(t *testing.T) {
		_, err := Builder("invalid-url")
		if err == nil {
			t.Error("Should get error on invalid url")
		}
	})

}

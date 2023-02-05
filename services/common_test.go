package services

import "testing"

func TestRandomString(t *testing.T) {
	t.Run("Should get a random string", func(t *testing.T) {
		str := getRandomString()
		if len(str) != 20 {
			t.Error("Should get a string with 20 characters")
		}
	})
}

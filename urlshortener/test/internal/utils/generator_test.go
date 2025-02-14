package utils

import (
	"github.com/Bl00mGuy/url-shortener/internal/utils"
	"testing"
)

func TestGenerate(t *testing.T) {
	shortURL := utils.Generate()
	if len(shortURL) != utils.ShortUrlLength {
		t.Errorf("Expected URL length of %d, but got %d", utils.ShortUrlLength, len(shortURL))
	}
}

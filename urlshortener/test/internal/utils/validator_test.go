package utils

import (
	"github.com/Bl00mGuy/url-shortener/internal/utils"
	"testing"
)

func TestValidateUrl(t *testing.T) {
	validURL := "https://example.com"

	if err := utils.ValidateUrl(validURL); err != nil {
		t.Errorf("Expected valid URL, but got error: %v", err)
	}
}

func TestValidateShortUrl(t *testing.T) {
	validShortURL := "abc123_DEF"
	invalidShortURL := "invalid-url!"

	if err := utils.ValidateShortUrl(validShortURL); err != nil {
		t.Errorf("Expected valid short URL, but got error: %v", err)
	}

	if err := utils.ValidateShortUrl(invalidShortURL); err == nil {
		t.Error("Expected error for invalid short URL, but got none")
	}
}

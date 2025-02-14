package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Bl00mGuy/url-shortener/internal/repository"
)

func TestSave_Success(t *testing.T) {
	store := repository.NewInMemoryURLStore()

	err := store.Save("https://example.com", "short.ly/mock123")

	assert.NoError(t, err)

	shortURL, err := store.GetShort("https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "short.ly/mock123", shortURL)

	originalURL, err := store.GetOriginal("short.ly/mock123")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", originalURL)
}

func TestSave_OriginalURLAlreadyMapped(t *testing.T) {
	store := repository.NewInMemoryURLStore()

	_ = store.Save("https://example.com", "short.ly/mock123")

	err := store.Save("https://example.com", "short.ly/mock456")

	assert.EqualError(t, err, "original URL already mapped to a different short URL: short.ly/mock123")
}

func TestSave_ShortURLExists(t *testing.T) {
	store := repository.NewInMemoryURLStore()

	_ = store.Save("https://example.com", "short.ly/mock123")

	err := store.Save("https://another.com", "short.ly/mock123")

	assert.EqualError(t, err, "short URL already exists")
}

func TestGetOriginal_Success(t *testing.T) {
	store := repository.NewInMemoryURLStore()

	_ = store.Save("https://example.com", "short.ly/mock123")

	originalURL, err := store.GetOriginal("short.ly/mock123")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", originalURL)
}

func TestGetOriginal_NotFound(t *testing.T) {
	store := repository.NewInMemoryURLStore()

	originalURL, err := store.GetOriginal("short.ly/nonexistent")
	assert.EqualError(t, err, "short URL not found")
	assert.Empty(t, originalURL)
}

func TestGetShort_Success(t *testing.T) {
	store := repository.NewInMemoryURLStore()

	_ = store.Save("https://example.com", "short.ly/mock123")

	shortURL, err := store.GetShort("https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, "short.ly/mock123", shortURL)
}

func TestGetShort_NotFound(t *testing.T) {
	store := repository.NewInMemoryURLStore()

	shortURL, err := store.GetShort("https://nonexistent.com")
	assert.EqualError(t, err, "original URL not found")
	assert.Empty(t, shortURL)
}

func TestClose(t *testing.T) {
	store := repository.NewInMemoryURLStore()

	err := store.Close()
	assert.NoError(t, err)
}

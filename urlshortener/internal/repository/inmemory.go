package repository

import (
	"errors"
	"sync"
)

type InMemoryURLStore struct {
	mu      sync.RWMutex
	urlData map[string]*URLData
}

type URLData struct {
	OriginalURL string
	ShortURL    string
}

func NewInMemoryURLStore() URLStorage {
	return &InMemoryURLStore{
		urlData: make(map[string]*URLData),
	}
}

func (store *InMemoryURLStore) Save(originalURL, shortURL string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if err := store.isOriginalURLMapped(originalURL, shortURL); err != nil {
		return err
	}

	if err := store.isShortURLExists(shortURL); err != nil {
		return err
	}

	return store.storeURLData(originalURL, shortURL)
}

func (store *InMemoryURLStore) isOriginalURLMapped(originalURL, shortURL string) error {
	if existingData, exists := store.urlData[originalURL]; exists && existingData.ShortURL != shortURL {
		return errors.New("original URL already mapped to a different short URL: " + existingData.ShortURL)
	}
	return nil
}

func (store *InMemoryURLStore) isShortURLExists(shortURL string) error {
	if _, exists := store.urlData[shortURL]; exists {
		return errors.New("short URL already exists")
	}
	return nil
}

func (store *InMemoryURLStore) storeURLData(originalURL, shortURL string) error {
	store.urlData[originalURL] = &URLData{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}
	store.urlData[shortURL] = &URLData{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}
	return nil
}

func (store *InMemoryURLStore) GetOriginal(shortURL string) (string, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	data, exists := store.urlData[shortURL]
	if !exists {
		return "", errors.New("short URL not found")
	}

	return data.OriginalURL, nil
}

func (store *InMemoryURLStore) GetShort(originalURL string) (string, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	data, exists := store.urlData[originalURL]
	if !exists {
		return "", errors.New("original URL not found")
	}

	return data.ShortURL, nil
}

func (store *InMemoryURLStore) Close() error {
	return nil
}

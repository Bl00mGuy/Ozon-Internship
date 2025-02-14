package repository

type URLStorage interface {
	Save(originalURL, shortURL string) error
	GetOriginal(shortURL string) (string, error)
	GetShort(originalURL string) (string, error)
	Close() error
}

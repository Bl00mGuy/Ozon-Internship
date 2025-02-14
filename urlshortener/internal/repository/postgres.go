package repository

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	Db *sql.DB
}

func NewPostgresRepository(connStr string) (URLStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &PostgresRepository{Db: db}, nil
}

func (r *PostgresRepository) Save(originalURL, shortURL string) error {
	if exists, err := r.isOriginalURLExists(originalURL); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("original URL already exists")
	}

	if err := r.insertShortenedURL(originalURL, shortURL); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetOriginal(shortURL string) (string, error) {
	originalURL, err := r.fetchOriginalURL(shortURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (r *PostgresRepository) GetShort(originalURL string) (string, error) {
	shortURL, err := r.fetchShortURL(originalURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (r *PostgresRepository) Close() error {
	return r.Db.Close()
}

func (r *PostgresRepository) isOriginalURLExists(originalURL string) (bool, error) {
	query := `SELECT short_url FROM urls WHERE original_url = $1`
	var existingShortURL string
	err := r.Db.QueryRow(query, originalURL).Scan(&existingShortURL)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check if original URL exists: %w", err)
	}
	return true, nil
}

func (r *PostgresRepository) insertShortenedURL(originalURL, shortURL string) error {
	query := `INSERT INTO urls (original_url, short_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING`
	_, err := r.Db.Exec(query, originalURL, shortURL)
	if err != nil {
		return fmt.Errorf("failed to insert shortened URL: %w", err)
	}
	return nil
}

func (r *PostgresRepository) fetchOriginalURL(shortURL string) (string, error) {
	query := `SELECT original_url FROM urls WHERE short_url = $1`
	var originalURL string
	err := r.Db.QueryRow(query, shortURL).Scan(&originalURL)
	if err == sql.ErrNoRows {
		return "", errors.New("short URL not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to fetch original URL: %w", err)
	}
	return originalURL, nil
}

func (r *PostgresRepository) fetchShortURL(originalURL string) (string, error) {
	query := `SELECT short_url FROM urls WHERE original_url = $1`
	var shortURL string
	err := r.Db.QueryRow(query, originalURL).Scan(&shortURL)
	if err == sql.ErrNoRows {
		return "", errors.New("original URL not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to fetch short URL: %w", err)
	}
	return shortURL, nil
}

package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/Bl00mGuy/url-shortener/internal/repository"
)

func TestNewPostgresRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := &repository.PostgresRepository{Db: db}

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.Db)
}

func TestPostgresSave_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := &repository.PostgresRepository{Db: db}

	mock.ExpectQuery(`SELECT short_url FROM urls WHERE original_url = \$1`).
		WithArgs("https://example.com").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec(`INSERT INTO urls \(original_url, short_url\) VALUES \(\$1, \$2\)`).
		WithArgs("https://example.com", "short.ly/mock123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Save("https://example.com", "short.ly/mock123")

	assert.NoError(t, err)
}

func TestPostgresSave_OriginalURLExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := &repository.PostgresRepository{Db: db}

	mock.ExpectQuery(`SELECT short_url FROM urls WHERE original_url = \$1`).
		WithArgs("https://example.com").
		WillReturnRows(sqlmock.NewRows([]string{"short_url"}).AddRow("short.ly/mock123"))

	err = repo.Save("https://example.com", "short.ly/mock456")

	assert.EqualError(t, err, "original URL already exists")
}

func TestPostgresGetOriginal_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := &repository.PostgresRepository{Db: db}

	mock.ExpectQuery(`SELECT original_url FROM urls WHERE short_url = \$1`).
		WithArgs("short.ly/mock123").
		WillReturnRows(sqlmock.NewRows([]string{"original_url"}).AddRow("https://example.com"))

	originalURL, err := repo.GetOriginal("short.ly/mock123")

	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", originalURL)
}

func TestPostgresGetOriginal_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := &repository.PostgresRepository{Db: db}

	mock.ExpectQuery(`SELECT original_url FROM urls WHERE short_url = \$1`).
		WithArgs("short.ly/mock123").
		WillReturnError(sql.ErrNoRows)

	originalURL, err := repo.GetOriginal("short.ly/mock123")

	assert.EqualError(t, err, "short URL not found")
	assert.Empty(t, originalURL)
}

func TestPostgresGetShort_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := &repository.PostgresRepository{Db: db}

	mock.ExpectQuery(`SELECT short_url FROM urls WHERE original_url = \$1`).
		WithArgs("https://example.com").
		WillReturnRows(sqlmock.NewRows([]string{"short_url"}).AddRow("short.ly/mock123"))

	shortURL, err := repo.GetShort("https://example.com")

	assert.NoError(t, err)
	assert.Equal(t, "short.ly/mock123", shortURL)
}

func TestPostgresGetShort_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := &repository.PostgresRepository{Db: db}

	mock.ExpectQuery(`SELECT short_url FROM urls WHERE original_url = \$1`).
		WithArgs("https://nonexistent.com").
		WillReturnError(sql.ErrNoRows)

	shortURL, err := repo.GetShort("https://nonexistent.com")

	assert.EqualError(t, err, "original URL not found")
	assert.Empty(t, shortURL)
}

func TestPostgresClose(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := &repository.PostgresRepository{Db: db}

	mock.ExpectClose()

	err = repo.Close()
	assert.NoError(t, err)
}

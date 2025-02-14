package services

import (
	"context"
	"github.com/Bl00mGuy/url-shortener/internal/services"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Bl00mGuy/url-shortener/proto/gen/go"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(originalUrl, shortUrl string) error {
	args := m.Called(originalUrl, shortUrl)
	return args.Error(0)
}

func (m *MockRepository) GetOriginal(shortUrl string) (string, error) {
	args := m.Called(shortUrl)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) GetShort(originalUrl string) (string, error) {
	args := m.Called(originalUrl)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestShortenUrl_AlreadyExists(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := logrus.New()

	service := services.NewUrlShortener(mockRepo, mockLogger)
	mockRepo.On("GetShort", "http://example.com").Return("short.ly/abc", nil)

	req := &_go.UrlShorteningRequest{OriginalUrl: "http://example.com"}
	resp, err := service.ShortenUrl(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "short.ly/abc", resp.ShortenedUrl)
	mockRepo.AssertExpectations(t)
}

func TestShortenUrl_InvalidUrl(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := logrus.New()

	service := services.NewUrlShortener(mockRepo, mockLogger)

	req := &_go.UrlShorteningRequest{OriginalUrl: "invalid-url"}
	resp, err := service.ShortenUrl(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid original URL")
	assert.Nil(t, resp)
}

func TestExpandUrl_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := logrus.New()

	service := services.NewUrlShortener(mockRepo, mockLogger)
	mockRepo.On("GetOriginal", "short.ly/abc").Return("", status.Error(codes.NotFound, "short URL not found"))

	req := &_go.ShortenedUrlRequest{ShortenedUrl: "short.ly/abc"}
	resp, err := service.ExpandUrl(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "short URL not found")
	assert.Nil(t, resp)
}

func TestShortenUrl_SaveError(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := logrus.New()

	service := services.NewUrlShortener(mockRepo, mockLogger)
	mockRepo.On("GetShort", "http://example.com").Return("", nil)
	mockRepo.On("Save", "http://example.com", mock.Anything).Return(status.Error(codes.Internal, "failed to save"))

	req := &_go.UrlShorteningRequest{OriginalUrl: "http://example.com"}
	resp, err := service.ShortenUrl(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save URL in the repository")
	assert.Nil(t, resp)
}

package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Bl00mGuy/url-shortener/internal/repository"
	"github.com/Bl00mGuy/url-shortener/internal/utils"
	"github.com/Bl00mGuy/url-shortener/proto/gen/go"
)

type UrlShortener struct {
	_go.UnimplementedUrlManagerServer
	repo   repository.URLStorage
	logger *logrus.Logger
}

func NewUrlShortener(repo repository.URLStorage, logger *logrus.Logger) *UrlShortener {
	return &UrlShortener{
		repo:   repo,
		logger: logger,
	}
}

func (u *UrlShortener) ShortenUrl(_ context.Context, req *_go.UrlShorteningRequest) (*_go.ShorteningResult, error) {
	originalUrl := req.GetOriginalUrl()

	if err := u.validateOriginalUrl(originalUrl); err != nil {
		return nil, err
	}

	if shortUrl, err := u.repo.GetShort(originalUrl); err == nil {
		return &_go.ShorteningResult{ShortenedUrl: shortUrl}, nil
	}

	shortUrl := u.generateShortUrl()

	if err := u.repo.Save(originalUrl, shortUrl); err != nil {
		u.logger.Errorf("Error saving URL: %v", err)
		return nil, status.Error(codes.Internal, "failed to save URL in the repository")
	}

	return &_go.ShorteningResult{ShortenedUrl: shortUrl}, nil
}

func (u *UrlShortener) ExpandUrl(_ context.Context, req *_go.ShortenedUrlRequest) (*_go.ExpandedUrlResult, error) {
	shortUrl := req.GetShortenedUrl()

	if err := u.validateShortUrl(shortUrl); err != nil {
		return nil, err
	}

	originalUrl, err := u.repo.GetOriginal(shortUrl)
	if err != nil {
		u.logger.Errorf("Error retrieving original URL: %v", err)
		return nil, status.Error(codes.NotFound, "short URL not found")
	}

	return &_go.ExpandedUrlResult{OriginalUrl: originalUrl}, nil
}

func (u *UrlShortener) validateOriginalUrl(originalUrl string) error {
	if err := utils.ValidateUrl(originalUrl); err != nil {
		u.logger.Warnf("Invalid URL: %s, Error: %v", originalUrl, err)
		return status.Error(codes.FailedPrecondition, "invalid original URL")
	}
	return nil
}

func (u *UrlShortener) validateShortUrl(shortUrl string) error {
	if err := utils.ValidateShortUrl(shortUrl); err != nil {
		u.logger.Warnf("Invalid short URL: %s, Error: %v", shortUrl, err)
		return status.Error(codes.FailedPrecondition, "invalid short URL")
	}
	return nil
}

func (u *UrlShortener) generateShortUrl() string {
	return utils.Generate()
}

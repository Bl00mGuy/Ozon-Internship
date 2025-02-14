package utils

import (
	"net/url"
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	InvalidUrlErrorMsg    = "incorrect URL format"
	InvalidShortUrlFormat = "incorrect short URL format"
)

var (
	urlRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{5,20}$`)
)

func ValidateUrl(originalUrl string) error {
	parsedURL, err := url.ParseRequestURI(originalUrl)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return status.Error(codes.FailedPrecondition, InvalidUrlErrorMsg)
	}

	return nil
}

func ValidateShortUrl(shortUrl string) error {
	if !urlRegex.MatchString(shortUrl) {
		return status.Error(codes.FailedPrecondition, InvalidShortUrlFormat)
	}
	return nil
}

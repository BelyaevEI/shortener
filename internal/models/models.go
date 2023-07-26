package models

import "context"

type (
	Request struct {
		URL string `json:"url"`
	}

	Response struct {
		Result string `json:"result"`
	}

	StorageURL struct {
		UserID      uint64 `json:"userID"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	Storage interface {
		Save(ctx context.Context, url1, url2 string, userID uint64) error
		GetOriginURL(ctx context.Context, shortURL string) (string, error)
		GetShortURL(ctx context.Context, longURL string) (string, error)
		Ping(ctx context.Context) error
		GetUrlsUser(ctx context.Context, userID uint64) ([]StorageURL, error)
	}

	Batch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url,omitempty"`
		ShortURL      string `json:"short_url,omitempty"`
	}
)

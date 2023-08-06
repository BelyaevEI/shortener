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
		UserID      uint32 `json:"userID"`
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
		DeletedFlag bool   `json:"deleted"`
	}

	DeleteURL struct {
		ShortURL string
		UserID   uint32
	}

	Storage interface {
		Save(ctx context.Context, url1, url2 string, userID uint32) error
		GetOriginURL(ctx context.Context, shortURL string) (string, bool, error)
		GetShortURL(ctx context.Context, longURL string) (string, error)
		Ping(ctx context.Context) error
		GetUrlsUser(ctx context.Context, userID uint32) ([]StorageURL, error)
		UpdateDeletedFlag(ctx context.Context, data []DeleteURL) error
	}

	Batch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url,omitempty"`
		ShortURL      string `json:"short_url,omitempty"`
		DeletedFlag   bool   `json:"deleted"`
	}

	KeyID string
)

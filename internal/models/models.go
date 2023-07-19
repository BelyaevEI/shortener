package models

type (
	Request struct {
		URL string `json:"url"`
	}

	Response struct {
		Result string `json:"result"`
	}

	StorageURL struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	Storage interface {
		Save(url1, url2 string) error
		GetOriginURL(inputURL string) (string, error)
		GetShortURL(inputURL string) (string, error)
		Ping() error
	}

	Batch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url,omitempty"`
		ShortURL      string `json:"short_url,omitempty"`
	}
)

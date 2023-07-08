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
)

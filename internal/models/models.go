package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type StorageURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

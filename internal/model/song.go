package model

type Song struct {
	ID          uint   `json:"id"`
	Group       string `json:"group"`
	Title       string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

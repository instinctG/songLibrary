package model

const DateFormat = "01.02.2006"

type Song struct {
	ID          string `json:"id,omitempty" example:"1"`
	Group       string `json:"group,omitempty" example:"Muse"`
	Title       string `json:"song,omitempty" example:"Supermassive Black Hole"`
	ReleaseDate string `json:"releaseDate,omitempty" example:"MM.DD.YYYY"`
	Link        string `json:"link,omitempty" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
	Text        string `json:"text,omitempty" example:"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?"`
}

type UpdateSong struct {
	Group       *string `json:"group"`
	Title       *string `json:"song"`
	ReleaseDate *string `json:"releaseDate"`
	Text        *string `json:"text"`
	Link        *string `json:"link"`
}

type SongRequest struct {
	Group       string `json:"group,omitempty" example:"Muse"`
	Title       string `json:"song,omitempty" example:"Supermassive Black Hole"`
	ReleaseDate string `json:"releaseDate,omitempty" example:"MM.DD.YYYY"`
	Link        string `json:"link,omitempty" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
	Text        string `json:"text,omitempty" example:"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?"`
}

type Verse struct {
	Text string `json:"text" example:"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?"`
}

type SongAddRequest struct {
	Group string `json:"group,omitempty" example:"Muse"`
	Title string `json:"song,omitempty" example:"Supermassive Black Hole"`
}

type MessageResponse struct {
	Message string `json:"message" example:"success"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

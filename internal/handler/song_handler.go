package handler

import (
	"encoding/json"
	"fmt"
	"github.com/instinctG/songLibrary/internal/model"
	"net/http"
)

type SongService interface {
	GetLibrary() ([]model.Song, error)
	GetSongText() (string, error)
	DeleteSong() error
	UpdateSong() (model.Song, error)
	AddSong(group, song string) (model.Song, error)
}

type Response struct {
	Message string `json:"message"`
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

/*
func (h *Handler) GetLibrary(w http.ResponseWriter, r *http.Request) error {
	group := r.URL.Query().Get("group")
	title := r.URL.Query().Get("song")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	songs, err := h.Service.GetLibrary()
	if err != nil {
		return fmt.Errorf("could not get song library: %w", err)
	}

	return WriteJSON(w, http.StatusOK, songs)
}
*/

/*
func (h *Handler) GetSongText(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	songText, err := h.Service.GetSongText()
	if err != nil {
		return fmt.Errorf("could not get the song text :%w", err)
	}

	return WriteJSON(w, http.StatusOK, songText)
}
*/

/*
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.Service.DeleteSong(); err != nil {
		return fmt.Errorf("could not delete song :  %w", err)
	}

	return WriteJSON(w, http.StatusOK, Response{Message: "succesfully deleted the song"})
}
*/

/*
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return fmt.Errorf("id is required")
	}
	numId, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("could not parse num id: %w", err)
	}

	var song model.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		return fmt.Errorf("failed to decode updated song body :%w", err)
	}

	updatedSong, err := h.Service.UpdateSong(id, song)
	if err != nil {
		return fmt.Errorf("could not update the song: %w", err)
	}

	return WriteJSON(w, http.StatusOK, updatedSong)
}
*/

func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) error {

	var song model.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		return fmt.Errorf("could not decode song body: %w", err)
	}

	newSong, err := h.Service.AddSong(song.Group, song.Text)
	if err != nil {
		return fmt.Errorf("could not add song: %w", err)
	}

	return WriteJSON(w, http.StatusOK, newSong)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

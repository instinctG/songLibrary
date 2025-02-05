package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/instinctG/songLibrary/internal/model"
	sl "github.com/instinctG/songLibrary/pkg/logger"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

type SongService interface {
	GetLibrary(ctx context.Context, filter model.Song, limit, offset int) ([]model.Song, error)
	GetSongText(ctx context.Context, id, limit, offset int) (string, error)
	DeleteSong(ctx context.Context, id int) error
	UpdateSong(ctx context.Context, id int, update model.Song) (model.Song, error)
	AddSong(ctx context.Context, song model.Song) (model.Song, error)
}

// @Summary Получить список песен
// @Description Возвращает список песен с возможностью фильтрации по группе, названию, дате выпуска, тексту и ссылке
// @Tags songs
// @Accept  json
// @Produce  json
// @Param group query string false "Фильтр по группе"
// @Param title query string false "Фильтр по названию песни"
// @Param release-date query string false "Фильтр по дате выпуска в формате MM.DD.YYYY"
// @Param text query string false "Фильтр по тексту песни"
// @Param link query string false "Фильтр по ссылке"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество записей на странице (по умолчанию 10)"
// @Success 200 {array} model.Song
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /songs [get]
func (h *Handler) GetLibrary(w http.ResponseWriter, r *http.Request) {

	filter := model.Song{
		Group:       r.URL.Query().Get("group"),
		Title:       r.URL.Query().Get("title"),
		ReleaseDate: r.URL.Query().Get("release-date"),
		Text:        r.URL.Query().Get("text"),
		Link:        r.URL.Query().Get("link"),
	}

	if filter.ReleaseDate != "" {

		if err := validateReleaseDate(filter.ReleaseDate); err != nil {

			h.Log.Debug("invalid release date", sl.Err(err))

			jsonRespond(w, r, http.StatusBadRequest, Response{Error: "invalid release date"})

			return
		}
	}

	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		pageString := os.Getenv("PAGE")
		page, _ = strconv.Atoi(pageString)
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limitString := os.Getenv("LIMIT")
		limit, _ = strconv.Atoi(limitString)
	}

	offset := (page - 1) * limit

	songs, err := h.Service.GetLibrary(r.Context(), filter, limit, offset)
	if err != nil {

		h.Log.Debug("failed to get songs from library", sl.Err(err))

		jsonRespond(w, r, http.StatusInternalServerError, Response{Error: err.Error()})

		return
	}

	jsonRespond(w, r, http.StatusOK, songs)

}

// @Summary Получить текст песни
// @Description Возвращает текст песни по её ID с пагинацией по куплетам
// @Tags song
// @Accept  json
// @Produce  json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество куплетов на странице (по умолчанию 10)"
// @Success 200 {object} model.Verse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /song/{id} [get]
func (h *Handler) GetSongText(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idString)
	if err != nil || id < 1 {

		h.Log.Debug("id is required,must be integer and cannot be less than 1")

		jsonRespond(w, r, http.StatusBadRequest, Response{Error: "id is required,must be integer and cannot be less than 1"})

		return
	}

	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		pageString := os.Getenv("PAGE")
		page, _ = strconv.Atoi(pageString)
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limitString := os.Getenv("VERSE")
		limit, _ = strconv.Atoi(limitString)
	}

	offset := (page - 1) * limit

	text, err := h.Service.GetSongText(r.Context(), id, limit, offset)

	if err != nil {

		h.Log.Debug("failed to get song verse by id", sl.Err(err))

		jsonRespond(w, r, http.StatusInternalServerError, Response{Error: err.Error()})

		return

	}

	jsonRespond(w, r, http.StatusOK, model.Verse{Text: text})
}

// @Summary Добавить песню
// @Description Добавляет новую песню в библиотеку
// @Tags song
// @Accept  json
// @Produce  json
// @Param song body model.SongAddRequest true "Данные песни"
// @Success 200 {object} model.Song
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /song [post]
func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {
	var song model.Song

	err := render.DecodeJSON(r.Body, &song)
	if errors.Is(err, io.EOF) {

		h.Log.Debug("POST /song - add song request body is empty", sl.Err(err))

		jsonRespond(w, r, http.StatusBadRequest, Response{Error: "error decoding POST /song json body"})

		return
	}

	if song.Group == "" || song.Title == "" {

		h.Log.Debug("song group or title cannot be empty")

		jsonRespond(w, r, http.StatusBadRequest, Response{Error: "song group or title cannot be empty"})

		return
	}

	infoURL := fmt.Sprintf("http://localhost:63342/info?group=%s&song=%s", song.Group, song.Title)

	response, err := http.Get(infoURL)
	if err != nil {

		h.Log.Debug("Failed to send GET /info swagger API request", sl.Err(err))

		jsonRespond(w, r, http.StatusInternalServerError, Response{Error: "Failed to send GET /info swagger API request"})

		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {

		h.Log.Debug("Failed to get song info by swagger API", "status code", response.StatusCode)

		jsonRespond(w, r, http.StatusInternalServerError, Response{Error: "Failed to get song info by swagger API"})

		return
	}

	err = render.DecodeJSON(response.Body, &song)
	if errors.Is(err, io.EOF) {

		h.Log.Debug("/info - request body is empty", sl.Err(err))

		jsonRespond(w, r, http.StatusBadRequest, Response{Error: "error decoding song info json body"})

		return
	}

	if song.Link == "" || song.Text == "" {

		h.Log.Debug("song link or text cannot be empty")

		jsonRespond(w, r, http.StatusBadRequest, Response{Error: "song link or text cannot be empty"})

		return
	}

	if err = validateReleaseDate(song.ReleaseDate); err != nil {

		h.Log.Debug("invalid release date", sl.Err(err))

		jsonRespond(w, r, http.StatusBadRequest, Response{Error: "invalid release date"})

		return
	}

	h.Log.Debug("request body decoded", slog.Any("request", song))

	newSong, err := h.Service.AddSong(r.Context(), song)
	if err != nil {

		h.Log.Debug("failed to add song", sl.Err(err))

		jsonRespond(w, r, http.StatusInternalServerError, Response{Error: err.Error()})

		return
	}

	h.Log.Debug("song added successfully")

	jsonRespond(w, r, http.StatusOK, newSong)
}

// @Summary Удалить песню
// @Description Удаляет песню по её ID
// @Tags song
// @Accept  json
// @Produce  json
// @Param id path int true "ID песни"
// @Success 200 {object} model.MessageResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /song/{id} [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idString)

	if err != nil || id < 1 {
		h.Log.Debug("id is required,must be integer and cannot be less than 1")

		jsonRespond(w, r, http.StatusBadRequest, Response{Error: "id is required,must be integer and cannot be less than 1"})

		return
	}

	err = h.Service.DeleteSong(r.Context(), id)

	if err != nil {
		h.Log.Debug("failed to delete song", sl.Err(err))

		jsonRespond(w, r, http.StatusInternalServerError, Response{Error: err.Error()})

		return
	}

	h.Log.Debug("song deleted successfully")

	jsonRespond(w, r, http.StatusOK, Response{Message: "song deleted successfully"})
}

// @Summary Обновить песню
// @Description Обновляет информацию о песне по её ID
// @Tags song
// @Accept  json
// @Produce  json
// @Param id path int true "ID песни"
// @Param song body model.SongRequest true "Данные для обновления песни"
// @Success 200 {object} model.Song
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /song/{id} [put]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {

	var song model.Song

	idString := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idString)
	if err != nil || id < 1 {
		h.Log.Debug("id is required,must be integer and cannot be less than 1")

		jsonRespond(w, r, http.StatusBadRequest, Response{Error: "id is required and cannot be less than 1"})

		return
	}

	err = render.DecodeJSON(r.Body, &song)
	if errors.Is(err, io.EOF) {

		h.Log.Debug("request body is empty", sl.Err(err))

		jsonRespond(w, r, http.StatusBadRequest, Response{Error: "error decoding json body"})

		return
	}

	if song.ReleaseDate != "" {

		if err = validateReleaseDate(song.ReleaseDate); err != nil {

			h.Log.Debug("invalid release date", sl.Err(err))

			jsonRespond(w, r, http.StatusBadRequest, Response{Error: "invalid release date"})

			return
		}
	}

	h.Log.Debug("request body decoded", slog.Any("request", song))

	updatedSong, err := h.Service.UpdateSong(r.Context(), id, song)

	if err != nil {

		h.Log.Debug("failed to update song", sl.Err(err))

		jsonRespond(w, r, http.StatusInternalServerError, Response{Error: err.Error()})

		return
	}

	h.Log.Debug("song updated successfully")

	jsonRespond(w, r, http.StatusOK, updatedSong)
}

package handler

import (
	"github.com/go-chi/render"
	"github.com/instinctG/songLibrary/internal/model"
	"net/http"
	"time"
)

// Response структура для ответа с сообщением или ошибкой.
type Response struct {
	Message string `json:"message,omitempty"` // Сообщение.
	Error   string `json:"error,omitempty"`   // Описание ошибки.
}

func jsonRespond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	render.JSON(w, r, data)
}

func validateReleaseDate(releaseDate string) error {
	_, err := time.Parse(model.DateFormat, releaseDate)
	if err != nil {
		return err
	}

	return nil
}

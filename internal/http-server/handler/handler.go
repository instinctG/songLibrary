package handler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/instinctG/songLibrary/docs"
	logger "github.com/instinctG/songLibrary/internal/http-server/middleware/logger"
	sl "github.com/instinctG/songLibrary/pkg/logger"
	"github.com/swaggo/http-swagger"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Handler struct {
	Log     *slog.Logger
	Router  *chi.Mux
	Service SongService
	Server  *http.Server
}

func NewHandler(service SongService, address string, log *slog.Logger) *Handler {
	h := &Handler{
		Log:     log,
		Service: service,
		Router:  chi.NewRouter(),
	}

	// Настройка middleware
	h.Router.Use(middleware.RequestID) // Генерация идентификаторов запросов.
	//h.Router.Use(middleware.Logger)  //Можно использовать логгер от chi, но решил написать свой для удобства логов
	h.Router.Use(logger.New(log))      // Логирование запросов.
	h.Router.Use(middleware.Recoverer) // Восстановление после паники.
	h.Router.Use(middleware.URLFormat) // Поддержка форматов URL.

	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    address,
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {
	h.Router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	h.Router.Get("/songs", h.GetLibrary)
	h.Router.Get("/song/{id}", h.GetSongText)
	h.Router.Post("/song", h.AddSong)
	h.Router.Put("/song/{id}", h.UpdateSong)
	h.Router.Delete("/song/{id}", h.DeleteSong)

	h.Router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))
}

// Serve запускает HTTP-сервер и обрабатывает сигналы завершения работы(graceful-shutdown).
// Возвращает: ошибку в случае, если сервер не может быть запущен.
func (h *Handler) Serve() error {
	h.Log.Info("starting server on port: " + h.Server.Addr)

	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Ожидание сигнала завершения (graceful shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	start := time.Now()

	h.Log.Info("Received shutdown signal")

	// Завершение работы сервера с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.Log.Info("Server shutdown complete", slog.String("shutdown_time", time.Since(start).String()), slog.String("signal", sig.String()))
	if err := h.Server.Shutdown(ctx); err != nil {
		h.Log.Error("Server Shutdown:", sl.Err(err))
	}

	// Обработка завершения контекста
	select {
	case <-ctx.Done():
		h.Log.Debug("timeout of 5 seconds.")
	}
	h.Log.Debug("Server exiting")

	return nil
}

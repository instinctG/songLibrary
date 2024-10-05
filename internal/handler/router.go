package handler

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Handler struct {
	Router  *mux.Router
	Service SongService
	Server  *http.Server
}

func NewHandler(service SongService, endPoint string) *Handler {
	h := &Handler{Service: service}

	h.Router = mux.NewRouter()
	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    endPoint,
		Handler: h.Router,
	}

	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "pong"}`))
	})
	//h.Router.HandleFunc("/songs/library", makeHTTPHandleFunc(h.GetLibrary)).Methods("GET")
	//h.Router.HandleFunc("/songs/{id}/text", makeHTTPHandleFunc(h.GetSongText)).Methods("GET")
	h.Router.HandleFunc("/songs", makeHTTPHandleFunc(h.AddSong)).Methods("POST")
	//h.Router.HandleFunc("/songs/{id}", makeHTTPHandleFunc(h.UpdateSong)).Methods("PUT")
	//h.Router.HandleFunc("/songs/{id}", makeHTTPHandleFunc(h.DeleteSong)).Methods("DELETE")
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

	return nil
}

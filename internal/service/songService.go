package service

import (
	"context"
	"errors"
	"github.com/instinctG/songLibrary/internal/model"
	"strings"
)

type Store interface {
	GetLibrary(ctx context.Context, filter model.Song, limit, offset int) ([]model.Song, error)
	GetSongText(ctx context.Context, id, limit, offset int) (string, error)
	DeleteSong(ctx context.Context, id int) error
	UpdateSong(ctx context.Context, id int, update model.Song) (model.Song, error)
	AddSong(ctx context.Context, song model.Song) (model.Song, error)
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{Store: store}
}

func (s *Service) AddSong(ctx context.Context, song model.Song) (model.Song, error) {

	newSong, err := s.Store.AddSong(ctx, song)
	if err != nil {
		return model.Song{}, err
	}

	return newSong, nil
}

func (s *Service) GetLibrary(ctx context.Context, filter model.Song, limit, offset int) ([]model.Song, error) {
	songs, err := s.Store.GetLibrary(ctx, filter, limit, offset)

	if err != nil {
		return []model.Song{}, err
	}

	return songs, nil
}

func (s *Service) GetSongText(ctx context.Context, id, limit, offset int) (string, error) {

	text, err := s.Store.GetSongText(ctx, id, limit, offset)
	if err != nil {
		return "", err
	}

	verses := strings.Split(text, "\n\n")

	start := offset
	end := offset + limit

	if start > len(verses) {
		return "", errors.New("verse not found")
	}

	if end > len(verses) {
		end = len(verses)
	}

	paginatedVerses := verses[start:end]
	combinedVerses := strings.Join(paginatedVerses, "\n\n")

	return combinedVerses, nil
}

func (s *Service) DeleteSong(ctx context.Context, id int) error {
	err := s.Store.DeleteSong(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateSong(ctx context.Context, id int, update model.Song) (model.Song, error) {

	updatedSong, err := s.Store.UpdateSong(ctx, id, update)
	if err != nil {
		return model.Song{}, err
	}

	return updatedSong, nil
}

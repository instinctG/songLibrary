package service

import (
	"errors"
	"fmt"
	"github.com/instinctG/songLibrary/internal/model"
)

type Store interface {
	GetLibrary() ([]model.Song, error)
	GetSongText() (string, error)
	DeleteSong() error
	UpdateSong() (model.Song, error)
	AddSong(group, song string) (model.Song, error)
}

var (
	ErrFetching       = errors.New("cannot fetching a song by id")
	ErrNotImplemented = errors.New("not implemented")
)

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{Store: store}
}

func (s *Service) AddSong(group, song string) (model.Song, error) {
	fmt.Println("adding song...")

	songDetail, err := s.Store.AddSong(group, song)
	if err != nil {
		fmt.Println(err) // todo:need to log
		return model.Song{}, fmt.Errorf("error adding song: %w", err)
	}

	return songDetail, nil
}

func (s *Service) GetLibrary() ([]model.Song, error) {

	return []model.Song{}, nil
}

func (s *Service) GetSongText() (string, error) {
	return "", nil
}

func (s *Service) DeleteSong() error {
	err := s.Store.DeleteSong()
	if err != nil {
		return ErrNotImplemented
	}

	return nil
}

func (s *Service) UpdateSong() (model.Song, error) {
	return model.Song{}, nil
}

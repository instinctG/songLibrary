package db

import "github.com/instinctG/songLibrary/internal/model"

func (d *Database) GetLibrary() ([]model.Song, error) {
	return []model.Song{}, nil
}

func (d *Database) GetSongText() (string, error) {
	return "", nil
}

func (d *Database) DeleteSong() error {
	return nil
}

func (d *Database) UpdateSong() (model.Song, error) {
	return model.Song{}, nil
}

func (d *Database) AddSong(group, song string) (model.Song, error) {

	//row, err := d.Client.QueryRow()
	return model.Song{}, nil
}

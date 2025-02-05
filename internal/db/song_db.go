package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/instinctG/songLibrary/internal/model"
	sl "github.com/instinctG/songLibrary/pkg/logger"
	"log"
	"strings"
	"time"
)

func (d *Database) GetLibrary(ctx context.Context, filter model.Song, limit, offset int) ([]model.Song, error) {
	d.Log.Debug("Fetching songs with filters", "group", filter.Group, "song", filter.Title, "releaseDate", filter.ReleaseDate, "limit", limit, "offset", offset)
	var songs []model.Song

	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		d.Log.Warn("Unable to acquire a database connection", sl.Err(err))
		return nil, err
	}
	defer conn.Release()

	listCmd := `SELECT * from song`

	conditions := make([]string, 0, 5)
	args := make([]any, 0, 7)

	if filter.Group != "" {
		conditions = append(conditions, fmt.Sprintf("group_name = $%d", len(args)+1))
		args = append(args, filter.Group)
	}

	if filter.Title != "" {
		conditions = append(conditions, fmt.Sprintf("title = $%d", len(args)+1))
		args = append(args, filter.Title)
	}

	if filter.ReleaseDate != "" {
		conditions = append(conditions, fmt.Sprintf("release_date = $%d", len(args)+1))
		args = append(args, filter.ReleaseDate)
	}

	if filter.Text != "" {
		conditions = append(conditions, fmt.Sprintf("text = $%d", len(args)+1))
		args = append(args, filter.Text)
	}

	if filter.Link != "" {
		conditions = append(conditions, fmt.Sprintf("link = $%d", len(args)+1))
		args = append(args, filter.Link)
	}

	var conditionPart string
	if len(conditions) > 0 {
		condition := strings.Join(conditions, " AND ")
		conditionPart = fmt.Sprintf(" WHERE %s", condition)
		fmt.Println(condition, conditionPart)
	}

	var limitPart string
	if limit > 0 {
		limitPart = fmt.Sprintf(" LIMIT $%d", len(args)+1)
		args = append(args, limit)
	}

	if offset > 0 {
		limitPart += fmt.Sprintf(" OFFSET $%d", len(args)+1)
		args = append(args, offset)
	}

	cmd := listCmd + conditionPart + limitPart

	d.Log.Debug("SQL command", "command", cmd, "args", args)

	rows, err := conn.Query(ctx, cmd, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var releaseDate time.Time
	for rows.Next() {
		var song model.Song
		if err = rows.Scan(
			&song.ID,
			&song.Group,
			&song.Title,
			&releaseDate,
			&song.Text,
			&song.Link,
		); err != nil {
			return nil, fmt.Errorf("GET SONGS DB : %w", err)
		}

		song.ReleaseDate = releaseDate.Format(model.DateFormat)
		songs = append(songs, song)
	}

	if len(songs) == 0 {
		return nil, errors.New("no songs found")
	}

	return songs, nil
}

func (d *Database) GetSongText(ctx context.Context, id, limit, offset int) (string, error) {

	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return "", err
	}
	defer conn.Release()

	args := make([]any, 0, 3)
	args = append(args, id)

	query := `SELECT text FROM song WHERE id = $1`

	var text string
	err = conn.QueryRow(ctx, query, id).Scan(&text)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("unable to find song by id")
		}
		return "", fmt.Errorf("SELECT text song : %w", err)
	}

	return text, nil
}

func (d *Database) DeleteSong(ctx context.Context, id int) error {

	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
		return err
	}
	defer conn.Release()

	query := `DELETE FROM song WHERE id = $1`
	result, err := conn.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("DELETE song: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("song with id %d not found", id)
	}

	return nil
}

func (d *Database) UpdateSong(ctx context.Context, id int, update model.Song) (model.Song, error) {

	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		d.Log.Warn("Unable to acquire a database connection", sl.Err(err))
		return model.Song{}, err
	}
	defer conn.Release()

	query := `UPDATE song
			  SET group_name   = COALESCE($1, group_name),
	              title        = COALESCE($2, title),
	              release_date = COALESCE($3, release_date),
	              text         = COALESCE($4, text),
	              link         = COALESCE($5, link)
			  WHERE id = $6
			  RETURNING id, group_name, title, release_date, text, link`

	var updatedSong model.Song
	var song model.UpdateSong

	if update.Group != "" {
		song.Group = &update.Group
	}
	if update.Title != "" {
		song.Title = &update.Title
	}
	if update.ReleaseDate != "" {
		song.ReleaseDate = &update.ReleaseDate
	}
	if update.Text != "" {
		song.Text = &update.Text
	}
	if update.Link != "" {
		song.Link = &update.Link
	}

	var releaseDate time.Time
	err = conn.QueryRow(ctx, query, song.Group, song.Title, song.ReleaseDate, song.Text, song.Link, id).
		Scan(
			&updatedSong.ID,
			&updatedSong.Group,
			&updatedSong.Title,
			&releaseDate,
			&updatedSong.Text,
			&updatedSong.Link,
		)

	updatedSong.ReleaseDate = releaseDate.Format(model.DateFormat)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Song{}, fmt.Errorf("unable to find song by id")
		}
		return model.Song{}, fmt.Errorf("UPDATE song : %w", err)
	}

	return updatedSong, nil
}

func (d *Database) AddSong(ctx context.Context, song model.Song) (model.Song, error) {

	conn, err := d.Client.Acquire(ctx)
	if err != nil {
		d.Log.Warn("Unable to acquire a database connection", sl.Err(err))
		return model.Song{}, err
	}
	defer conn.Release()

	query := `
		INSERT INTO song (group_name,title,release_date,text,link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id,group_name,title,release_date,text,link
		`

	row := conn.QueryRow(ctx, query, song.Group, song.Title, song.ReleaseDate, song.Text, song.Link)

	var newSong model.Song
	var date time.Time
	err = row.Scan(
		&newSong.ID,
		&newSong.Group,
		&newSong.Title,
		&date,
		&newSong.Text,
		&newSong.Link,
	)

	newSong.ReleaseDate = date.Format(model.DateFormat)
	if err != nil {
		return model.Song{}, fmt.Errorf("INSERT song: %w", err)
	}

	return newSong, nil
}

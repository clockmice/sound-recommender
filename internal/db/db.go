package db

import (
	"database/sql"
)

type Service interface {
}

var _ Service = (*service)(nil)

type service struct {
	name string
}

func New(name string) (Service, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Ping the database to ensure it's accessible
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS playlists (id TEXT PRIMARY KEY, name TEXT)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS sounds (id TEXT PRIMARY KEY, name TEXT)")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS playlist_sounds (playlist_id INTEGER, sound_id)")
	if err != nil {
		return nil, err
	}

	return db, nil
}

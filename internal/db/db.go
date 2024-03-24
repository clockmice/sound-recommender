package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	api "github.com/clockmice/sound-recommender/gen"
	_ "github.com/mattn/go-sqlite3"
)

type Service interface {
	CreateSound(ctx context.Context, sounds []api.Sound) error
	CreatePlaylist() error
}

var _ Service = (*service)(nil)

type service struct {
	db *sql.DB
}

func New(name string) (Service, error) {
	willCreateDatabase := false

	_, err := os.Stat(name)
	willCreateDatabase = err != nil

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	if willCreateDatabase {
		createDB(db)
	}

	// // Ping the database to ensure it's accessible
	// err = db.Ping()
	// if err != nil {
	// 	return nil, err
	// }
	return &service{db: db}, nil
}

func (s *service) CreateSound(ctx context.Context, sounds []api.Sound) error {
	now := time.Now().UTC()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}

	for _, sound := range sounds {

		res1, err := tx.Exec("INSERT INTO sounds (title, bpm, duration_in_seconds, createdAt) VALUES (?, ?, ?, ?)", sound.Title, sound.Bpm, sound.DurationInSeconds, now)
		if err != nil {
			fmt.Printf("%v", err)
			return err
		}

		for _, genre := range sound.Genres {
			res2, err := tx.Exec("INSERT INTO genres(name) VALUES(?) ON CONFLICT DO NOTHING RETURNING id", genre)
			if err != nil {
				fmt.Printf("%v", err)
				return err
			}

			fmt.Printf("%v", res2)

			sound_id, _ := res1.LastInsertId()
			genre_id, _ := res2.LastInsertId()

			_, err = tx.Exec("INSERT INTO sound_genres VALUES(?, ?) ON CONFLICT DO NOTHING", sound_id, genre_id)
			if err != nil {
				fmt.Printf("%v", err)
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreatePlaylist() error {
	/*
		_, err = db.Exec("INSERT INTO users (name) VALUES (?)", "John Doe")
		if err != nil {
			log.Fatal(err)
		}
	*/
	return nil
}

func createDB(db *sql.DB) {
	db.Exec(`
		CREATE TABLE sounds (
			id string,
			title string,
        	bpm integer,
	        duration_in_seconds integer,
	        createdAt datetime,
        	updatedAt datetime DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		);
		CREATE TABLE genres (name string NOT NULL, PRIMARY KEY (name), UNIQUE(name));
		CREATE TABLE playlists (
			id string,
			name string,
			created_at datetime,
			updated_at datetime DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		);
		CREATE TABLE sound_genres (
			sound_id string REFERENCES sounds(id) ON DELETE CASCADE,
			genre_id string REFERENCES genres(name) ON DELETE CASCADE
		);
		CREATE TABLE playlist_sounds (
			playlist_id string REFERENCES playlists(id) ON DELETE CASCADE,
			sound_id string REFERENCES sounds(id) ON DELETE CASCADE
		);
		CREATE TABLE credits(
			id string,
			sound_id string REFERENCES sounds(id) ON DELETE CASCADE,
			credit_name string,
			credit_role string,
			PRIMARY KEY (id)
		);
	`)
}

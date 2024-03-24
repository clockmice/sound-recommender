package service

import (
	"context"
	"time"

	api "github.com/clockmice/sound-recommender/gen"
	"github.com/clockmice/sound-recommender/internal/db"
	"github.com/google/uuid"
)

type RestController struct {
	dbClient db.Service
}

var _ api.StrictServerInterface = (*RestController)(nil)

// Obs, not a thread safe map
var soundsMap = make(map[string]api.Sound)
var playlistsMap = make(map[string]api.Playlist)

func (r RestController) GetSounds(ctx context.Context, request api.GetSoundsRequestObject) (api.GetSoundsResponseObject, error) {
	var resp = []api.Sound{}

	for _, sound := range soundsMap {
		resp = append(resp, sound)
	}

	return api.GetSounds200JSONResponse{Data: &resp}, nil
}

func (r RestController) PostAdminSounds(ctx context.Context, request api.PostAdminSoundsRequestObject) (api.PostAdminSoundsResponseObject, error) {
	var created = []api.Sound{}
	data := request.Body.Data
	now := time.Now().UTC()

	for _, item := range data {
		id := uuid.New().String()

		sound := api.Sound{
			Bpm:               *item.Bpm,
			CreatedAt:         now,
			Credits:           *item.Credits,
			DurationInSeconds: *item.DurationInSeconds,
			Genres:            *item.Genres,
			Id:                id,
			Title:             *item.Title,
			UpdatedAt:         now,
		}
		soundsMap[id] = sound
		created = append(created, sound)
	}

	return api.PostAdminSounds201JSONResponse{Data: &created}, nil
}

func (r RestController) PostPlaylists(ctx context.Context, request api.PostPlaylistsRequestObject) (api.PostPlaylistsResponseObject, error) {
	var created = []api.Playlist{}
	data := request.Body.Data
	now := time.Now().UTC()

	for _, item := range data {
		id := uuid.New().String()
		var sounds []api.Sound

		for _, id := range *item.Sounds {
			sound, exists := soundsMap[id]
			if !exists {
				return api.PostPlaylists400JSONResponse{N400JSONResponse: api.N400JSONResponse{Detail: "Sound Id not found"}}, nil
			}

			sounds = append(sounds, sound)
		}

		playlist := api.Playlist{
			CreatedAt: &now,
			Id:        &id,
			Sounds:    &sounds,
			Title:     item.Title,
			UpdatedAt: &now,
		}
		playlistsMap[id] = playlist
		created = append(created, playlist)
	}

	return api.PostPlaylists201JSONResponse{Data: &created}, nil
}

func (r RestController) GetSoundsRecommended(ctx context.Context, request api.GetSoundsRecommendedRequestObject) (api.GetSoundsRecommendedResponseObject, error) {
	playlistId := request.Params.PlaylistId
	var genrePopular = []string{}

	for _, sound := range *playlistsMap[playlistId].Sounds {
		genrePopular = append(genrePopular, sound.Genres...)
	}

}

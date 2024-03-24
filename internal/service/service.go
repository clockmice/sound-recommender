package service

import (
	"context"
	"sort"
	"time"

	api "github.com/clockmice/sound-recommender/gen"
	"github.com/google/uuid"
)

type RestController struct{}

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

func (r RestController) GetPlaylists(ctx context.Context, request api.GetPlaylistsRequestObject) (api.GetPlaylistsResponseObject, error) {
	var resp = []api.Playlist{}

	for _, playlist := range playlistsMap {
		resp = append(resp, playlist)
	}

	return api.GetPlaylists200JSONResponse{Data: &resp}, nil
}

func (r RestController) GetSoundsRecommended(ctx context.Context, request api.GetSoundsRecommendedRequestObject) (api.GetSoundsRecommendedResponseObject, error) {
	playlistId := request.Params.PlaylistId
	var genresInPlaylist = []api.Genre{}
	var soundsInPlaylist = make(map[string]bool)

	for _, sound := range *playlistsMap[playlistId].Sounds {
		genresInPlaylist = append(genresInPlaylist, sound.Genres...)
		soundsInPlaylist[sound.Id] = true
	}

	topGenres := getTopGenres(genresInPlaylist)
	recommendations := getRecommendations(topGenres, soundsInPlaylist)

	return api.GetSoundsRecommended200JSONResponse{Data: &recommendations}, nil
}

func getTopGenres(genres []api.Genre) []api.Genre {
	var topGenres []api.Genre
	// Stores count of each genre
	genreCounts := make(map[string]int)

	// Count the occurrences of each genre
	for _, genre := range genres {
		genreCounts[genre]++
	}

	// Convert map to slice of genre-count pairs
	var genreCountPairs []struct {
		genre api.Genre
		count int
	}
	for genre, count := range genreCounts {
		genreCountPairs = append(genreCountPairs, struct {
			genre api.Genre
			count int
		}{genre, count})
	}

	// Sort the slice by count in descending order
	sort.Slice(genreCountPairs, func(i, j int) bool {
		return genreCountPairs[i].count > genreCountPairs[j].count
	})

	for _, pair := range genreCountPairs {
		topGenres = append(topGenres, pair.genre)
		if len(topGenres) == 3 {
			break
		}
	}

	return topGenres
}

func getRecommendations(topGenres []api.Genre, soundsInPlaylist map[string]bool) []api.Sound {
	var recommendations = []api.Sound{}

	for id, sound := range soundsMap {
		_, exists := soundsInPlaylist[id]
		if exists {
			continue
		}

		for _, genre := range sound.Genres {
			for _, i := range topGenres {
				if i == genre {
					recommendations = append(recommendations, sound)
				}
			}
		}
	}

	return recommendations
}

package service

import (
	"context"

	api "github.com/clockmice/sound-recommender/gen"
)

type RestController struct{}

var _ api.StrictServerInterface = (*RestController)(nil)
var maxPageSize = 3

// Obs, not a thread safe map
var assets = make(map[string]api.Sound)

func (r RestController) GetSounds(ctx context.Context, request api.GetSoundsRequestObject) (api.GetSoundsResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (r RestController) PostAdminSounds(ctx context.Context, request api.PostAdminSoundsRequestObject) (api.PostAdminSoundsResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (r RestController) PostPlaylists(ctx context.Context, request api.PostPlaylistsRequestObject) (api.PostPlaylistsResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (r RestController) GetSoundsRecommended(ctx context.Context, request api.GetSoundsRecommendedRequestObject) (api.GetSoundsRecommendedResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

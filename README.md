## Sound Recommender API
### Overview
Go-based web service containerized in Docker.
### How to run:
Make sure you have docker running in your system and `docker-compose` installed. Then run:
``` shell
docker-compose up
``````
This will start the service at `http://localhost:8080` so the Postman test collection should run as it is.

### Endpoints:
See spec/api.yaml for OpenAPI definitions.

- `POST /admin/sounds` - Create sounds.
- `GET /sounds` - List all sounds.
- `POST /playlists` - Create a new playlist.
- `GET /playlists` - List all playlists.
- `GET /sounds/recommended/{playlistId}` - Get recommended sounds for playlist ID based on genre.

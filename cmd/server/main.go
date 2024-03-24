package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	api "github.com/clockmice/sound-recommender/gen"
	db "github.com/clockmice/sound-recommender/internal/db"
	svc "github.com/clockmice/sound-recommender/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	port := flag.String("port", "8080", "Port for HTTP server")
	flag.Parse()

	dbClient, err := db.New("sound_recommender")
	if err != nil {
		log.Fatal(err)
	}

	s := createServer(*port, dbClient)

	log.Fatal(s.ListenAndServe())
}

func createServer(port string, dbClient db.Service) *http.Server {
	swag, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("failed to load swagger spec, err: %v", err)
	}
	swag.Servers = nil

	r := chi.NewRouter()
	r.Use(oapimiddleware.OapiRequestValidator(swag))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handler := api.NewStrictHandler(&svc.RestController{dbClient: dbClient}, nil)
	api.HandlerFromMux(handler, r)

	return &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("127.0.0.1", port),
	}
}

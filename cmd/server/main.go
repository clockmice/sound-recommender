package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	api "github.com/clockmice/sound-recommender/gen"
	svc "github.com/clockmice/sound-recommender/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	port := flag.String("port", "8080", "Port for HTTP server")
	flag.Parse()

	s := createServer(*port)

	log.Fatal(s.ListenAndServe())
}

func createServer(port string) *http.Server {
	swag, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("failed to load swagger spec, err: %v", err)
	}
	swag.Servers = nil

	r := chi.NewRouter()
	r.Use(oapimiddleware.OapiRequestValidator(swag))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handler := api.NewStrictHandler(&svc.RestController{}, nil)
	api.HandlerFromMux(handler, r)

	return &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("0.0.0.0", port),
	}
}

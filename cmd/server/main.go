package main

import (
	"context"
	"log"
	"net/http"

	"coinflip/internal/config"
	"coinflip/internal/db"
	"coinflip/internal/game"
	"coinflip/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	if err := db.Init(ctx, cfg.DBUrl); err != nil {
		log.Fatal(err)
	}
	if err := game.LoadRTP(ctx); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API endpoints
	r.Post("/api/play", handlers.PlayHandler)
	r.Get("/api/balance", handlers.BalanceHandler)
	r.Get("/api/history", handlers.HistoryHandler)

	// Admin endpoints
	r.Route("/api/admin", func(r chi.Router) {
		r.Get("/rtp", handlers.GetRTPHandler)
		r.Post("/rtp", handlers.SetRTPHandler)
	})

	// Статика клиента
	fs := http.FileServer(http.Dir("./web/client"))
	r.Handle("/*", fs)

	log.Println("listening on :" + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
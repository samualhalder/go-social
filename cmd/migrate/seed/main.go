package main

import (
	"log"

	"github.com/samualhalder/go-social/internal/db"
	"github.com/samualhalder/go-social/internal/env"
	"github.com/samualhalder/go-social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgresql://samualhalder:samualpass@localhost:5433/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	store := store.NewStore(conn)
	db.Seed(store)
}

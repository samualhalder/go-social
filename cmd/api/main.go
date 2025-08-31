package main

import (
	"log"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/samualhalder/go-social/internal/db"
	"github.com/samualhalder/go-social/internal/env"
	"github.com/samualhalder/go-social/internal/store"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while laoding .env")
	}

	cnf := config{addr: env.GetString("ADDR", ":8080"), db: dbConfig{
		addr:        env.GetString("DB_ADDR", "postgresql://samualhalder:samualpass@localhost:5433/social?sslmode=disable"),
		maxOpenConn: env.GetInt("MAX_OPEN_CONN", 30),
		maxIdleConn: env.GetInt("MAX_IDLE_CONN", 30),
		maxIdleTime: env.GetString("MAX_IDLE_TIME", "15m"),
	}}
	db, err := db.New(cnf.db.addr, cnf.db.maxOpenConn, cnf.db.maxIdleConn, cnf.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	slog.Info("🗃️ DB connection is stablished")
	store := store.NewStore(db)
	app := application{config: cnf, store: store}
	mux := app.mount()
	slog.Info("🛣️ Route setup is done")
	log.Fatal(app.run(mux))
}

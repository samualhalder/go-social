// @title Go Social API
// @version 1.0
// @description This is the API documentation for Go Social project.
// @termsOfService http://swagger.io/terms/

// @contact.name Samual Halder
// @contact.email your-email@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/samualhalder/go-social/docs"
	"github.com/samualhalder/go-social/internal/db"
	"github.com/samualhalder/go-social/internal/env"
	"github.com/samualhalder/go-social/internal/store" // swagger docs
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while laoding .env")
	}

	cnf := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:        env.GetString("DB_ADDR", "postgresql://samualhalder:samualpass@localhost:5433/social?sslmode=disable"),
			maxOpenConn: env.GetInt("MAX_OPEN_CONN", 30),
			maxIdleConn: env.GetInt("MAX_IDLE_CONN", 30),
			maxIdleTime: env.GetString("MAX_IDLE_TIME", "15m"),
		},
		mail: mailConfig{
			exp: time.Hour * 24 * 3,
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(cnf.db.addr, cnf.db.maxOpenConn, cnf.db.maxIdleConn, cnf.db.maxIdleTime)
	if err != nil {
		logger.Panic(err)
	}
	defer db.Close()
	logger.Info("🗃️ DB connection is stablished")
	store := store.NewStore(db)
	app := application{config: cnf, store: store, logger: logger}
	mux := app.mount()
	logger.Info("🛣️ Route setup is done")
	logger.Fatal(app.run(mux))
}

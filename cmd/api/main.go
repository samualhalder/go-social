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

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/samualhalder/go-social/docs"
	"github.com/samualhalder/go-social/internal/auth"
	"github.com/samualhalder/go-social/internal/db"
	"github.com/samualhalder/go-social/internal/env"
	"github.com/samualhalder/go-social/internal/mailer"
	"github.com/samualhalder/go-social/internal/store" // swagger docs
	"github.com/samualhalder/go-social/internal/store/cache"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while laoding .env")
	}

	cnf := config{
		addr:        env.GetString("ADDR", ":8080"),
		env:         env.GetString("ENV", "development"),
		frontEndURL: env.GetString("FRONT_END_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:        env.GetString("DB_ADDR", "postgresql://samualhalder:samualpass@localhost:5433/social?sslmode=disable"),
			maxOpenConn: env.GetInt("MAX_OPEN_CONN", 30),
			maxIdleConn: env.GetInt("MAX_IDLE_CONN", 30),
			maxIdleTime: env.GetString("MAX_IDLE_TIME", "15m"),
		},
		mail: mailConfig{
			exp:      time.Hour * 24 * 3,
			fromUser: env.GetString("FROM_USER", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				username: "sam",
				pass:     "sam",
			},
			token: tokenConfig{
				secret: env.GetString("JWT_SECRET", "itsasecretok2323"),
				expiry: time.Hour * 3 * 24,
				issuer: env.GetString("TOKEN_ISSUER", "GO_SOCIAL"),
			},
		},
		redisCfg: RedisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PASSWORD", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", true),
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

	var rdb *redis.Client

	if cnf.redisCfg.enabled {
		rdb = cache.NewRedisClient(cnf.redisCfg.addr, cnf.redisCfg.pw, cnf.redisCfg.db)
		logger.Info("🗄️ redis cache connection is stablished")
	}
	cacheStore := cache.NewRedisStore(rdb)
	rdb.SetEX(rdb.Context(), "user-1", "test", time.Hour)
	app := application{
		config: cnf,
		store:  store, logger: logger,
		cacheStorage:  cacheStore,
		mailer:        mailer.NewSendGrid(cnf.mail.fromUser, cnf.mail.sendGrid.apiKey),
		authenticator: auth.NewJWTAuthenticator(cnf.auth.token.secret, cnf.auth.token.issuer, cnf.auth.token.issuer),
	}
	mux := app.mount()
	logger.Info("🛣️ Route setup is done")
	logger.Fatal(app.run(mux))
}

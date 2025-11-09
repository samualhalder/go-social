package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/samualhalder/go-social/internal/auth"
	"github.com/samualhalder/go-social/internal/mailer"
	"github.com/samualhalder/go-social/internal/store"
	"github.com/samualhalder/go-social/internal/store/cache"
	"go.uber.org/zap"

	httpSwagger "github.com/swaggo/http-swagger"
)

type application struct {
	config        config
	store         store.Store
	cacheStorage  cache.Store
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator *auth.JWTAuthenticator
}

type config struct {
	addr        string
	db          dbConfig
	mail        mailConfig
	env         string
	frontEndURL string
	auth        authConfig
	redisCfg    RedisConfig
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}
type tokenConfig struct {
	secret string
	expiry time.Duration
	issuer string
}
type basicConfig struct {
	username string
	pass     string
}

type dbConfig struct {
	addr        string
	maxOpenConn int
	maxIdleConn int
	maxIdleTime string
}
type mailConfig struct {
	exp      time.Duration
	fromUser string
	sendGrid sendGridConfig
}
type sendGridConfig struct {
	apiKey string
}
type RedisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	fmt.Printf("hti here")
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheck)
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/create", app.createPost)
			r.Route("/{postId}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.checkPostOwnerShip("admin", app.deletePostById))
				r.Patch("/", app.checkPostOwnerShip("moderator", app.updatePostById))
			})
		})
		r.Route("/comments", func(r chi.Router) {
			r.Post("/create/{postId}", app.createComment)
		})
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHanlder)
			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.GetFeedForUser)
			})
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)
				r.Post("/follow", app.followUserHandler)
				//TODO: will make it delete req when we add authintication via tokens
				r.Put("/unfollow", app.unFollowUserHandler)
			})
		})
		r.Route("/authinticate", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})

	})
	return r
}

func (app *application) run(mux http.Handler) error {
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second,
	}
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Info("signal caught", "signal", s.String())
		shutdown <- srv.Shutdown(ctx)

	}()
	app.logger.Info("🌐Server is running on", slog.String("port", srv.Addr))
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Info("server has stopped", "addr", app.config.addr)
	return nil
}

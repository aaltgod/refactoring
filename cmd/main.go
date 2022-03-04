package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aaltgod/refactoring/internal/apperror"
	handler "github.com/aaltgod/refactoring/internal/handlers"
	"github.com/aaltgod/refactoring/internal/user/repository"
	"github.com/aaltgod/refactoring/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

const store = `users.json`

func main() {

	logger := logger.NewLogger()

	server := &http.Server{Addr: "0.0.0.0:3333", Handler: service(logger)}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Fatal(err)
		}
		serverStopCtx()
	}()

	logger.Info("Starting service")
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
	logger.Info("Stoping service")

	<-serverCtx.Done()
}

func service(logger *logrus.Logger) http.Handler {

	r := chi.NewRouter()

	storage := repository.NewRepository(store, logger)
	h := handler.NewHandler(storage, logger)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", apperror.Middleware(h.SearchUsers))
				r.Post("/", apperror.Middleware(h.CreateUser))

				r.Route("/{id}", func(r chi.Router) {
					r.Use(UserCtx)
					r.Get("/", apperror.Middleware(h.GetUser))
					r.Patch("/", apperror.Middleware(h.UpdateUser))
					r.Delete("/", apperror.Middleware(h.DeleteUser))
				})
			})
		})
	})

	return r
}

func UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id := chi.URLParam(r, "id"); id != "" {
			ctx := context.WithValue(r.Context(), "id", id)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		render.Render(w, r, apperror.ErrUserNotFound)
	})
}

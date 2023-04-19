package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"refactoring/internal/users/delivery/httpchi"
	"refactoring/internal/users/repository/ujson"
	"refactoring/internal/users/usecase"

	"github.com/go-chi/chi/v5"
)

const store = `../users.json`

func main() {
	f, err := os.Open(store)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	router := chi.NewRouter()

	userRepo := ujson.NewUserRepository(f)
	uu := usecase.NewUserUsecase(userRepo)
	httpchi.NewUserHandler(router, uu)

	log.Println("starting server at :3333")
	srv := &http.Server{
		Addr:    ":3333",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("error listen and serve: %s", err.Error())
		}
	}()

	// wait signal to shutdown server with a timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server. ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err.Error())
	}

	log.Println("Server exiting")

}

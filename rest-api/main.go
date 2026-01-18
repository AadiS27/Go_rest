package main

import (
	"context"
	// "fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AadiS27/Go_rest/internal/config"
	"github.com/AadiS27/Go_rest/internal/http/handlers/student"
	"github.com/AadiS27/Go_rest/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database connection

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	slog.Info("connected to database", slog.String("database", cfg.StoragePath))
	defer storage.Db.Close()

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetStudent(storage))
	
	//setup servver
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	slog.Info("server is starting", slog.String("address", cfg.Address))
	// fmt.Println("server is running on port",cfg.Address)

	//graceful shutdown
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}

	}()

	<-done

	slog.Info("server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //this is used for graceful shutdown where context is used to cancel the server shutdown
	defer cancel()                                                          //this is used to cancel the context if the server is not shutdown within the timeout
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown complete")
}

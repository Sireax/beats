package main

import (
	"beats/api"
	"beats/db"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db.Connect("postgres://backend:supersecret@127.0.0.1:6432/beats")
	a := api.NewAPI("localhost:8080")
	go a.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down...")
}

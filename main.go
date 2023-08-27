package main

import (
	"os"
	"os/signal"

	"bibit.id/challenge/handler"
	"bibit.id/challenge/repo"
	"bibit.id/challenge/usecase"
)

func main() {
	stockRepo := repo.New()
	stockUsecase := usecase.New(stockRepo)
	stockHandler := handler.New(stockUsecase)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go serveGRPC(stockHandler)
	go serveKafka(stockHandler)

	<-signals
}

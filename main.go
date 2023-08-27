package main

import (
	"bibit.id/challenge/handler"
	"bibit.id/challenge/repo"
	"bibit.id/challenge/usecase"
)

func main() {
	stockRepo := repo.New()
	stockUsecase := usecase.New(stockRepo)
	stockHandler := handler.New(stockUsecase)

	go serveGRPC(stockHandler)
	serveKafka(stockHandler)
}

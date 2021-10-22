package main

import (
	"context"
	"homework6/internal/http"
	"homework6/internal/store/inmemory"
	"log"
)

func main(){
	store := inmemory.NewDB()
	srv := http.NewServer(context.Background(), ":8080", store)

	if err := srv.Run(); err !=nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}
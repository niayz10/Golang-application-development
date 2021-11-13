package main

import (
	"context"
	"project/internal/http"
	"project/internal/store/postgres"
)

func main(){
	urlExample := "postgres://goproject:goproject@localhost:5432/goproject"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	srv := http.NewServer(context.Background(), ":8080", store)
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
}
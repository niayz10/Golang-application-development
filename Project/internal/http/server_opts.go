package http

import (
	"github.com/go-redis/redis/v8"
	"project/internal/store"
)

type ServerOption func(srv *Server)


func WithAddress(address string) ServerOption{
	return func(srv *Server) {
		srv.Address = address
	}
}

func WithStore(store store.Store) ServerOption{
	return func(srv *Server){
		srv.store = store
	}
}

func WithRedis(redis *redis.Client) ServerOption {
	return func(srv *Server) {
		srv.redis = redis
	}
}
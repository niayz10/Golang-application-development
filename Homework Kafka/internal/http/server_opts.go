package http

import (
	"homework-kafka/internal/message_broker"
	"homework-kafka/internal/store"
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

func WithBroker(broker message_broker.MessageBroker) ServerOption {
	return func(srv *Server) {
		srv.broker = broker
	}
}
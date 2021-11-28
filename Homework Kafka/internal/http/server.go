package http

import (
	"context"
	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
	"homework-kafka/internal/message_broker"
	"homework-kafka/internal/store"
	"log"
	"net/http"
	"time"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store
	cache       *lru.TwoQueueCache
	broker      message_broker.MessageBroker
	Address     string
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()
	mr := NewMangaResource(s.store.Mangas(), s.broker, s.cache)
	br := NewBookResource(s.store.Books(), s.broker, s.cache)
	r.Mount("/manga", mr.Routes())
	r.Mount("/book",br.Routes())

	return r
}

func (s *Server) Run() error {

	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on ", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}
	log.Println("[HTTP] Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
}

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-ozzo/ozzo-validation/v4"
	"log"
	"net/http"
	"project/internal/models"
	"project/internal/store"
	"strconv"
	"time"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	store       store.Store
	Address     string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
		store:       store,
		
		Address:     address,	
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.Post("/allmanga", func(w http.ResponseWriter, r *http.Request) {
		manga := new(models.Manga)
		if err := json.NewDecoder(r.Body).Decode(manga); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Mangas().Create(r.Context(), manga); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	r.Get("/allmanga", func(w http.ResponseWriter, r *http.Request) {
		mangas, err := s.store.Mangas().All(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, mangas)
	})

	r.Get("/allmanga/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		manga, err := s.store.Mangas().ByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, manga)
	})

	r.Put("/allmanga", func(w http.ResponseWriter, r *http.Request) {
		manga := new(models.Manga)
		if err := json.NewDecoder(r.Body).Decode(manga); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		err := validation.ValidateStruct(
			manga,
			validation.Field(&manga.ID, validation.Required),
			validation.Field(&manga.Title, validation.Required),
		)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		if err := s.store.Mangas().Update(r.Context(), manga); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	r.Delete("/allmanga/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		
		if err := s.store.Mangas().Delete(r.Context(), id); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	r.Post("/books", func(w http.ResponseWriter, r *http.Request) {
		book := new(models.Book)
		if err := json.NewDecoder(r.Body).Decode(book); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Books().Create(r.Context(), book); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	r.Get("/books", func(w http.ResponseWriter, r *http.Request) {
		books, err := s.store.Books().All(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, books)
	})

	r.Get("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		book, err := s.store.Books().ByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, book)
	})

	r.Put("/books", func(w http.ResponseWriter, r *http.Request) {
		book := new(models.Book)
		if err := json.NewDecoder(r.Body).Decode(book); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		err := validation.ValidateStruct(
			book,
			validation.Field(&book.ID, validation.Required),
			validation.Field(&book.Title, validation.Required),
		)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		if err := s.store.Books().Update(r.Context(), book); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

	r.Delete("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}
		
		if err := s.store.Books().Delete(r.Context(), id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB err: %v", err)
			return
		}
	})

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

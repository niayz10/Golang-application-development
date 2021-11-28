package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	lru "github.com/hashicorp/golang-lru"
	"homework-kafka/internal/message_broker"
	"homework-kafka/internal/models"
	"homework-kafka/internal/store"
	"net/http"
	"strconv"
)

type BookResource struct {
	repo store.BooksRepository
	broker message_broker.MessageBroker
	cache  *lru.TwoQueueCache
}

func NewBookResource(repo store.BooksRepository, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *BookResource {
	return &BookResource{
		repo: repo,
		broker: broker,
		cache:  cache,
	}
}


func (b *BookResource) Routes() chi.Router {
	r :=chi.NewRouter()

	r.Post("/", b.CreateBook)
	r.Get("/", b.AllBooks)
	r.Get("/{id}", b.ByID)
	r.Put("/", b.UpdateBook)
	r.Delete("/{id}", b.DeleteBook)
	return r
}

func (b *BookResource) CreateBook(w http.ResponseWriter, r *http.Request) {
	book := new(models.Book)
	if err := json.NewDecoder(r.Body).Decode(book); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := b.repo.Create(r.Context(), book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (b *BookResource) AllBooks(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.BookFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		titles, ok := b.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, titles)
			return
		}
		filter.Query = &searchQuery
	}
	titles, err := b.repo.All(r.Context(), filter)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if searchQuery != "" {
		b.cache.Add(searchQuery, titles)
	}
	render.JSON(w, r, titles)
}

func (b *BookResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	book, err := b.repo.ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	render.JSON(w, r, book)
}

func (b *BookResource) UpdateBook(w http.ResponseWriter, r *http.Request) {
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
	if err := b.repo.Update(r.Context(), book); err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
}

func (b *BookResource) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := b.repo.Delete(r.Context(), id); err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
}

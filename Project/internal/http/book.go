package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-redis/redis/v8"
	"net/http"
	"project/internal/models"
	"project/internal/store"
	"strconv"
)

type BookResource struct {
	repo store.BooksRepository
	red *redis.Client
}

func NewBookResource(repo store.BooksRepository, red *redis.Client) *BookResource {
	return &BookResource{
		repo: repo,
		red: red,
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

	searchQuery :=queryValues.Get("query")
	if searchQuery != "" {
		booksFromRed, err := b.red.Get(context.Background(), searchQuery).Result()
		fmt.Printf("redis = %s\n", booksFromRed)
		if err == nil {
			books := make([]*models.Book, 0)
			err := json.Unmarshal([]byte(booksFromRed), &books)
			if err != nil {
				return
			}
			render.JSON(w, r, books)
			return
		}
		filter.Query =&searchQuery
	}

	books, err := b.repo.All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if searchQuery != "" {
		fmt.Println(searchQuery)
		booksMarshal, _:= json.Marshal(books)
		b.red.Set(context.Background(), searchQuery, booksMarshal, 0)
	}

	render.JSON(w, r, books)
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
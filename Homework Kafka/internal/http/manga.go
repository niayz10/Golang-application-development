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

type MangaResource struct {
	repo store.MangasRepository
	broker message_broker.MessageBroker
	cache  *lru.TwoQueueCache
}

func NewMangaResource(repo store.MangasRepository, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *MangaResource {
	return &MangaResource{
		repo: repo,
		broker: broker,
		cache:  cache,
	}
}


func (m *MangaResource) Routes() chi.Router {
	r :=chi.NewRouter()

	r.Post("/", m.CreateManga)
	r.Get("/", m.AllManga)
	r.Get("/{id}", m.ByID)
	r.Put("/", m.UpdateManga)
	r.Delete("/{id}", m.DeleteManga)

	return r
}

func (m *MangaResource) CreateManga(w http.ResponseWriter, r *http.Request) {
	manga := new(models.Manga)
	if err := json.NewDecoder(r.Body).Decode(manga); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := m.repo.Create(r.Context(), manga); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (m *MangaResource) AllManga(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.MangaFilter{}

	searchQuery :=queryValues.Get("query")
	if searchQuery != "" {
		titles, ok := m.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, titles)
			return
		}
		filter.Query = &searchQuery
	}
	titles, err := m.repo.All(r.Context(), filter)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	if searchQuery != "" {
		m.cache.Add(searchQuery, titles)
	}
	render.JSON(w, r, titles)
}

func (m *MangaResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	manga, err := m.repo.ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	render.JSON(w, r, manga)
}

func (m *MangaResource) UpdateManga(w http.ResponseWriter, r *http.Request) {
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
	if err := m.repo.Update(r.Context(), manga); err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

}

func (m *MangaResource) DeleteManga(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := m.repo.Delete(r.Context(), id); err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

}

package inmemory

import (
	"context"
	"fmt"
	"homework6/internal/models"
	"homework6/internal/store"
	"sync"
)

type DB struct {
	data map[int]*models.Manga

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		data: make(map[int]*models.Manga),
		mu:   new(sync.RWMutex),
	}
}

func (db *DB) Create(ctx context.Context, manga *models.Manga) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[manga.ID] = manga
	return nil
}
func (db *DB) All(ctx context.Context) ([]*models.Manga, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	mangas := make([]*models.Manga, 0, len(db.data))
	for _, manga := range db.data {
		mangas = append(mangas, manga)
	}
	return mangas, nil
}
func (db *DB) ByID(ctx context.Context, id int) (*models.Manga, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	manga, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("No manga with id %d", id)
	}
	return manga, nil
}
func (db *DB) Update(ctx context.Context, manga *models.Manga) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[manga.ID] = manga
	return nil
}
func (db *DB) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}

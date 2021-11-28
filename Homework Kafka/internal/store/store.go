package store

import (
	"context"
	"homework-kafka/internal/models"
)


type Store interface {
	Connect(url string) error
	Close() error

	Mangas() MangasRepository
	Books() BooksRepository
}

type MangasRepository interface {
	Create(ctx context.Context, manga *models.Manga) error
	All(ctx context.Context, filter *models.MangaFilter) ([]*models.Manga, error)
	ByID(ctx context.Context, id int) (*models.Manga, error)
	Update(ctx context.Context, manga *models.Manga) error
	Delete(ctx context.Context, id int ) error
}

type BooksRepository interface {
	Create(ctx context.Context, book *models.Book) error
	All(ctx context.Context, filter *models.BookFilter) ([]*models.Book, error)
	ByID(ctx context.Context, id int) (*models.Book, error)
	Update(ctx context.Context, book *models.Book) error
	Delete(ctx context.Context, id int ) error
}
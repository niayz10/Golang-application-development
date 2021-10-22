package store

import (
	"context"
	"homework6/internal/models"
)

type Store interface {
	
	Create(ctx context.Context, manga *models.Manga) error
	All(ctx context.Context) ([]*models.Manga, error)
	ByID(ctx context.Context, id int) (*models.Manga, error)
	Update(ctx context.Context, manga *models.Manga) error
	Delete(ctx context.Context, id int) error
	
}

// type MangaRepository interface {
// 	Create(ctx context.Context, manga *models.Manga) error
// 	All(ctx context.Context) ([]*models.Manga, error)
// 	ByID(ctx context.Context, id int) (*models.Manga, error)
// 	Update(ctx context.Context, manga models.Manga) error
// 	Delete(ctx context.Context, id int ) error 
// }
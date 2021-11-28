package postgres
import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"homework-kafka/internal/models"
	"homework-kafka/internal/store"
)



type MangasRepository struct {
	conn *sqlx.DB
}

func NewMangaRepository(conn *sqlx.DB) store.MangasRepository {
	return &MangasRepository{
		conn: conn,
	}
}

func (c MangasRepository) Create(ctx context.Context, manga *models.Manga) error {
	_, err := c.conn.Exec("INSERT INTO manga(title, description, genre, numberofchapters, author, country) " +
		"VALUES ($1, $2, $3, $4, $5, $6)", manga.Title, manga.Description, manga.Genre, manga.NumberOfChapters, manga.Author, manga.Country)
	if err != nil {
		return err
	}

	return nil
}

func (c MangasRepository) All(ctx context.Context, filter *models.MangaFilter) ([]*models.Manga, error) {
	mangas := make([]*models.Manga, 0)
	basicQuery := "SELECT * FROM manga"
	if filter.Query != nil {
		basicQuery = fmt.Sprintf("%s WHERE title ILIKE $1", basicQuery)
		if err := c.conn.Select(&mangas, basicQuery, "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}
		return mangas, nil
	}

	if err := c.conn.Select(&mangas, basicQuery); err != nil {
		return nil, err
	}

	return mangas, nil
}

func (c MangasRepository) ByID(ctx context.Context, id int) (*models.Manga, error) {
	manga := new(models.Manga)
	if err := c.conn.Get(manga, "SELECT id, title, description, genre FROM manga WHERE id=$1", id); err != nil {
		return nil, err
	}

	return manga, nil
}

func (c MangasRepository) Update(ctx context.Context, manga *models.Manga) error {
	_, err := c.conn.Exec("UPDATE manga SET title = $1 WHERE id = $2", manga.Title, manga.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c MangasRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM manga WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

package postgres
import (
	"context"
	"github.com/jmoiron/sqlx"
	"project/internal/models"
	"project/internal/store"
)

func (db *DB) Mangas() store.MangasRepository {
	if db.mangas == nil {
		db.mangas = NewMangaRepository(db.conn)
	}

	return db.mangas
}

type MangasRepository struct {
	conn *sqlx.DB
}

func NewMangaRepository(conn *sqlx.DB) store.MangasRepository {
	return &MangasRepository{conn: conn}
}

func (c MangasRepository) Create(ctx context.Context, manga *models.Manga) error {
	_, err := c.conn.Exec("INSERT INTO manga(title, description, genre, numberofchapters, author, country) VALUES " +
		"($1, $2, $3, $4, $5, $6)", manga.Title, manga.Description, manga.Genre, manga.NumberOfChapters, manga.Author, manga.Country)
	if err != nil {
		return err
	}

	return nil
}

func (c MangasRepository) All(ctx context.Context) ([]*models.Manga, error) {
	categories := make([]*models.Manga, 0)
	if err := c.conn.Select(&categories, "SELECT * FROM manga"); err != nil {
		return nil, err
	}

	return categories, nil
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

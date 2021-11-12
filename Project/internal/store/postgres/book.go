package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"project/internal/models"
	"project/internal/store"
)

func (db *DB) Books() store.BooksRepository {
	if db.books == nil {
		db.books = NewBooksRepository(db.conn)
	}

	return db.books
}

type BooksRepository struct {
	conn *sqlx.DB
}

func NewBooksRepository(conn *sqlx.DB) store.BooksRepository {
	return &BooksRepository{conn: conn}
}

func (c BooksRepository) Create(ctx context.Context, book *models.Book) error {
	_, err := c.conn.Exec("INSERT INTO book(title, description, genre, numberofchapters, author, country) VALUES " +
		"($1, $2, $3, $4, $5, $6)", book.Title, book.Description, book.Genre, book.NumberOfChapters, book.Author, book.Country)
	if err != nil {
		return err
	}

	return nil
}

func (c BooksRepository) All(ctx context.Context) ([]*models.Book, error) {
	books := make([]*models.Book, 0)
	if err := c.conn.Select(&books, "SELECT * FROM book"); err != nil {
		return nil, err
	}

	return books, nil
}

func (c BooksRepository) ByID(ctx context.Context, id int) (*models.Book, error) {
	book := new(models.Book)
	if err := c.conn.Get(book, "SELECT id, title, description, genre FROM book WHERE id=$1", id); err != nil {
		return nil, err
	}

	return book, nil
}

func (c BooksRepository) Update(ctx context.Context, book *models.Book) error {
	_, err := c.conn.Exec("UPDATE book SET title = $1 WHERE id = $2", book.Title, book.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c BooksRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM book WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
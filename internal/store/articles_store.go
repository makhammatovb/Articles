package store

import (
	"database/sql"
	"fmt"
	"time"
)

type Article struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	Image           string         `json:"image"`
	AuthorID        int            `json:"author_id"`
	Paraghraps      []Paraghraph    `json:"paraghraps"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type Paraghraph struct {
	ID              int      `json:"id"`
	Headline        string   `json:"headline"`
	Body            string   `json:"body"`
	OrderIndex      int      `json:"order_index"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type PostgresArticleStore struct {
	db *sql.DB
}

func NewPostgresArticleStore(db *sql.DB) *PostgresArticleStore {
	return &PostgresArticleStore{db: db}
}

type ArticleStore interface {
	CreateArticle(article *Article) (*Article, error)
	GetArticleByID(id int64) (*Article, error)
	UpdateArticle(article *Article) error
	DeleteArticle(id int64) error
}

func (pg *PostgresArticleStore) CreateArticle(article *Article) (*Article, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := 
	`INSERT INTO articles (title, description, image, author_id)
	VALUES ($1, $2, $3, $4) RETURNING id;
	`
	err = tx.QueryRow(query, article.Title, article.Description, article.Image, article.AuthorID).Scan(&article.ID)
	if err != nil {
		return nil, err
	}
	for _, paraghraph := range article.Paraghraps {
		query := 
		`INSERT INTO paraghraps (article_id, headline, body, order_index)
		VALUES ($1, $2, $3, $4) RETURNING id;
		`
		err = tx.QueryRow(query, article.ID, paraghraph.Headline, paraghraph.Body, paraghraph.OrderIndex).Scan(&paraghraph.ID)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (pg *PostgresArticleStore) GetArticleByID(id int64) (*Article, error) {
	article := &Article{}
	return article, nil
}

func (pg *PostgresArticleStore) UpdateArticle(article *Article) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
	UPDATE articles SET title = $1, description = $2, image = $3, author_id = $4, updated_at = NOW()
	WHERE id = $5;
	`
	result, err := tx.Exec(query, article.Title, article.Description, article.Image, article.AuthorID, article.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("article with ID %d not found", article.ID)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM paragraphs WHERE article_id = $1;`, article.ID)
	if err != nil {
		return err
	}
	for _, paragraph := range article.Paraghraps {
		query :=
		`INSERT INTO paragraphs (article_id, headline, body, order_index)
		VALUES ($1, $2, $3, $4);
		`
		_, err := tx.Exec(query, article.ID, paragraph.Headline, paragraph.Body, paragraph.OrderIndex)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (pg *PostgresArticleStore) DeleteArticle(id int64) error {
	query := `
	DELETE FROM articles WHERE id = $1;
	`
	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("article with ID %d not found", id)
	}
	return nil
}

package store

import (
	"database/sql"
	"fmt"
	"time"
)

type Review struct {
	ID        int64     `json:"id"`
	Stars     int       `json:"stars"`
	Note      *string   `json:"note"`
	AuthorID  int64     `json:"author_id"`
	ArticleID int64     `json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostgresReviewStore struct {
	db *sql.DB
}

func NewPostgresReviewStore(db *sql.DB) *PostgresReviewStore {
	return &PostgresReviewStore{db: db}
}

type ReviewStore interface {
	CreateReview(review *Review) (*Review, error)
	GetReviewByID(id int64) (*Review, error)
	UpdateReview(review *Review) error
	DeleteReview(id int64) error
	GetReviewByUserAndArticle(userID, articleID int64) (*Review, error)
}

func (pg *PostgresReviewStore) CreateReview(review *Review) (*Review, error) {
	query :=
		`INSERT INTO reviews (article_id, author_id, stars, note, created_at, updated_at)
	VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id;
	`
	if review.Stars == 0 || review.Stars > 5 || review.Stars < 0 {
		return nil, fmt.Errorf("invalid stars value: %d", review.Stars)
	}
	err := pg.db.QueryRow(query, review.ArticleID, review.AuthorID, review.Stars, review.Note).Scan(&review.ID)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (pg *PostgresReviewStore) UpdateReview(review *Review) error {
	query := `
	UPDATE reviews SET stars = $1, note = $2, updated_at = NOW()
	WHERE id = $3;
	`
	result, err := pg.db.Exec(query, review.Stars, review.Note, review.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("review with ID %d not found", review.ID)
	}
	return nil
}

func (pg *PostgresReviewStore) GetReviewByID(id int64) (*Review, error) {
	review := &Review{}
	query := `
	SELECT * FROM reviews WHERE id = $1;
	`
	row := pg.db.QueryRow(query, id)
	err := row.Scan(&review.ID, &review.AuthorID, &review.ArticleID, &review.Note, &review.Stars, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return review, nil
}

func (pg *PostgresReviewStore) DeleteReview(id int64) error {
	query := `
	DELETE FROM reviews WHERE id = $1;
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
		return fmt.Errorf("review with ID %d not found", id)
	}
	return nil
}

func (pg *PostgresReviewStore) GetReviewByUserAndArticle(userID, articleID int64) (*Review, error) {
    review := &Review{}
    query := `
    SELECT id, stars, note, author_id, article_id, created_at, updated_at 
    FROM reviews WHERE author_id = $1 AND article_id = $2;
    `
    row := pg.db.QueryRow(query, userID, articleID)
    err := row.Scan(&review.ID, &review.Stars, &review.Note, &review.AuthorID, &review.ArticleID, &review.CreatedAt, &review.UpdatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return review, nil
}

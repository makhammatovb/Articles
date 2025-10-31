package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=postgres dbname=articles sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	_, err = db.Exec("TRUNCATE TABLE articles, paragraphs CASCADE")
	if err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}
	return db
}

func TestCreateArticle(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresArticleStore(db)

	tests := []struct {
		name    string
		article *Article
		wantErr bool
	}{
		{
			name: "Valid Article",
			article: &Article{
				Title:       "Valid Article",
				Description: "Description of valid article",
				Image:       "https://example.com/image.jpg",
				AuthorID:    1,
				Paraghraps: []Paraghraph{
					{
						Headline:   "Headline 1",
						Body:       "Body 1",
						OrderIndex: 1,
					},
					{
						Headline:   "Headline 2",
						Body:       "Body 2",
						OrderIndex: 2,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Article with empty paragraphs",
			article: &Article{
				Title:       "Article with empty paragraphs",
				Description: "Description",
				Image:       "https://example.com/image.jpg",
				AuthorID:    1,
				Paraghraps:  []Paraghraph{},
			},
			wantErr: false,
		},
		{
			name: "Article with nil paragraphs",
			article: &Article{
				Title:       "Article with nil paragraphs",
				Description: "Description",
				Image:       "https://example.com/image.jpg",
				AuthorID:    1,
				Paraghraps:  nil,
			},
			wantErr: false,
		},
		{
			name: "Article with no paragraphs",
			article: &Article{
				Title:       "Article with no paragraphs",
				Description: "Description",
				Image:       "https://example.com/image.jpg",
				AuthorID:    1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdArticle, err := store.CreateArticle(tt.article)
			if tt.wantErr {
				assert.Error(t, err, "Expected error but got none")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.article.Title, createdArticle.Title)
			assert.Equal(t, tt.article.Description, createdArticle.Description)
			assert.Equal(t, tt.article.Image, createdArticle.Image)
			assert.Equal(t, tt.article.AuthorID, createdArticle.AuthorID)
			
			assert.NotZero(t, createdArticle.ID)
			assert.Equal(t, len(tt.article.Paraghraps), len(createdArticle.Paraghraps))
			
			for i, paragraph := range createdArticle.Paraghraps {
				assert.Equal(t, tt.article.Paraghraps[i].Headline, paragraph.Headline)
				assert.Equal(t, tt.article.Paraghraps[i].Body, paragraph.Body)
				assert.Equal(t, tt.article.Paraghraps[i].OrderIndex, paragraph.OrderIndex)
				assert.NotZero(t, paragraph.ID)
			}
			
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func StringPtr(s string) *string {
	return &s
}
package mocks

import (
	"time"

	"github.com/kayden-vs/snippetbox/internal/models"
)

var mockSnippet = &models.Snippet{
	ID:         1,
	Title:      "An old silent Pond",
	Content:    "An old silent Pond...",
	Created:    time.Now(),
	Expires:    time.Now(),
	AuthorName: "Test User",
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int, authorName string) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}

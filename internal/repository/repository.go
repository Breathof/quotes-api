package repository

import (
	"context"
	"time"
)

// Author represents an author in the system
type Author struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Bio       *string   `json:"bio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Quote represents a quote in the system
type Quote struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	AuthorID  int64     `json:"author_id"`
	Author    *Author   `json:"author,omitempty"`
	Source    *string   `json:"source,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// QuoteWithAuthor represents a quote with its author information
type QuoteWithAuthor struct {
	Quote
	AuthorName string  `json:"author_name"`
	AuthorBio  *string `json:"author_bio,omitempty"`
}

// CreateAuthorParams represents parameters for creating an author
type CreateAuthorParams struct {
	Name string  `json:"name" validate:"required,min=1,max=255"`
	Bio  *string `json:"bio,omitempty"`
}

// UpdateAuthorParams represents parameters for updating an author
type UpdateAuthorParams struct {
	Name string  `json:"name" validate:"required,min=1,max=255"`
	Bio  *string `json:"bio,omitempty"`
}

// CreateQuoteParams represents parameters for creating a quote
type CreateQuoteParams struct {
	Content  string   `json:"content" validate:"required,min=1"`
	AuthorID int64    `json:"author_id" validate:"required,min=1"`
	Source   *string  `json:"source,omitempty" validate:"omitempty,max=500"`
	Tags     []string `json:"tags,omitempty"`
}

// UpdateQuoteParams represents parameters for updating a quote
type UpdateQuoteParams struct {
	Content  string   `json:"content" validate:"required,min=1"`
	AuthorID int64    `json:"author_id" validate:"required,min=1"`
	Source   *string  `json:"source,omitempty" validate:"omitempty,max=500"`
	Tags     []string `json:"tags,omitempty"`
}

// ListParams represents pagination parameters
type ListParams struct {
	Limit  int32 `json:"limit" validate:"min=1,max=100"`
	Offset int32 `json:"offset" validate:"min=0"`
}

// AuthorRepository defines the interface for author data access
type AuthorRepository interface {
	Create(ctx context.Context, params CreateAuthorParams) (*Author, error)
	GetByID(ctx context.Context, id int64) (*Author, error)
	List(ctx context.Context, params ListParams) ([]*Author, error)
	Update(ctx context.Context, id int64, params UpdateAuthorParams) (*Author, error)
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int64, error)
	Search(ctx context.Context, query string, params ListParams) ([]*Author, error)
}

// QuoteRepository defines the interface for quote data access
type QuoteRepository interface {
	Create(ctx context.Context, params CreateQuoteParams) (*Quote, error)
	GetByID(ctx context.Context, id int64) (*QuoteWithAuthor, error)
	List(ctx context.Context, params ListParams) ([]*QuoteWithAuthor, error)
	ListByAuthor(ctx context.Context, authorID int64, params ListParams) ([]*QuoteWithAuthor, error)
	Update(ctx context.Context, id int64, params UpdateQuoteParams) (*Quote, error)
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int64, error)
	Search(ctx context.Context, query string, params ListParams) ([]*QuoteWithAuthor, error)
	GetRandom(ctx context.Context) (*QuoteWithAuthor, error)
}

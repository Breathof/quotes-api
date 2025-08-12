package postgres

import (
	"context"
	"fmt"

	"github.com/igferreira/quotes-api/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository wraps the sqlc queries and implements the repository interfaces
type Repository struct {
	db      *pgxpool.Pool
	queries *Queries
}

// NewRepository creates a new PostgreSQL repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:      db,
		queries: New(db),
	}
}

// AuthorRepo returns the author repository
func (r *Repository) AuthorRepo() repository.AuthorRepository {
	return &authorRepository{
		db:      r.db,
		queries: r.queries,
	}
}

// QuoteRepo returns the quote repository
func (r *Repository) QuoteRepo() repository.QuoteRepository {
	return &quoteRepository{
		db:      r.db,
		queries: r.queries,
	}
}

// WithTx executes a function within a database transaction
func (r *Repository) WithTx(ctx context.Context, fn func(repository.AuthorRepository, repository.QuoteRepository) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)
	txRepo := &Repository{
		db:      r.db,
		queries: qtx,
	}

	if err := fn(txRepo.AuthorRepo(), txRepo.QuoteRepo()); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// authorRepository implements repository.AuthorRepository
type authorRepository struct {
	db      *pgxpool.Pool
	queries *Queries
}

// Create creates a new author
func (r *authorRepository) Create(ctx context.Context, params repository.CreateAuthorParams) (*repository.Author, error) {
	author, err := r.queries.CreateAuthor(ctx, CreateAuthorParams{
		Name: params.Name,
		Bio:  params.Bio,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create author: %w", err)
	}

	return &repository.Author{
		ID:        author.ID,
		Name:      author.Name,
		Bio:       author.Bio,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	}, nil
}

// GetByID retrieves an author by ID
func (r *authorRepository) GetByID(ctx context.Context, id int64) (*repository.Author, error) {
	author, err := r.queries.GetAuthor(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("author not found")
		}
		return nil, fmt.Errorf("failed to get author: %w", err)
	}

	return &repository.Author{
		ID:        author.ID,
		Name:      author.Name,
		Bio:       author.Bio,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	}, nil
}

// List retrieves a paginated list of authors
func (r *authorRepository) List(ctx context.Context, params repository.ListParams) ([]*repository.Author, error) {
	authors, err := r.queries.ListAuthors(ctx, ListAuthorsParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list authors: %w", err)
	}

	result := make([]*repository.Author, len(authors))
	for i, author := range authors {
		result[i] = &repository.Author{
			ID:        author.ID,
			Name:      author.Name,
			Bio:       author.Bio,
			CreatedAt: author.CreatedAt,
			UpdatedAt: author.UpdatedAt,
		}
	}

	return result, nil
}

// Update updates an existing author
func (r *authorRepository) Update(ctx context.Context, id int64, params repository.UpdateAuthorParams) (*repository.Author, error) {
	author, err := r.queries.UpdateAuthor(ctx, UpdateAuthorParams{
		ID:   id,
		Name: params.Name,
		Bio:  params.Bio,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("author not found")
		}
		return nil, fmt.Errorf("failed to update author: %w", err)
	}

	return &repository.Author{
		ID:        author.ID,
		Name:      author.Name,
		Bio:       author.Bio,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	}, nil
}

// Delete deletes an author
func (r *authorRepository) Delete(ctx context.Context, id int64) error {
	err := r.queries.DeleteAuthor(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete author: %w", err)
	}
	return nil
}

// Count returns the total number of authors
func (r *authorRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountAuthors(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count authors: %w", err)
	}
	return count, nil
}

// Search searches for authors by name
func (r *authorRepository) Search(ctx context.Context, query string, params repository.ListParams) ([]*repository.Author, error) {
	authors, err := r.queries.SearchAuthorsByName(ctx, SearchAuthorsByNameParams{
		Column1: query,
		Limit:   params.Limit,
		Offset:  params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search authors: %w", err)
	}

	result := make([]*repository.Author, len(authors))
	for i, author := range authors {
		result[i] = &repository.Author{
			ID:        author.ID,
			Name:      author.Name,
			Bio:       author.Bio,
			CreatedAt: author.CreatedAt,
			UpdatedAt: author.UpdatedAt,
		}
	}

	return result, nil
}

// quoteRepository implements repository.QuoteRepository
type quoteRepository struct {
	db      *pgxpool.Pool
	queries *Queries
}

// Create creates a new quote
func (r *quoteRepository) Create(ctx context.Context, params repository.CreateQuoteParams) (*repository.Quote, error) {
	quote, err := r.queries.CreateQuote(ctx, CreateQuoteParams{
		Content:  params.Content,
		AuthorID: params.AuthorID,
		Source:   params.Source,
		Tags:     params.Tags,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	return &repository.Quote{
		ID:        quote.ID,
		Content:   quote.Content,
		AuthorID:  quote.AuthorID,
		Source:    quote.Source,
		Tags:      quote.Tags,
		CreatedAt: quote.CreatedAt,
		UpdatedAt: quote.UpdatedAt,
	}, nil
}

// GetByID retrieves a quote by ID with author information
func (r *quoteRepository) GetByID(ctx context.Context, id int64) (*repository.QuoteWithAuthor, error) {
	row, err := r.queries.GetQuote(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("quote not found")
		}
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	return &repository.QuoteWithAuthor{
		Quote: repository.Quote{
			ID:        row.ID,
			Content:   row.Content,
			AuthorID:  row.AuthorID,
			Source:    row.Source,
			Tags:      row.Tags,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		},
		AuthorName: row.AuthorName,
		AuthorBio:  row.AuthorBio,
	}, nil
}

// List retrieves a paginated list of quotes with author information
func (r *quoteRepository) List(ctx context.Context, params repository.ListParams) ([]*repository.QuoteWithAuthor, error) {
	rows, err := r.queries.ListQuotes(ctx, ListQuotesParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list quotes: %w", err)
	}

	result := make([]*repository.QuoteWithAuthor, len(rows))
	for i, row := range rows {
		result[i] = &repository.QuoteWithAuthor{
			Quote: repository.Quote{
				ID:        row.ID,
				Content:   row.Content,
				AuthorID:  row.AuthorID,
				Source:    row.Source,
				Tags:      row.Tags,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			},
			AuthorName: row.AuthorName,
			AuthorBio:  row.AuthorBio,
		}
	}

	return result, nil
}

// ListByAuthor retrieves quotes by a specific author
func (r *quoteRepository) ListByAuthor(ctx context.Context, authorID int64, params repository.ListParams) ([]*repository.QuoteWithAuthor, error) {
	rows, err := r.queries.ListQuotesByAuthor(ctx, ListQuotesByAuthorParams{
		AuthorID: authorID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list quotes by author: %w", err)
	}

	result := make([]*repository.QuoteWithAuthor, len(rows))
	for i, row := range rows {
		result[i] = &repository.QuoteWithAuthor{
			Quote: repository.Quote{
				ID:        row.ID,
				Content:   row.Content,
				AuthorID:  row.AuthorID,
				Source:    row.Source,
				Tags:      row.Tags,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			},
			AuthorName: row.AuthorName,
			AuthorBio:  row.AuthorBio,
		}
	}

	return result, nil
}

// Update updates an existing quote
func (r *quoteRepository) Update(ctx context.Context, id int64, params repository.UpdateQuoteParams) (*repository.Quote, error) {
	quote, err := r.queries.UpdateQuote(ctx, UpdateQuoteParams{
		ID:       id,
		Content:  params.Content,
		AuthorID: params.AuthorID,
		Source:   params.Source,
		Tags:     params.Tags,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("quote not found")
		}
		return nil, fmt.Errorf("failed to update quote: %w", err)
	}

	return &repository.Quote{
		ID:        quote.ID,
		Content:   quote.Content,
		AuthorID:  quote.AuthorID,
		Source:    quote.Source,
		Tags:      quote.Tags,
		CreatedAt: quote.CreatedAt,
		UpdatedAt: quote.UpdatedAt,
	}, nil
}

// Delete deletes a quote
func (r *quoteRepository) Delete(ctx context.Context, id int64) error {
	err := r.queries.DeleteQuote(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete quote: %w", err)
	}
	return nil
}

// Count returns the total number of quotes
func (r *quoteRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.queries.CountQuotes(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count quotes: %w", err)
	}
	return count, nil
}

// Search searches for quotes by content
func (r *quoteRepository) Search(ctx context.Context, query string, params repository.ListParams) ([]*repository.QuoteWithAuthor, error) {
	rows, err := r.queries.SearchQuotesByContent(ctx, SearchQuotesByContentParams{
		Column1: query,
		Limit:   params.Limit,
		Offset:  params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search quotes: %w", err)
	}

	result := make([]*repository.QuoteWithAuthor, len(rows))
	for i, row := range rows {
		result[i] = &repository.QuoteWithAuthor{
			Quote: repository.Quote{
				ID:        row.ID,
				Content:   row.Content,
				AuthorID:  row.AuthorID,
				Source:    row.Source,
				Tags:      row.Tags,
				CreatedAt: row.CreatedAt,
				UpdatedAt: row.UpdatedAt,
			},
			AuthorName: row.AuthorName,
			AuthorBio:  row.AuthorBio,
		}
	}

	return result, nil
}

// GetRandom retrieves a random quote
func (r *quoteRepository) GetRandom(ctx context.Context) (*repository.QuoteWithAuthor, error) {
	row, err := r.queries.GetRandomQuote(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no quotes found")
		}
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	return &repository.QuoteWithAuthor{
		Quote: repository.Quote{
			ID:        row.ID,
			Content:   row.Content,
			AuthorID:  row.AuthorID,
			Source:    row.Source,
			Tags:      row.Tags,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		},
		AuthorName: row.AuthorName,
		AuthorBio:  row.AuthorBio,
	}, nil
}

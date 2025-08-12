package service

import (
	"context"
	"fmt"

	"github.com/igferreira/quotes-api/internal/repository"
)

// Service provides business logic for the quotes API
type Service struct {
	authorRepo repository.AuthorRepository
	quoteRepo  repository.QuoteRepository
}

// NewService creates a new service instance
func NewService(authorRepo repository.AuthorRepository, quoteRepo repository.QuoteRepository) *Service {
	return &Service{
		authorRepo: authorRepo,
		quoteRepo:  quoteRepo,
	}
}

// CreateAuthor creates a new author
func (s *Service) CreateAuthor(ctx context.Context, params repository.CreateAuthorParams) (*repository.Author, error) {
	// Check if author already exists
	authors, err := s.authorRepo.Search(ctx, params.Name, repository.ListParams{Limit: 1, Offset: 0})
	if err != nil {
		return nil, fmt.Errorf("failed to check existing authors: %w", err)
	}
	if len(authors) > 0 && authors[0].Name == params.Name {
		return nil, fmt.Errorf("author with name %q already exists", params.Name)
	}

	author, err := s.authorRepo.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create author: %w", err)
	}

	return author, nil
}

// GetAuthor retrieves an author by ID
func (s *Service) GetAuthor(ctx context.Context, id int64) (*repository.Author, error) {
	author, err := s.authorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get author: %w", err)
	}
	return author, nil
}

// ListAuthors retrieves a paginated list of authors
func (s *Service) ListAuthors(ctx context.Context, params repository.ListParams) ([]*repository.Author, int64, error) {
	// Get total count
	total, err := s.authorRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count authors: %w", err)
	}

	// Get authors
	authors, err := s.authorRepo.List(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list authors: %w", err)
	}

	return authors, total, nil
}

// UpdateAuthor updates an existing author
func (s *Service) UpdateAuthor(ctx context.Context, id int64, params repository.UpdateAuthorParams) (*repository.Author, error) {
	// Check if author exists
	_, err := s.authorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}

	author, err := s.authorRepo.Update(ctx, id, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update author: %w", err)
	}

	return author, nil
}

// DeleteAuthor deletes an author
func (s *Service) DeleteAuthor(ctx context.Context, id int64) error {
	// Check if author has quotes
	quotes, err := s.quoteRepo.ListByAuthor(ctx, id, repository.ListParams{Limit: 1, Offset: 0})
	if err != nil {
		return fmt.Errorf("failed to check author quotes: %w", err)
	}
	if len(quotes) > 0 {
		return fmt.Errorf("cannot delete author with existing quotes")
	}

	err = s.authorRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete author: %w", err)
	}

	return nil
}

// SearchAuthors searches for authors by name
func (s *Service) SearchAuthors(ctx context.Context, query string, params repository.ListParams) ([]*repository.Author, int64, error) {
	authors, err := s.authorRepo.Search(ctx, query, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search authors: %w", err)
	}

	// For search, we don't have an efficient way to get total count
	// In a real application, you might want to implement a separate count search method
	return authors, int64(len(authors)), nil
}

// CreateQuote creates a new quote
func (s *Service) CreateQuote(ctx context.Context, params repository.CreateQuoteParams) (*repository.Quote, error) {
	// Verify author exists
	_, err := s.authorRepo.GetByID(ctx, params.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}

	quote, err := s.quoteRepo.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	return quote, nil
}

// GetQuote retrieves a quote by ID
func (s *Service) GetQuote(ctx context.Context, id int64) (*repository.QuoteWithAuthor, error) {
	quote, err := s.quoteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}
	return quote, nil
}

// ListQuotes retrieves a paginated list of quotes
func (s *Service) ListQuotes(ctx context.Context, params repository.ListParams) ([]*repository.QuoteWithAuthor, int64, error) {
	// Get total count
	total, err := s.quoteRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count quotes: %w", err)
	}

	// Get quotes
	quotes, err := s.quoteRepo.List(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list quotes: %w", err)
	}

	return quotes, total, nil
}

// ListQuotesByAuthor retrieves quotes by a specific author
func (s *Service) ListQuotesByAuthor(ctx context.Context, authorID int64, params repository.ListParams) ([]*repository.QuoteWithAuthor, error) {
	// Verify author exists
	_, err := s.authorRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}

	quotes, err := s.quoteRepo.ListByAuthor(ctx, authorID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list quotes by author: %w", err)
	}

	return quotes, nil
}

// UpdateQuote updates an existing quote
func (s *Service) UpdateQuote(ctx context.Context, id int64, params repository.UpdateQuoteParams) (*repository.Quote, error) {
	// Check if quote exists
	_, err := s.quoteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("quote not found: %w", err)
	}

	// Verify new author exists
	_, err = s.authorRepo.GetByID(ctx, params.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}

	quote, err := s.quoteRepo.Update(ctx, id, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update quote: %w", err)
	}

	return quote, nil
}

// DeleteQuote deletes a quote
func (s *Service) DeleteQuote(ctx context.Context, id int64) error {
	err := s.quoteRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete quote: %w", err)
	}
	return nil
}

// SearchQuotes searches for quotes by content
func (s *Service) SearchQuotes(ctx context.Context, query string, params repository.ListParams) ([]*repository.QuoteWithAuthor, int64, error) {
	quotes, err := s.quoteRepo.Search(ctx, query, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search quotes: %w", err)
	}

	return quotes, int64(len(quotes)), nil
}

// GetRandomQuote retrieves a random quote
func (s *Service) GetRandomQuote(ctx context.Context) (*repository.QuoteWithAuthor, error) {
	quote, err := s.quoteRepo.GetRandom(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}
	return quote, nil
}

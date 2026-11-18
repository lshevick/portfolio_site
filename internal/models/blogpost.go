package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lshevick/internal/store"
)

type BlogPostModel struct {
	Store *store.Store
}

func NewBlogPostModel(s *store.Store) (*BlogPostModel, error) {
	return &BlogPostModel{Store: s}, nil
}

// BlogPost is the data we expose. CreatedAt/UpdatedAt are parsed from the
// sqlite datetime strings returned by the driver.
type BlogPost struct {
	ID         int
	Title      string
	Slug       string
	AuthorID   int
	AuthorName string // denormalized for display; avoids a JOIN in templates
	Body       string
	Published  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Slugify turns arbitrary text into a URL-friendly slug: lowercase, trim, and
// collapse every run of whitespace/punctuation into a single hyphen.
func Slugify(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(s))), "-")
}

func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	// modernc.org/sqlite returns ISO 8601 / RFC3339 by default; fall back to the
	// classic sqlite "YYYY-MM-DD HH:MM:SS" layout for other drivers.
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("parsing time %q", s)
}

// GetAll returns every post newest-first. Callers filter Published as needed.
func (m *BlogPostModel) GetAll() ([]*BlogPost, error) {
	rows, err := m.Store.DB.Query(
		"SELECT id, title, slug, author_id, author_name, body, published, created_at, updated_at FROM blog_posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*BlogPost
	for rows.Next() {
		p := &BlogPost{}
		if err := scanBlogPost(rows, p); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

// Get returns a post by slug. Returns ErrNoRecord when not found.
func (m *BlogPostModel) Get(slug string) (*BlogPost, error) {
	rows, err := m.Store.DB.Query(
		"SELECT id, title, slug, author_id, author_name, body, published, created_at, updated_at FROM blog_posts WHERE slug = ?", slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var p BlogPost
	found := false
	for rows.Next() {
		if err := scanBlogPost(rows, &p); err != nil {
			return nil, err
		}
		found = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !found { // no row -> not found
		return nil, ErrNoRecord
	}
	return &p, nil
}

// Create inserts a new post. Slug is UNIQUE, so we pre-check to return
// ErrDuplicateSlug cleanly instead of parsing the driver's constraint message.
func (m *BlogPostModel) Create(p *BlogPost) error {
	if _, err := m.Get(p.Slug); err == nil {
		return ErrDuplicateSlug
	}

	res, err := m.Store.DB.Exec(
		"INSERT INTO blog_posts (title, slug, author_id, author_name, body, published, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		p.Title, p.Slug, p.AuthorID, p.AuthorName, p.Body, p.Published, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

// Update modifies an existing post by id.
func (m *BlogPostModel) Update(id int, p *BlogPost) error {
	_, err := m.Store.DB.Exec(
		"UPDATE blog_posts SET title = ?, slug = ?, author_id = ?, author_name = ?, body = ?, published = ?, updated_at = ? WHERE id = ?",
		p.Title, p.Slug, p.AuthorID, p.AuthorName, p.Body, p.Published, time.Now().UTC(), id)
	return err
}

// Delete removes a post by id.
func (m *BlogPostModel) Delete(id int) error {
	_, err := m.Store.DB.Exec("DELETE FROM blog_posts WHERE id = ?", id)
	return err
}

func scanBlogPost(row *sql.Rows, p *BlogPost) error {
	var published int
	createdAtRaw, updatedAtRaw := "", "" // distinct locals: Scan args can't reuse a name twice
	if err := row.Scan(
		&p.ID, &p.Title, &p.Slug, &p.AuthorID, &p.AuthorName, &p.Body, &published, &createdAtRaw, &updatedAtRaw); err != nil {
		return err
	}
	p.Published = published != 0

	createdAt, err := parseTime(createdAtRaw)
	if err != nil {
		return err
	}
	p.CreatedAt = createdAt

	updatedAt, err := parseTime(updatedAtRaw)
	if err != nil {
		return err
	}
	p.UpdatedAt = updatedAt
	return nil
}

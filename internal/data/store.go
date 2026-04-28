package data

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var variablePattern = regexp.MustCompile(`\{\{([A-Za-z_][A-Za-z0-9_]*)\}\}`)

type Snippet struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Language    string    `json:"language"`
	Type        string    `json:"type"`
	UsageCount  int       `json:"usage_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []string  `json:"tags,omitempty"`
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) CreateSnippet(snippet Snippet) (Snippet, error) {
	if strings.TrimSpace(snippet.Name) == "" {
		return Snippet{}, errors.New("snippet name is required")
	}
	if snippet.Type == "" {
		snippet.Type = "inline"
	}

	result, err := s.db.Exec(`
		INSERT INTO snippets (name, description, content, language, type)
		VALUES (?, ?, ?, ?, ?)
	`, snippet.Name, snippet.Description, snippet.Content, snippet.Language, snippet.Type)
	if err != nil {
		return Snippet{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Snippet{}, err
	}
	snippet.ID = id

	if err := s.replaceVariables(id, ExtractVariables(snippet.Content)); err != nil {
		return Snippet{}, err
	}
	for _, tag := range snippet.Tags {
		if err := s.TagSnippet(snippet.Name, tag); err != nil {
			return Snippet{}, err
		}
	}

	return s.GetSnippet(snippet.Name)
}

func (s *Store) UpdateSnippet(name string, snippet Snippet) (Snippet, error) {
	existing, err := s.GetSnippet(name)
	if err != nil {
		return Snippet{}, err
	}
	if snippet.Name == "" {
		snippet.Name = existing.Name
	}
	if snippet.Type == "" {
		snippet.Type = existing.Type
	}

	_, err = s.db.Exec(`
		INSERT INTO snippet_versions (snippet_id, content)
		VALUES (?, ?)
	`, existing.ID, existing.Content)
	if err != nil {
		return Snippet{}, err
	}

	_, err = s.db.Exec(`
		UPDATE snippets
		SET name = ?, description = ?, content = ?, language = ?, type = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, snippet.Name, snippet.Description, snippet.Content, snippet.Language, snippet.Type, existing.ID)
	if err != nil {
		return Snippet{}, err
	}

	if err := s.replaceVariables(existing.ID, ExtractVariables(snippet.Content)); err != nil {
		return Snippet{}, err
	}
	if snippet.Tags != nil {
		if err := s.ReplaceSnippetTags(existing.ID, snippet.Tags); err != nil {
			return Snippet{}, err
		}
	}

	return s.GetSnippet(snippet.Name)
}

func (s *Store) UpsertSnippet(snippet Snippet) (Snippet, error) {
	existing, err := s.GetSnippet(snippet.Name)
	if err == nil {
		if snippet.Tags == nil {
			snippet.Tags = existing.Tags
		}
		return s.UpdateSnippet(existing.Name, snippet)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return s.CreateSnippet(snippet)
	}
	return Snippet{}, err
}

func (s *Store) GetSnippet(name string) (Snippet, error) {
	row := s.db.QueryRow(`
		SELECT id, name, COALESCE(description, ''), content, COALESCE(language, ''), COALESCE(type, 'inline'),
		       COALESCE(usage_count, 0), created_at, updated_at
		FROM snippets
		WHERE name = ?
	`, name)
	return s.scanSnippet(row)
}

func (s *Store) UseSnippet(name string) (Snippet, error) {
	_, err := s.db.Exec("UPDATE snippets SET usage_count = usage_count + 1 WHERE name = ?", name)
	if err != nil {
		return Snippet{}, err
	}
	return s.GetSnippet(name)
}

func (s *Store) ListSnippets() ([]Snippet, error) {
	rows, err := s.db.Query(`
		SELECT id, name, COALESCE(description, ''), content, COALESCE(language, ''), COALESCE(type, 'inline'),
		       COALESCE(usage_count, 0), created_at, updated_at
		FROM snippets
		ORDER BY updated_at DESC, name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanSnippets(rows)
}

func (s *Store) SearchSnippets(query string) ([]Snippet, error) {
	like := "%" + query + "%"
	rows, err := s.db.Query(`
		SELECT id, name, COALESCE(description, ''), content, COALESCE(language, ''), COALESCE(type, 'inline'),
		       COALESCE(usage_count, 0), created_at, updated_at
		FROM snippets
		WHERE name LIKE ? OR description LIKE ? OR content LIKE ? OR language LIKE ?
		ORDER BY updated_at DESC, name ASC
	`, like, like, like, like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanSnippets(rows)
}

func (s *Store) DeleteSnippet(name string) error {
	result, err := s.db.Exec("DELETE FROM snippets WHERE name = ?", name)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Store) TagSnippet(name string, tag string) error {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return errors.New("tag is required")
	}
	snippet, err := s.GetSnippet(name)
	if err != nil {
		return err
	}

	result, err := s.db.Exec(`
		INSERT INTO tags (name, usage_count)
		VALUES (?, 1)
		ON CONFLICT(name) DO UPDATE SET usage_count = usage_count + 1, updated_at = CURRENT_TIMESTAMP
	`, tag)
	if err != nil {
		return err
	}

	tagID, err := result.LastInsertId()
	if err != nil || tagID == 0 {
		err = s.db.QueryRow("SELECT id FROM tags WHERE name = ?", tag).Scan(&tagID)
	}
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT OR IGNORE INTO snippet_tags (snippet_id, tag_id)
		VALUES (?, ?)
	`, snippet.ID, tagID)
	return err
}

func (s *Store) ReplaceSnippetTags(snippetID int64, tags []string) error {
	if _, err := s.db.Exec("DELETE FROM snippet_tags WHERE snippet_id = ?", snippetID); err != nil {
		return err
	}
	var name string
	if err := s.db.QueryRow("SELECT name FROM snippets WHERE id = ?", snippetID).Scan(&name); err != nil {
		return err
	}
	seen := map[string]bool{}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		if err := s.TagSnippet(name, tag); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) SnippetTags(snippetID int64) ([]string, error) {
	rows, err := s.db.Query(`
		SELECT tags.name
		FROM tags
		JOIN snippet_tags ON snippet_tags.tag_id = tags.id
		WHERE snippet_tags.snippet_id = ?
		ORDER BY tags.name ASC
	`, snippetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func ExtractVariables(content string) []string {
	matches := variablePattern.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool, len(matches))
	var variables []string
	for _, match := range matches {
		key := match[1]
		if !seen[key] {
			seen[key] = true
			variables = append(variables, key)
		}
	}
	sort.Strings(variables)
	return variables
}

func RenderTemplate(content string, values map[string]string, defaults map[string]string) (string, error) {
	var missing []string
	rendered := variablePattern.ReplaceAllStringFunc(content, func(token string) string {
		key := strings.TrimSuffix(strings.TrimPrefix(token, "{{"), "}}")
		if value, ok := values[key]; ok {
			return value
		}
		if value, ok := defaults[key]; ok {
			return value
		}
		missing = append(missing, key)
		return token
	})
	if len(missing) > 0 {
		return "", fmt.Errorf("missing values for variables: %s", strings.Join(missing, ", "))
	}
	return rendered, nil
}

func (s *Store) VariableDefaults(snippetID int64) (map[string]string, error) {
	rows, err := s.db.Query("SELECT key, default_value FROM variables WHERE snippet_id = ?", snippetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	defaults := map[string]string{}
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		defaults[key] = value
	}
	return defaults, rows.Err()
}

func (s *Store) replaceVariables(snippetID int64, variables []string) error {
	if _, err := s.db.Exec("DELETE FROM variables WHERE snippet_id = ?", snippetID); err != nil {
		return err
	}
	for _, variable := range variables {
		_, err := s.db.Exec(`
			INSERT INTO variables (snippet_id, key, default_value)
			VALUES (?, ?, '')
		`, snippetID, variable)
		if err != nil {
			return err
		}
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func (s *Store) scanSnippet(row rowScanner) (Snippet, error) {
	var snippet Snippet
	err := row.Scan(
		&snippet.ID,
		&snippet.Name,
		&snippet.Description,
		&snippet.Content,
		&snippet.Language,
		&snippet.Type,
		&snippet.UsageCount,
		&snippet.CreatedAt,
		&snippet.UpdatedAt,
	)
	if err != nil {
		return Snippet{}, err
	}
	snippet.Tags, err = s.SnippetTags(snippet.ID)
	if err != nil {
		return Snippet{}, err
	}
	return snippet, nil
}

func (s *Store) scanSnippets(rows *sql.Rows) ([]Snippet, error) {
	snippets := []Snippet{}
	for rows.Next() {
		var snippet Snippet
		if err := rows.Scan(
			&snippet.ID,
			&snippet.Name,
			&snippet.Description,
			&snippet.Content,
			&snippet.Language,
			&snippet.Type,
			&snippet.UsageCount,
			&snippet.CreatedAt,
			&snippet.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tags, err := s.SnippetTags(snippet.ID)
		if err != nil {
			return nil, err
		}
		snippet.Tags = tags
		snippets = append(snippets, snippet)
	}
	return snippets, rows.Err()
}

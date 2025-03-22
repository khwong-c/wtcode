package codecard

import "github.com/google/uuid"

// CodeCard is the main model of this service, representing snippets of code, descriptions, and examples.
type CodeCard struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Code        string    `json:"code"`
	Language    string    `json:"language"`
	Description string    `json:"description"`
	Example     string    `json:"example"`

	Author      *string `json:"author"`
	AuthorEmail *string `json:"author_email"`
	AuthorURL   *string `json:"author_url"`

	Verified bool `json:"verified"`
}

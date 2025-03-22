package card

import "github.com/google/uuid"

type Card struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Code        string    `json:"code"`
	Language    string    `json:"language"`
	Description string    `json:"description"`
	Example     string    `json:"example"`

	Author      *string `json:"author"`
	AuthorEmail *string `json:"author_email"`
	AuthorURL   *string `json:"author_url"`
}

type CreateCardInput struct {
	Title       *string `json:"title,omitempty"`
	Code        *string `json:"code,omitempty"`
	Language    *string `json:"language,omitempty"`
	Description *string `json:"description,omitempty"`
	Example     *string `json:"example,omitempty"`

	Author      *string `json:"author,omitempty"`
	AuthorEmail *string `json:"author_email,omitempty"`
	AuthorURL   *string `json:"author_url,omitempty"`
}

func (c *CreateCardInput) Validate() error {
	// TODO: Implement validation
	return nil
}

type UpdateCardInput struct {
	Title       *string `json:"title,omitempty"`
	Code        *string `json:"code,omitempty"`
	Language    *string `json:"language,omitempty"`
	Description *string `json:"description,omitempty"`
	Example     *string `json:"example,omitempty"`

	Author      *string `json:"author,omitempty"`
	AuthorEmail *string `json:"author_email,omitempty"`
	AuthorURL   *string `json:"author_url,omitempty"`
}

func (c *UpdateCardInput) Validate() error {
	// TODO: Implement validation
	return nil
}

type SupportedLanguage struct {
	Language string `json:"language"`
}

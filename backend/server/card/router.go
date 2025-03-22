package card

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/unrolled/render"

	"github.com/khwong-c/wtcode/authentication"
	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/server/middlewares"
	"github.com/khwong-c/wtcode/tooling/di"
)

type adminAPI interface {
	GetCards(ids []uuid.UUID) ([]*Card, error)
	ListID(batchSize, page int) ([]uuid.UUID, error)
	CreateCard(card *CreateCardInput) (*Card, error)
	VerifyCard(id uuid.UUID) (*Card, error)
	DeleteCard(id uuid.UUID) error
	UpdateCard(id uuid.UUID, card *UpdateCardInput) (*Card, error)
}

type publicAPI interface {
	GetRandomCardIDs() ([]uuid.UUID, error)
	GetCard(id uuid.UUID) (Card, error)
	GetSupportedLanguages() []SupportedLanguage
}

type API struct {
	config *config.Config
	auth   authentication.Authenticator
	admin  adminAPI
	public publicAPI
	render *render.Render
}

func Create(injector *do.Injector) (*API, error) {
	return &API{
		config: di.InvokeOrProvide(injector, config.LoadConfig),
		auth:   di.InvokeOrProvide(injector, authentication.CreateAPIKeyAuthenticator),
		render: render.New(),
	}, nil
}

func (c *API) Mount(r chi.Router) {
	r.Use(middleware.Recoverer)

	// Admin API
	r.Group(func(r chi.Router) {
		r.Use(middlewares.RequireAdminAccess(c.config, c.auth))
		r.Post("/admin/cards/details", c.getCardsHandler)
		r.Get("/admin/cards", c.listCardIDsHandler)
		r.Post("/admin/card", c.createCardHandler)
		r.Put("/admin/card/{id}", c.updateCardHandler)
		r.Put("/admin/card/{id}/verify", c.verifyCardHandler)
		r.Delete("/admin/card/{id}", c.deleteCardHandler)
	})
	// Public API
	r.Group(func(r chi.Router) {
		r.Get("/cards", c.getRandomCardIDsHandler)
		r.Get("/card/{id}", c.getCardHandler)
		r.Get("/languages", c.getSupportedLanguages)
	})
}

package card

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juju/errors"

	"github.com/khwong-c/wtcode/server/common"
)

func (c *API) getCardsHandler(w http.ResponseWriter, r *http.Request) {
	type cardList []uuid.UUID
	var ids cardList
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		panic(errors.Annotate(err, "Failed to decode card id list"))
	}
	cards, err := c.admin.GetCards(ids)
	if err != nil {
		panic(errors.Annotate(err, "Failed to get cards"))
	}
	if err := c.render.JSON(w, http.StatusOK, cards); err != nil {
		panic(errors.Trace(err))
	}
}

func (c *API) listCardIDsHandler(w http.ResponseWriter, r *http.Request) {
	const (
		defaultBatchSize = "50"
		defaultPage      = "0"
	)
	var bs, page int
	var err error
	bsURL, pageURL := r.URL.Query().Get("batch_size"), r.URL.Query().Get("page")
	if bsURL == "" {
		bsURL = defaultBatchSize
	}
	if pageURL == "" {
		pageURL = defaultPage
	}
	if bs, err = strconv.Atoi(bsURL); err != nil {
		panic(errors.Annotate(err, "Failed to parse batch_size"))
	}
	if page, err = strconv.Atoi(pageURL); err != nil {
		panic(errors.Annotate(err, "Failed to parse page"))
	}
	ids, err := c.admin.ListID(bs, page)
	if err != nil {
		panic(errors.Annotate(err, "Failed to list card ids"))
	}
	if err := c.render.JSON(w, http.StatusOK, ids); err != nil {
		panic(errors.Trace(err))
	}
}

func (c *API) createCardHandler(w http.ResponseWriter, r *http.Request) {
	newCard := CreateCardInput{}
	if err := json.NewDecoder(r.Body).Decode(&newCard); err != nil {
		panic(
			errors.Annotate(
				errors.WithType(err, errors.NotValid),
				"Failed to decode new card",
			),
		)
	}
	if err := newCard.Validate(); err != nil {
		panic(errors.WithType(err, errors.NotValid))
	}

	card, err := c.admin.CreateCard(&newCard)
	if err != nil {
		panic(errors.Annotate(err, "Failed to create card"))
	}
	if err := c.render.JSON(w, http.StatusOK, card); err != nil {
		panic(errors.Trace(err))
	}
}

func (c *API) verifyCardHandler(w http.ResponseWriter, r *http.Request) {
	idURL := chi.URLParam(r, "id")
	id, err := uuid.Parse(idURL)
	if err != nil {
		panic(errors.Annotate(
			errors.WithType(err, errors.NotValid),
			"Failed to parse id",
		))
	}

	card, err := c.admin.VerifyCard(id)
	if err != nil {
		panic(errors.Annotate(
			errors.WithType(err, errors.NotValid),
			"Failed to verify card",
		))
	}
	if err := c.render.JSON(w, http.StatusOK, card); err != nil {
		panic(errors.Trace(err))
	}
}

func (c *API) updateCardHandler(w http.ResponseWriter, r *http.Request) {
	idURL := chi.URLParam(r, "id")
	id, err := uuid.Parse(idURL)
	if err != nil {
		panic(errors.Annotate(
			errors.WithType(err, errors.NotValid),
			"Failed to parse id",
		))
	}

	newCard := UpdateCardInput{}
	if err := json.NewDecoder(r.Body).Decode(&newCard); err != nil {
		panic(errors.Annotate(
			errors.WithType(err, errors.NotValid),
			"Failed to decode new card",
		))
	}
	if err := newCard.Validate(); err != nil {
		panic(errors.WithType(err, errors.NotValid))
	}

	card, err := c.admin.UpdateCard(id, &newCard)
	if err != nil {
		panic(errors.Annotate(err, "Failed to create card"))
	}
	if err := c.render.JSON(w, http.StatusOK, card); err != nil {
		panic(errors.Trace(err))
	}
}

func (c *API) deleteCardHandler(w http.ResponseWriter, r *http.Request) {
	idURL := chi.URLParam(r, "id")
	id, err := uuid.Parse(idURL)
	if err != nil {
		panic(errors.Annotate(
			errors.WithType(err, errors.NotValid),
			"Failed to parse id",
		))
	}
	if err := c.admin.DeleteCard(id); err != nil {
		panic(errors.Annotate(err, "Failed to delete card"))
	}
	if err := c.render.JSON(w, http.StatusOK, common.SuccessResponse); err != nil {
		panic(errors.Trace(err))
	}
}

func (c *API) getRandomCardIDsHandler(w http.ResponseWriter, _ *http.Request) {
	cardIDs, err := c.public.GetRandomCardIDs()
	if err != nil {
		panic(errors.Annotate(err, "Failed to get random card ids"))
	}
	if err := c.render.JSON(w, http.StatusOK, cardIDs); err != nil {
		panic(errors.Trace(err))
	}
}

func (c *API) getCardHandler(w http.ResponseWriter, r *http.Request) {
	idURL := chi.URLParam(r, "id")
	id, err := uuid.Parse(idURL)
	if err != nil {
		panic(errors.Annotate(
			errors.WithType(err, errors.NotValid),
			"Failed to parse id",
		))
	}
	card, err := c.public.GetCard(id)
	if err != nil {
		panic(errors.Annotate(err, "Failed to get card"))
	}
	if err := c.render.JSON(w, http.StatusOK, card); err != nil {
		panic(errors.Trace(err))
	}
}

func (c *API) getSupportedLanguages(w http.ResponseWriter, _ *http.Request) {
	langs := c.public.GetSupportedLanguages()
	if err := c.render.JSON(w, http.StatusOK, langs); err != nil {
		panic(errors.Trace(err))
	}
}

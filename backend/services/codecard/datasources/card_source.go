package datasources

import (
	"context"

	"github.com/google/uuid"
	"github.com/juju/errors"
	"github.com/samber/do"
	"gorm.io/gorm"

	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/drivers"
	"github.com/khwong-c/wtcode/tooling/di"
)

type CodeCardSource interface {
	ListCardIDs(ctx context.Context, page int, size int) (uuid.UUIDs, error)
	GetCards(ctx context.Context, ids uuid.UUIDs) ([]*CodeCard, error)
	AddCard(ctx context.Context, card *CodeCard) error
	DeleteCard(ctx context.Context, id uuid.UUID) error
	ModifyCard(ctx context.Context, id uuid.UUID, card *CodeCard) error
}

type codeCardSource struct {
	injector *do.Injector
	cfg      *config.Config
	db       *gorm.DB
}

func CreateCodeCardSource(injector *do.Injector) (CodeCardSource, error) {
	src := &codeCardSource{
		injector: injector,
		cfg:      di.InvokeOrProvide(injector, config.LoadConfig),
		db: di.InvokeOrProvideNamed(
			injector, string(drivers.SQLTargetDefault), drivers.DialSQL,
		).DB(),
	}
	if src.cfg.DBSetup.AutoMigrate {
		if err := src.migrateModels(); err != nil {
			return nil, err
		}
	}

	return src, nil
}

func (s *codeCardSource) ListCardIDs(ctx context.Context, page int, size int) (uuid.UUIDs, error) {
	var ids uuid.UUIDs
	if size < 0 || page < 0 {
		return ids, errors.NotValidf("page and size must be non-negative")
	}
	q := s.db.WithContext(ctx).Model(&CodeCard{})
	if size > 0 {
		off := page * size
		q = q.Offset(off).Limit(size)
	}
	err := q.Pluck("id", &ids).Error
	return ids, err
}

func (s *codeCardSource) GetCards(ctx context.Context, ids uuid.UUIDs) ([]*CodeCard, error) {
	var cards []*CodeCard
	if len(ids) == 0 {
		return cards, nil
	}

	if len(ids) == 1 {
		var card CodeCard
		err := s.db.WithContext(ctx).First(&card, ids[0]).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.WithType(err, errors.NotFound)
		}
		return []*CodeCard{&card}, err
	}

	err := s.db.WithContext(ctx).Where("id IN ?", ids).Find(&cards).Error
	return cards, err
}

func (s *codeCardSource) AddCard(ctx context.Context, card *CodeCard) error {
	var err error
	if card.ID, err = uuid.NewV7(); err != nil {
		return errors.Trace(err)
	}
	return s.db.WithContext(ctx).Create(card).Error
}

func (s *codeCardSource) DeleteCard(ctx context.Context, id uuid.UUID) error {
	err := s.db.WithContext(ctx).Delete(&CodeCard{}, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithType(err, errors.NotFound)
	}
	return err
}

func (s *codeCardSource) ModifyCard(ctx context.Context, id uuid.UUID, card *CodeCard) error {
	if card.ID != id {
		newCard := *card
		newCard.ID = id
		card = &newCard
	}

	err := s.db.WithContext(ctx).Save(card).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.WithType(err, errors.NotFound)
	}
	return err
}

package datasources

import (
	"context"
	"testing"

	"github.com/samber/do"

	"github.com/khwong-c/wtcode/config"
	"github.com/khwong-c/wtcode/tests/fixtures"
	"github.com/khwong-c/wtcode/tooling/di"
)

func TestPluck(t *testing.T) {
	injector := di.CreateInjector(false, false)
	di.InvokeOrProvide(injector, func(injector *do.Injector) (*config.Config, error) {
		cfg, err := fixtures.CreateDefaultConfig(false)
		if err != nil {
			return nil, err
		}
		// Patch the Config for testing.
		cfg.SQLTarget.Default = "sqlite::memory:"
		cfg.DBSetup.AutoMigrate = true
		return cfg, nil
	})
	src := di.InvokeOrProvide(injector, CreateCodeCardSource)
	ctx := context.Background()
	src.AddCard(ctx, &CodeCard{Title: "Card 1"})
	src.AddCard(ctx, &CodeCard{Title: "Card 2"})
	l, _ := src.ListCardIDs(ctx, 0, 0)
	t.Log(l)
	for _, id := range l {
		src.ModifyCard(ctx, id, &CodeCard{Title: "Card " + id.String()})
	}
	cs, _ := src.GetCards(ctx, l)
	for _, c := range cs {
		t.Log(c)
	}

}

package datasources

import (
	"github.com/juju/errors"
)

var allModels = []interface{}{
	&CodeCard{},
}

func (s *codeCardSource) migrateModels() error {
	if err := s.db.AutoMigrate(allModels...); err != nil {
		return errors.Annotate(err, "failed to auto migrate models")
	}
	return nil
}

package builder

import (
	"fmt"
	"github.com/kilianstallz/stage-sync/internal/constants"
)

type ConnectionStringStrategy interface {
	VerifyString(connString string) error
}

type ConnectionString struct {
	strategies map[constants.DBType]ConnectionStringStrategy
}

func NewDBMetaBuilder() ConnectionString {
	return ConnectionString{strategies: make(map[constants.DBType]ConnectionStringStrategy)}
}

func (s *ConnectionString) RegisterStrategy(name constants.DBType, strategy ConnectionStringStrategy) {
	s.strategies[name] = strategy
}

func (s *ConnectionString) VerifyString(strat constants.DBType, connString string) error {
	strategy := s.strategies[strat]
	if err := strategy.VerifyString(connString); err != nil {
		return fmt.Errorf("ConnectionStringStrategy.VerifyString: %w", err)
	}

	return nil
}

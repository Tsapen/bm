package migrator

import (
	"database/sql"
	"fmt"
	"os"
)

type (
	Executor interface {
		Exec(query string, args ...interface{}) (sql.Result, error)
	}

	Migrator struct {
		executor       Executor
		migrationsPath string
	}
)

func New(migrationsPath string, executor Executor) *Migrator {
	return &Migrator{
		migrationsPath: migrationsPath,
		executor:       executor,
	}
}

func (m *Migrator) Apply() error {
	content, err := os.ReadFile(m.migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get migrations by path '%s'", m.migrationsPath)
	}

	_, err = m.executor.Exec(string(content))
	if err != nil {
		return fmt.Errorf("execute migrations: %w", err)
	}

	return nil
}

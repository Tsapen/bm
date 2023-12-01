// Package migrator provides functionality to apply database migrations.
package migrator

import (
	"database/sql"
	"fmt"
	"os"
)

// Executor is an interface that defines the methods required for executing SQL queries.
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Migrator represents the migration manager.
type Migrator struct {
	executor       Executor
	migrationsPath string
}

// New creates a new instance of Migrator.
func New(migrationsPath string, executor Executor) *Migrator {
	return &Migrator{
		migrationsPath: migrationsPath,
		executor:       executor,
	}
}

// Apply reads and executes SQL migrations from the specified file path.
func (m *Migrator) Apply() error {
	// Read migration content from the specified file path
	content, err := os.ReadFile(m.migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get migrations by path '%s'", m.migrationsPath)
	}

	// Execute the SQL migrations
	_, err = m.executor.Exec(string(content))
	if err != nil {
		return fmt.Errorf("execute migrations: %w", err)
	}

	return nil
}

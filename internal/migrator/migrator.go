// Package migrator provides functionality to apply database migrations.
package migrator

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// // Executor is an interface that defines the methods required for executing SQL queries.
// type Executor interface {
// 	Exec(query string, args ...interface{}) (sql.Result, error)
// }

// // Migrator represents the migration manager.
// type Migrator struct {
// 	// executor       Executor
// 	m              *migrate.Migrate
// 	migrationsPath string
// }

// New creates a new instance of Migrator.
// func New(migrationsPath string, executor Executor) (*Migrator, error) {
func ApplyMigrations(migrationsPath string, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{
		SchemaName: "public",
	})
	if err != nil {
		return fmt.Errorf("init instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("init migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// // Apply reads and executes SQL migrations from the specified file path.
// func (m *Migrator) ApplyDeprecated() error {
// 	// Read migration content from the specified file path
// 	content, err := os.ReadFile(m.migrationsPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to get migrations by path '%s'", m.migrationsPath)
// 	}

// 	// Execute the SQL migrations
// 	_, err = m.executor.Exec(string(content))
// 	if err != nil {
// 		return fmt.Errorf("execute migrations: %w", err)
// 	}

// 	return nil
// }

// func (m *Migrator) Apply() error {
// 	// Apply all available migrations
// 	if err := m.m.Up(); err != nil && err != migrate.ErrNoChange {
// 		return err
// 	}

// 	return nil
// }

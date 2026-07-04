package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Migration struct {
	Version string
	Path    string
	SQL     string
}

func RunMigrations(ctx context.Context, db *sql.DB, dir string) error {
	if err := ensureSchemaMigrations(ctx, db); err != nil {
		return err
	}

	migrations, err := readUpMigrations(dir)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		applied, err := isMigrationApplied(ctx, db, migration.Version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		if err := applyMigration(ctx, db, migration); err != nil {
			return err
		}
	}

	return nil
}

func ensureSchemaMigrations(ctx context.Context, db *sql.DB) error {
	query := `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version VARCHAR(190) NOT NULL,
				applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				PRIMARY KEY (version)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`
	_, err := db.ExecContext(ctx, query)
	return err
}

func readUpMigrations(dir string) ([]Migration, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)

	migrations := make([]Migration, 0, len(files))
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		base := filepath.Base(file)
		version := strings.TrimSuffix(base, ".up.sql")
		migrations = append(migrations, Migration{
			Version: version,
			Path:    file,
			SQL:     string(content),
		})
	}

	return migrations, nil
}

func isMigrationApplied(ctx context.Context, db *sql.DB, version string) (bool, error) {
	var exists int
	err := db.QueryRowContext(ctx, "SELECT COUNT(1) FROM schema_migrations WHERE version = ?", version).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func applyMigration(ctx context.Context, db *sql.DB, migration Migration) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statements := splitSQLStatements(migration.SQL)
	for _, statement := range statements {
		if strings.TrimSpace(statement) == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("apply %s: %w", migration.Path, err)
		}
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)", migration.Version, time.Now()); err != nil {
		return err
	}

	return tx.Commit()
}

func splitSQLStatements(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	statements := make([]string, 0, len(parts))
	for _, part := range parts {
		statement := strings.TrimSpace(part)
		if statement != "" {
			statements = append(statements, statement)
		}
	}
	return statements
}

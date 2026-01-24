package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

// NOTE: go:embed does not allow ".." in patterns, so we keep a copy of SQL files
// under this package directory.
//
//go:embed sql/*.sql
var embeddedMigrations embed.FS

// Up applies all pending migrations. It is safe to call on every startup.
func Up(db *pgxpool.Pool) error {
	slog.Info("migrations: Up() called")

	// Some versions of golang-migrate's pgx/v5 driver expect a *sql.DB,
	// so we bridge from pgxpool.Pool via pgx/v5/stdlib.
	sqlDB := stdlib.OpenDBFromPool(db)
	defer func() { _ = sqlDB.Close() }()
	slog.Info("migrations: opened stdlib *sql.DB from pool")

	// Ensure the stdlib DB is actually reachable early (helps surface config issues).
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		slog.Error("migrations: ping failed", slog.String("error", err.Error()))
		return fmt.Errorf("ping db before migrations: %w", err)
	}
	slog.Info("migrations: ping succeeded")

	driver, err := migratepgx.WithInstance(sqlDB, &migratepgx.Config{})
	if err != nil {
		slog.Error("migrations: create driver failed", slog.String("error", err.Error()))
		return fmt.Errorf("create postgres migrate driver: %w", err)
	}
	slog.Info("migrations: driver created")

	migrationsFS, err := fs.Sub(embeddedMigrations, "sql")
	if err != nil {
		slog.Error("migrations: fs.Sub failed", slog.String("error", err.Error()))
		return fmt.Errorf("open embedded migrations fs: %w", err)
	}
	slog.Info("migrations: embedded FS ready")

	src, err := iofs.New(migrationsFS, ".")
	if err != nil {
		slog.Error("migrations: iofs.New failed", slog.String("error", err.Error()))
		return fmt.Errorf("create iofs source: %w", err)
	}
	slog.Info("migrations: iofs source created")

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		slog.Error("migrations: NewWithInstance failed", slog.String("error", err.Error()))
		return fmt.Errorf("create migrate instance: %w", err)
	}
	slog.Info("migrations: migrate instance created")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("migrations: Up() failed", slog.String("error", err.Error()))
		return fmt.Errorf("migrate up: %w", err)
	}
	slog.Info("migrations: Up() finished (no change or success)")

	return nil
}

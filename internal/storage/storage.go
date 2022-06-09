package storage

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/fancar/tmp_xm/internal/config"
)

// Migrations
//go:embed migrations/*
var migrations embed.FS

var (
	jwtsecret []byte
	// HashIterations denfines the number of times a password is hashed.
	HashIterations = 100000
)

// Setup configures the storage package.
func Setup(c config.Config) error {
	jwtsecret = []byte(c.ExternalAPI.JWTSecret)
	// HashIterations = c.General.PasswordHashIterations

	// setup timezone
	// if err := SetTimeLocation(c.Metrics.Timezone); err != nil {
	// 	return errors.Wrap(err, "set time location error")
	// }

	log.Info("storage: connecting to PostgreSQL database ...")
	log.Debugf("storage: PostgreSQL DSN: %s \n", c.PostgreSQL.DSN)
	d, err := sqlx.Open("postgres", c.PostgreSQL.DSN)
	if err != nil {
		return errors.Wrap(err, "storage: PostgreSQL connection error")
	}
	d.SetMaxOpenConns(c.PostgreSQL.MaxOpenConnections)
	d.SetMaxIdleConns(c.PostgreSQL.MaxIdleConnections)

	for {
		if err := d.Ping(); err != nil {
			log.WithError(err).Warning("storage: ping PostgreSQL database error, will retry in 2s")
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	db = &DBLogger{d}

	if c.PostgreSQL.Automigrate {
		if err := MigrateUp(d); err != nil {
			return err
		}
	}

	return nil
}

// MigrateUp configure postgres migration up
func MigrateUp(db *sqlx.DB) error {
	log.Info("storage: applying PostgreSQL data migrations from migrations dir ...")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return fmt.Errorf("storage: migrate postgres driver error: %w", err)
	}

	src, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return fmt.Errorf("new httpfs error: %w", err)
	}

	m, err := migrate.NewWithInstance("httpfs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("storage: new migrate instance error: %w", err)
	}

	oldVersion, _, _ := m.Version()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("storage: migrate up error: %w", err)
	}

	newVersion, _, _ := m.Version()

	if oldVersion != newVersion {
		log.WithFields(log.Fields{
			"from_version": oldVersion,
			"to_version":   newVersion,
		}).Info("storage: PostgreSQL data migrations applied")
	}

	return nil
}

// MigrateDown configure postgres migration down
func MigrateDown(db *sqlx.DB) error {
	log.Info("storage: reverting PostgreSQL data migrations from migrations dir ...")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return fmt.Errorf("storage: migrate postgres driver error: %w", err)
	}

	src, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return fmt.Errorf("new httpfs error: %w", err)
	}

	m, err := migrate.NewWithInstance("httpfs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("storage: new migrate instance error: %w", err)
	}

	oldVersion, _, _ := m.Version()

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("storage: migrate down error: %w", err)
	}

	newVersion, _, _ := m.Version()

	if oldVersion != newVersion {
		log.WithFields(log.Fields{
			"from_version": oldVersion,
			"to_version":   newVersion,
		}).Info("storage: reverted PostgreSQL data migrations applied")
	}

	return nil
}

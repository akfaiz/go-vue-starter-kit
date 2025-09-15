package migrate

import (
	"context"
	"database/sql"

	_ "github.com/akfaiz/go-vue-starter-kit/db/migrations"
	"github.com/akfaiz/go-vue-starter-kit/internal/config"
	"github.com/akfaiz/migris"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/urfave/cli/v3"
)

const migrationDir = "db/migrations"

var Command = &cli.Command{
	Name:  "migrate",
	Usage: "Database migration commands",
	Commands: []*cli.Command{
		{
			Name:  "create",
			Usage: "Create a new migration file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Aliases:  []string{"n"},
					Required: true,
					Usage:    "Name of the migration",
				},
			},
			Action: func(ctx context.Context, c *cli.Command) error {
				return migris.Create(migrationDir, c.String("name"))
			},
		},
		{
			Name:  "up",
			Usage: "Apply all up migrations",
			Action: func(ctx context.Context, c *cli.Command) error {
				migrate, err := newMigrator()
				if err != nil {
					return err
				}
				return migrate.UpContext(ctx)
			},
		},
		{
			Name:  "reset",
			Usage: "Rollback all migrations",
			Action: func(ctx context.Context, c *cli.Command) error {
				migrate, err := newMigrator()
				if err != nil {
					return err
				}
				return migrate.ResetContext(ctx)
			},
		},
		{
			Name:  "down",
			Usage: "Rollback the last migration",
			Action: func(ctx context.Context, c *cli.Command) error {
				migrate, err := newMigrator()
				if err != nil {
					return err
				}
				return migrate.DownContext(ctx)
			},
		},
		{
			Name:  "status",
			Usage: "Show the status of all migrations",
			Action: func(ctx context.Context, c *cli.Command) error {
				migrate, err := newMigrator()
				if err != nil {
					return err
				}
				return migrate.StatusContext(ctx)
			},
		},
	},
}

func newMigrator() (*migris.Migrate, error) {
	cfg := config.Load()
	db, err := sql.Open("pgx", cfg.Database.DSN())
	if err != nil {
		return nil, err
	}

	migrate, err := migris.New("pgx", migris.WithMigrationDir(migrationDir), migris.WithDB(db))
	if err != nil {
		return nil, err
	}
	return migrate, nil
}

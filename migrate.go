package rdb

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/theNullP0inter/googly/ingress"

	"github.com/spf13/cobra"
)

// MigrationIngress can be registered with googly and be used to run migrations
type MigrationIngress struct {
	*ingress.BaseIngress
}

// NewMigrationIngress creates a new MigrationIngress
func NewMigrationIngress(name string, path string, driver database.Driver) *MigrationIngress {

	cmd := NewMigrateCommand(
		&ingress.CommandConfig{
			Name:  name,
			Short: "DB Migrator",
		},
		fmt.Sprintf("file://%s", path),
		"mysql", driver,
	)
	ing := ingress.NewBaseIngress(cmd)
	return &MigrationIngress{ing}

}

// getMigration is a helper function to create a new instance of migrate.Migrate
func getMigration(path string, driverName string, driver database.Driver) *migrate.Migrate {

	m, err := migrate.NewWithDatabaseInstance(
		path,
		driverName,
		driver,
	)
	if err != nil {
		panic(err)
	}
	return m
}

// NewMigrateCommand creates a new MigrateCommand which is a *cobra.Command
//
// It also creates 2 sub commands: up & down
// up can be used to apply all your migrations
// down can be used to rollback all your migrations
func NewMigrateCommand(config *ingress.CommandConfig, path string, driverName string, driver database.Driver) *cobra.Command {

	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migration on a database driver",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(("Running Migrations"))
		},
	}

	var migrateUpCmd = &cobra.Command{
		Use:   "up",
		Short: "Ups all the migrations",
		Run: func(cmd *cobra.Command, args []string) {
			m := getMigration(path, driverName, driver)
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				panic(err)
			}
		},
	}

	var migrateDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Downs all the migrations",
		Run: func(cmd *cobra.Command, args []string) {
			m := getMigration(path, driverName, driver)
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				panic(err)
			}
		},
	}

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)

	return migrateCmd

}

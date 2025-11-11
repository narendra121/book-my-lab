/*
Copyright © 2025
*/
package main

import (
	"fmt"
	"log"
	"os"

	"booking.com/internal/config"
	"booking.com/internal/db/postgresql/dao"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var db *gorm.DB

// --------------------------------------------------
// main entry point
// --------------------------------------------------
func main() {
	Execute()
}

// --------------------------------------------------
// rootCmd: Base command
// --------------------------------------------------
var rootCmd = &cobra.Command{
	Use:   "bookmylab",
	Short: "BookMyLab CLI tool",
	Long: `BookMyLab CLI helps developers run migrations and generate database models.
It’s built using Cobra and integrates with GORM and golang-migrate.`,
}

// --------------------------------------------------
// migrateCmd: Run database migrations (up/down)
// --------------------------------------------------
var migrateCmd = &cobra.Command{
	Use:   "migrate [direction] [path]",
	Args:  cobra.RangeArgs(1, 2),
	Short: "Run database migrations (up or down)",
	Long: `Executes database migrations using golang-migrate.
Usage examples:
  bookmylab migrate up
  bookmylab migrate down
  bookmylab migrate up internal/db/migrations`,
	Run: func(cmd *cobra.Command, args []string) {
		direction := args[0]
		path := "internal/db/migrations"
		if len(args) == 2 {
			path = args[1]
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to get sql.DB instance: %v", err)
		}

		driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
		if err != nil {
			log.Fatalf("failed to create driver: %v", err)
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://"+path,
			"postgres",
			driver,
		)
		if err != nil {
			log.Fatalf("failed to create migrate instance: %v", err)
		}
		fmt.Println(m.Version())
		fmt.Println("🚀 Running migrations...")
		switch direction {
		case "down":
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Printf("❌ Migration failed: %v", err)
				os.Exit(1)
			}
			fmt.Println("✅ Rolled back successfully!")
		default:
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Printf("❌ Migration failed: %v", err)
				os.Exit(1)
			}
			fmt.Println("✅ Migrations applied successfully!")
		}
		os.Exit(0)
	},
}

// --------------------------------------------------
// genDbModelsCmd: Generate DB models using GORM Gen
// --------------------------------------------------
var genDbModelsCmd = &cobra.Command{
	Use:   "gen-db-models",
	Short: "Generate GORM models from PostgreSQL schema",
	Long: `Automatically generates model structs from your PostgreSQL schema
and stores them inside internal/db/postgresql/dao.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🧱 Generating DB models...")

		g := gen.NewGenerator(gen.Config{
			OutPath:           "internal/db/postgresql/dao",
			Mode:              gen.WithoutContext | gen.WithDefaultQuery,
			FieldWithIndexTag: true,
			FieldWithTypeTag:  true,
			// Uncomment these for more customization:
			// FieldNullable:     true,
			// FieldCoverable:    true,
		})

		g.UseDB(db)
		allTables := g.GenerateAllTable()
		g.ApplyBasic(allTables...)
		g.Execute()

		fmt.Println("✅ Model generation completed successfully!")
	},
}

// --------------------------------------------------
// Execute: Entry point for CLI execution
// --------------------------------------------------
func Execute() {
	rootCmd.AddCommand(migrateCmd, genDbModelsCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// --------------------------------------------------
// init: Setup configuration and connect to DB
// --------------------------------------------------
func init() {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Println("❌ Error loading app configuration:", err)
		return
	}

	db, err = dao.Connect(cfg.PostgresqlDb)
	if err != nil {
		log.Println("❌ Error connecting to database:", err)
		return
	}
}

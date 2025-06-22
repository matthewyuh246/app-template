package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"

	"app-template/pkg/database"
)

func main() {
	// 環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// データベース接続
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// マイグレーション実行
	if err := runMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}

func runMigrations(db *sql.DB) error {
	// MySQL ドライバーの設定
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create mysql driver: %w", err)
	}

	// マイグレーションファイルのパスを設定
	migrationsPath := "file://database/migrations"
	
	// migrate インスタンスを作成
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "mysql", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// コマンドライン引数をチェック
	args := os.Args[1:]
	if len(args) == 0 {
		// 引数がない場合は up を実行
		return runUp(m)
	}

	switch args[0] {
	case "up":
		return runUp(m)
	case "down":
		return runDown(m)
	case "drop":
		return runDrop(m)
	case "version":
		return showVersion(m)
	case "force":
		if len(args) < 2 {
			return fmt.Errorf("force command requires version argument")
		}
		return runForce(m, args[1])
	default:
		return fmt.Errorf("unknown command: %s. Available commands: up, down, drop, version, force", args[0])
	}
}

func runUp(m *migrate.Migrate) error {
	log.Println("Running migrations up...")
	err := m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations up: %w", err)
	}
	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}
	return nil
}

func runDown(m *migrate.Migrate) error {
	log.Println("Running migrations down...")
	err := m.Steps(-1)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations down: %w", err)
	}
	if err == migrate.ErrNoChange {
		log.Println("No migrations to rollback")
	} else {
		log.Println("Migration rolled back successfully")
	}
	return nil
}

func runDrop(m *migrate.Migrate) error {
	log.Println("Dropping all migrations...")
	err := m.Drop()
	if err != nil {
		return fmt.Errorf("failed to drop migrations: %w", err)
	}
	log.Println("All migrations dropped successfully")
	return nil
}

func showVersion(m *migrate.Migrate) error {
	version, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("failed to get migration version: %w", err)
	}
	log.Printf("Current migration version: %d, dirty: %v", version, dirty)
	return nil
}

func runForce(m *migrate.Migrate, versionStr string) error {
	var version int
	if _, err := fmt.Sscanf(versionStr, "%d", &version); err != nil {
		return fmt.Errorf("invalid version format: %s", versionStr)
	}
	
	log.Printf("Forcing migration to version %d...", version)
	err := m.Force(version)
	if err != nil {
		return fmt.Errorf("failed to force migration: %w", err)
	}
	log.Printf("Migration forced to version %d successfully", version)
	return nil
} 
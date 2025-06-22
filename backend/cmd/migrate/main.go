package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/go-sql-driver/mysql"
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
	// マイグレーションテーブルを作成
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// マイグレーションファイルを取得
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// 実行済みマイグレーションを取得
	executedMigrations, err := getExecutedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %w", err)
	}

	// 新しいマイグレーションを実行
	for _, file := range migrationFiles {
		if !contains(executedMigrations, file) {
			if err := executeMigration(db, file); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", file, err)
			}
			log.Printf("Executed migration: %s", file)
		}
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err := db.Exec(query)
	return err
}

func getMigrationFiles() ([]string, error) {
	files, err := filepath.Glob("database/migrations/*.sql")
	if err != nil {
		return nil, err
	}

	// ファイル名でソート
	sort.Strings(files
	
	// ファイル名のみを取得
	var migrationFiles []string
	for _, file := range files {
		migrationFiles = append(migrationFiles, filepath.Base(file))
	}

	return migrationFiles, nil
}

func getExecutedMigrations(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		migrations = append(migrations, version)
	}

	return migrations, nil
}

func executeMigration(db *sql.DB, filename string) error {
	// マイグレーションファイルを読み込み
	content, err := os.ReadFile(filepath.Join("database/migrations", filename))
	if err != nil {
		return err
	}

	// SQLを実行
	if _, err := db.Exec(string(content)); err != nil {
		return err
	}

	// マイグレーション記録を追加
	_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", filename)
	return err
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
} 
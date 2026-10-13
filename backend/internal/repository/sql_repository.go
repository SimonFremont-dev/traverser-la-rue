package repository

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

// SQLRepository handles SQL (PostgreSQL) database operations with Atlas integration
type SQLRepository struct {
	db     *sql.DB
	dsn    string
	migrationsDir string
}

// NewSQLRepository creates a new SQL repository
func NewSQLRepository(dsn string) (*SQLRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &SQLRepository{
		db: db,
		dsn: dsn,
		migrationsDir: "migrations",
	}, nil
}

// NewSQLRepositoryWithMigrationsDir creates a new SQL repository with custom migrations directory
func NewSQLRepositoryWithMigrationsDir(dsn string, migrationsDir string) (*SQLRepository, error) {
	repo, err := NewSQLRepository(dsn)
	if err != nil {
		return nil, err
	}
	repo.migrationsDir = migrationsDir
	return repo, nil
}

// Close closes the database connection
func (r *SQLRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// DB returns the underlying database connection
func (r *SQLRepository) DB() *sql.DB {
	return r.db
}

// RunMigrations runs all SQL migrations using Atlas-style naming
func (r *SQLRepository) RunMigrations(migrationsDir string) error {
	// Create atlas_schema_revisions table if it doesn't exist
	if _, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS atlas_schema_revisions (
			version TEXT PRIMARY KEY,
			description TEXT,
			executed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("failed to create atlas_schema_revisions table: %w", err)
	}

	// Get already applied migrations
	applied := make(map[string]bool)
	rows, err := r.db.Query("SELECT version FROM atlas_schema_revisions")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var version string
			if err := rows.Scan(&version); err != nil {
				log.Printf("Warning: failed to scan migration version: %v", err)
			}
			applied[version] = true
		}
	}

	// Get migration files
	files, err := getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Run migrations in order
	for _, file := range files {
		// Extract version from filename (e.g., "20261011_000001_create_skill_schema.sql" -> "20261011_000001")
		version := strings.TrimSuffix(file, ".sql")
		version = filepath.Base(version) // Remove path if present
		
		if applied[version] {
			log.Printf("Migration %s already applied, skipping", version)
			continue
		}

		// Read migration file
		content, err := os.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Extract description from +atlas:sum comment if present
		description := extractMigrationDescription(string(content))

		// Execute migration in a transaction
		tx, err := r.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", version, err)
		}

		log.Printf("Applying migration: %s - %s", version, description)
		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", version, err)
		}

		// Record applied migration
		if _, err := tx.Exec(
			"INSERT INTO atlas_schema_revisions (version, description) VALUES ($1, $2)",
			version,
			description,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", version, err)
		}

		log.Printf("Migration %s applied successfully", version)
	}

	return nil
}

// RunMigrationsFromFS runs migrations from an embedded filesystem
func (r *SQLRepository) RunMigrationsFromFS(fsys fs.FS, migrationsDir string) error {
	// Store the filesystem for later use
	r.migrationsDir = ""
	
	// Create atlas_schema_revisions table if it doesn't exist
	if _, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS atlas_schema_revisions (
			version TEXT PRIMARY KEY,
			description TEXT,
			executed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("failed to create atlas_schema_revisions table: %w", err)
	}

	// Get already applied migrations
	applied := make(map[string]bool)
	rows, err := r.db.Query("SELECT version FROM atlas_schema_revisions")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var version string
			if err := rows.Scan(&version); err != nil {
				log.Printf("Warning: failed to scan migration version: %v", err)
			}
			applied[version] = true
		}
	}

	// Get migration files from embedded filesystem
	files, err := getMigrationFilesFromFS(fsys, migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Run migrations in order
	for _, file := range files {
		// Extract version from filename
		version := strings.TrimSuffix(file, ".sql")
		version = filepath.Base(version)
		
		if applied[version] {
			log.Printf("Migration %s already applied, skipping", version)
			continue
		}

		// Read migration file from embedded FS
		content, err := fs.ReadFile(fsys, filepath.Join(migrationsDir, file))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Extract description
		description := extractMigrationDescription(string(content))

		// Execute migration in a transaction
		tx, err := r.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", version, err)
		}

		log.Printf("Applying migration: %s - %s", version, description)
		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", version, err)
		}

		// Record applied migration
		if _, err := tx.Exec(
			"INSERT INTO atlas_schema_revisions (version, description) VALUES ($1, $2)",
			version,
			description,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", version, err)
		}

		log.Printf("Migration %s applied successfully", version)
	}

	return nil
}

// getMigrationFiles returns sorted list of migration files
func getMigrationFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		// Match Atlas-style migration files: YYYYMMDD_HHMMSS_description.sql
		if strings.HasSuffix(name, ".sql") && len(name) > 4 {
			// Check if it starts with a timestamp-like pattern
			if name[0] >= '0' && name[0] <= '9' {
				migrationFiles = append(migrationFiles, name)
			}
		}
	}

	// Sort by filename (which includes timestamp)
	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

// getMigrationFilesFromFS returns sorted list of migration files from embedded filesystem
func getMigrationFilesFromFS(fsys fs.FS, dir string) ([]string, error) {
	files, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if strings.HasSuffix(name, ".sql") && len(name) > 4 {
			if name[0] >= '0' && name[0] <= '9' {
				migrationFiles = append(migrationFiles, name)
			}
		}
	}

	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

// extractMigrationDescription extracts the description from +atlas:sum comment
func extractMigrationDescription(content string) string {
	// Look for +atlas:sum comment
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "-- +atlas:sum") {
			// Extract the value
			parts := strings.SplitN(trimmed, "\"", 3)
			if len(parts) >= 3 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

// GetAppliedMigrations returns list of applied migration versions
func (r *SQLRepository) GetAppliedMigrations() ([]string, error) {
	var migrations []string

	rows, err := r.db.Query("SELECT version FROM atlas_schema_revisions ORDER BY version")
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		migrations = append(migrations, version)
	}

	return migrations, nil
}

// GetPendingMigrations returns list of pending migration files
func (r *SQLRepository) GetPendingMigrations() ([]string, error) {
	// Get applied migrations
	applied, err := r.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}
	appliedSet := make(map[string]bool)
	for _, v := range applied {
		appliedSet[v] = true
	}

	// Get all migration files
	files, err := getMigrationFiles(r.migrationsDir)
	if err != nil {
		return nil, err
	}

	// Filter for pending ones
	var pending []string
	for _, file := range files {
		version := strings.TrimSuffix(file, ".sql")
		version = filepath.Base(version)
		if !appliedSet[version] {
			pending = append(pending, file)
		}
	}

	return pending, nil
}

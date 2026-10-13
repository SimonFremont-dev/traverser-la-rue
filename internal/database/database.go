// Package database provides SQLite database operations for the CV Builder
package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database path constant
const (
	DefaultDBPath = "data/cv_builder.db"
	DBDir         = "data"
)

// Database holds the SQLite database connection
type Database struct {
	Path string
	Conn *sql.DB
}

// NewDatabase creates a new database instance
func NewDatabase(dbPath string) (*Database, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better performance
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, fmt.Errorf("failed to set journal mode: %w", err)
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return &Database{
		Path: dbPath,
		Conn: db,
	}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	if db.Conn != nil {
		return db.Conn.Close()
	}
	return nil
}

// InitSchema initializes the database schema
func (db *Database) InitSchema() error {
	// SQL statements for creating tables
	schema := `
	-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Resumes table
	CREATE TABLE IF NOT EXISTS resumes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		content_format TEXT NOT NULL DEFAULT 'json',
		is_current BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- Skills master table (all possible skills)
	CREATE TABLE IF NOT EXISTS skills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		category TEXT,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- User skills table (skills the user has)
	CREATE TABLE IF NOT EXISTS user_skills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		skill_id INTEGER NOT NULL,
		proficiency TEXT CHECK(proficiency IN ('Beginner', 'Intermediate', 'Advanced', 'Expert')),
		years_experience INTEGER DEFAULT 0,
		last_verified_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE,
		UNIQUE(user_id, skill_id)
	);

	-- Work experience table
	CREATE TABLE IF NOT EXISTS work_experience (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		company TEXT NOT NULL,
		position TEXT NOT NULL,
		start_date DATE NOT NULL,
		end_date DATE,
		current BOOLEAN DEFAULT FALSE,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- Education table
	CREATE TABLE IF NOT EXISTS education (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		institution TEXT NOT NULL,
		degree TEXT NOT NULL,
		field_of_study TEXT,
		start_date DATE,
		end_date DATE,
		gpa REAL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- Job offers table
	CREATE TABLE IF NOT EXISTS job_offers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		company TEXT NOT NULL,
		position TEXT NOT NULL,
		location TEXT,
		description TEXT,
		source TEXT,
		raw_content TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	-- Job requirements table
	CREATE TABLE IF NOT EXISTS job_requirements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_offer_id INTEGER NOT NULL,
		requirement_type TEXT NOT NULL CHECK(requirement_type IN ('skill', 'experience', 'education', 'certification')),
		name TEXT NOT NULL,
		category TEXT,
		is_required BOOLEAN DEFAULT TRUE,
		min_years INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (job_offer_id) REFERENCES job_offers(id) ON DELETE CASCADE
	);

	-- User responses to job requirements
	CREATE TABLE IF NOT EXISTS user_responses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		job_requirement_id INTEGER NOT NULL,
		has_skill BOOLEAN NOT NULL,
		proficiency TEXT CHECK(proficiency IN ('Beginner', 'Intermediate', 'Advanced', 'Expert')),
		years_experience INTEGER,
		notes TEXT,
		responded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (job_requirement_id) REFERENCES job_requirements(id) ON DELETE CASCADE,
		UNIQUE(user_id, job_requirement_id)
	);

	-- Indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_skills_name ON skills(name);
	CREATE INDEX IF NOT EXISTS idx_user_skills_user ON user_skills(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_skills_skill ON user_skills(skill_id);
	CREATE INDEX IF NOT EXISTS idx_user_skills_last_verified ON user_skills(last_verified_at);
	CREATE INDEX IF NOT EXISTS idx_job_offers_user ON job_offers(user_id);
	CREATE INDEX IF NOT EXISTS idx_job_requirements_job ON job_requirements(job_offer_id);
	CREATE INDEX IF NOT EXISTS idx_user_responses_user ON user_responses(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_responses_expires ON user_responses(expires_at);
	`

	// Execute the schema
	_, err := db.Conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// GetDBPath returns the default database path
func GetDBPath() string {
	return DefaultDBPath
}

// InitDefaultDB initializes the default database
func InitDefaultDB() (*Database, error) {
	db, err := NewDatabase(DefaultDBPath)
	if err != nil {
		return nil, err
	}

	if err := db.InitSchema(); err != nil {
		return nil, err
	}

	return db, nil
}

// Time helpers
var (
	// SkillExpirationDuration is the duration after which a skill needs to be re-verified (1 year)
	SkillExpirationDuration = 365 * 24 * time.Hour
)

// IsSkillExpired checks if a skill verification has expired
func IsSkillExpired(lastVerifiedAt *time.Time) bool {
	if lastVerifiedAt == nil {
		return true
	}
	return time.Since(*lastVerifiedAt) > SkillExpirationDuration
}

// GetExpirationDate returns the expiration date for a skill verification
func GetExpirationDate() time.Time {
	return time.Now().Add(SkillExpirationDuration)
}

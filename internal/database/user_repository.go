package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *Database
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db}
}

// UserDB represents a user in the database
type UserDB struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(name, email string) (*UserDB, error) {
	query := `INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)`
	result, err := r.db.Conn.Exec(query, name, email, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &UserDB{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*UserDB, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = ?`
	row := r.db.Conn.QueryRow(query, email)

	user := &UserDB{}
	var createdAtStr, updatedAtStr string
	err := row.Scan(&user.ID, &user.Name, &user.Email, &createdAtStr, &updatedAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	// Parse timestamps
	if createdAtStr != "" {
		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	}
	if updatedAtStr != "" {
		user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int64) (*UserDB, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?`
	row := r.db.Conn.QueryRow(query, id)

	user := &UserDB{}
	var createdAtStr, updatedAtStr string
	err := row.Scan(&user.ID, &user.Name, &user.Email, &createdAtStr, &updatedAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	// Parse timestamps
	if createdAtStr != "" {
		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	}
	if updatedAtStr != "" {
		user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	}

	return user, nil
}

// UpdateUser updates a user's information
func (r *UserRepository) UpdateUser(id int64, name, email string) error {
	query := `UPDATE users SET name = ?, email = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Conn.Exec(query, name, email, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user and all associated data (cascade delete)
func (r *UserRepository) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers lists all users
func (r *UserRepository) ListUsers() ([]*UserDB, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := r.db.Conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*UserDB
	for rows.Next() {
		user := &UserDB{}
		var createdAtStr, updatedAtStr string
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &createdAtStr, &updatedAtStr); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		// Parse timestamps
		if createdAtStr != "" {
			user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}
		if updatedAtStr != "" {
			user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		}

		users = append(users, user)
	}

	return users, nil
}

// ImportResume imports a resume to initialize or update a user's profile
func (r *UserRepository) ImportResume(userID int64, resume *models.Resume) error {
	// Create skill repository for this operation
	skillRepo := NewSkillRepository(r.db)

	// Get or create user skills
	for _, skill := range resume.Skills {
		if err := skillRepo.UpsertSkillAndUserSkill(userID, skill); err != nil {
			return fmt.Errorf("failed to upsert skill: %w", err)
		}
	}

	// Get or create work experience
	for _, exp := range resume.Experience {
		if err := skillRepo.UpsertWorkExperience(userID, exp); err != nil {
			return fmt.Errorf("failed to upsert work experience: %w", err)
		}
	}

	// Get or create education
	for _, edu := range resume.Education {
		if err := skillRepo.UpsertEducation(userID, edu); err != nil {
			return fmt.Errorf("failed to upsert education: %w", err)
		}
	}

	// Store the resume
	if err := r.StoreResume(userID, resume); err != nil {
		return fmt.Errorf("failed to store resume: %w", err)
	}

	return nil
}

// StoreResume stores a resume in the database
func (r *UserRepository) StoreResume(userID int64, resume *models.Resume) error {
	// First, mark all existing resumes as not current
	_, err := r.db.Conn.Exec(`UPDATE resumes SET is_current = FALSE WHERE user_id = ?`, userID)
	if err != nil {
		return fmt.Errorf("failed to update existing resumes: %w", err)
	}

	// Insert the new resume
	query := `INSERT INTO resumes (user_id, title, content, content_format, is_current, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = r.db.Conn.Exec(query, userID, resume.Metadata.Title, "", "json", true, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert resume: %w", err)
	}

	return nil
}

package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
)

// SkillRepository handles skill-related database operations
type SkillRepository struct {
	db *Database
}

// NewSkillRepository creates a new skill repository
func NewSkillRepository(db *Database) *SkillRepository {
	return &SkillRepository{db: db}
}

// SkillDB represents a skill in the database
type SkillDB struct {
	ID          int64
	Name        string
	Category    string
	Description string
	CreatedAt   time.Time
}

// UserSkillDB represents a user's skill with proficiency
type UserSkillDB struct {
	ID            int64
	UserID       int64
	SkillID      int64
	Proficiency  string
	Years        int
	LastVerified *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	SkillName    string
	Category     string
}

// GetOrCreateSkill gets a skill by name or creates it if it doesn't exist
func (r *SkillRepository) GetOrCreateSkill(name, category, description string) (*SkillDB, error) {
	// Try to get the skill first
	query := `SELECT id, name, category, description, created_at FROM skills WHERE name = ?`
	row := r.db.Conn.QueryRow(query, name)

	skill := &SkillDB{}
	var createdAtStr string
	err := row.Scan(&skill.ID, &skill.Name, &skill.Category, &skill.Description, &createdAtStr)
	if err == nil {
		if createdAtStr != "" {
			skill.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}
		return skill, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to scan skill: %w", err)
	}

	// Create the skill
	insertQuery := `INSERT INTO skills (name, category, description, created_at) VALUES (?, ?, ?, ?)`
	result, err := r.db.Conn.Exec(insertQuery, name, category, description, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create skill: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &SkillDB{
		ID:          id,
		Name:        name,
		Category:    category,
		Description: description,
		CreatedAt:   time.Now(),
	}, nil
}

// GetSkillByName retrieves a skill by name
func (r *SkillRepository) GetSkillByName(name string) (*SkillDB, error) {
	query := `SELECT id, name, category, description, created_at FROM skills WHERE name = ?`
	row := r.db.Conn.QueryRow(query, name)

	skill := &SkillDB{}
	var createdAtStr string
	err := row.Scan(&skill.ID, &skill.Name, &skill.Category, &skill.Description, &createdAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan skill: %w", err)
	}

	if createdAtStr != "" {
		skill.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	}

	return skill, nil
}

// GetSkillByID retrieves a skill by ID
func (r *SkillRepository) GetSkillByID(id int64) (*SkillDB, error) {
	query := `SELECT id, name, category, description, created_at FROM skills WHERE id = ?`
	row := r.db.Conn.QueryRow(query, id)

	skill := &SkillDB{}
	var createdAtStr string
	err := row.Scan(&skill.ID, &skill.Name, &skill.Category, &skill.Description, &createdAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan skill: %w", err)
	}

	if createdAtStr != "" {
		skill.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	}

	return skill, nil
}

// UpsertSkillAndUserSkill upserts a skill and the user's proficiency in that skill
func (r *SkillRepository) UpsertSkillAndUserSkill(userID int64, skill models.Skill) error {
	// Normalize the skill name
	skillName := strings.TrimSpace(skill.Name)
	if skillName == "" {
		return nil
	}

	// Get or create the skill
	skillDB, err := r.GetOrCreateSkill(skillName, skill.Category, "")
	if err != nil {
		return err
	}

	// Check if user already has this skill
	var existingID int64
	var lastVerifiedStr string
	err = r.db.Conn.QueryRow(
		`SELECT id, last_verified_at FROM user_skills WHERE user_id = ? AND skill_id = ?`,
		userID, skillDB.ID,
	).Scan(&existingID, &lastVerifiedStr)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing user skill: %w", err)
	}

	// Parse last_verified_at
	var lastVerifiedAt *time.Time
	if lastVerifiedStr != "" {
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", lastVerifiedStr)
		lastVerifiedAt = &parsedTime
	}

	// Check if skill needs to be updated (expired or doesn't exist)
	needsUpdate := IsSkillExpired(lastVerifiedAt)

	if existingID > 0 && !needsUpdate {
		// Skill exists and is still valid, no need to update
		return nil
	}

	// Upsert the user skill
	if existingID > 0 {
		// Update existing
		query := `UPDATE user_skills SET proficiency = ?, years_experience = ?, last_verified_at = ?, updated_at = ? WHERE id = ?`
		_, err = r.db.Conn.Exec(query, skill.Proficiency, skill.Years, time.Now(), time.Now(), existingID)
	} else {
		// Insert new
		query := `INSERT INTO user_skills (user_id, skill_id, proficiency, years_experience, last_verified_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err = r.db.Conn.Exec(query, userID, skillDB.ID, skill.Proficiency, skill.Years, time.Now(), time.Now(), time.Now())
	}

	if err != nil {
		return fmt.Errorf("failed to upsert user skill: %w", err)
	}

	return nil
}

// GetUserSkills retrieves all skills for a user
func (r *SkillRepository) GetUserSkills(userID int64) ([]*UserSkillDB, error) {
	query := `
		SELECT us.id, us.user_id, us.skill_id, us.proficiency, us.years_experience, us.last_verified_at, us.created_at, us.updated_at, 
		       s.name, s.category
		FROM user_skills us
		JOIN skills s ON us.skill_id = s.id
		WHERE us.user_id = ?
		ORDER BY s.name`

	rows, err := r.db.Conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user skills: %w", err)
	}
	defer rows.Close()

	var skills []*UserSkillDB
	for rows.Next() {
		skill := &UserSkillDB{}
		var lastVerifiedStr, createdAtStr, updatedAtStr string
		if err := rows.Scan(
			&skill.ID, &skill.UserID, &skill.SkillID, &skill.Proficiency, &skill.Years,
			&lastVerifiedStr, &createdAtStr, &updatedAtStr, &skill.SkillName, &skill.Category,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user skill: %w", err)
		}

		// Parse timestamps
		if lastVerifiedStr != "" {
			parsedTime, _ := time.Parse("2006-01-02 15:04:05", lastVerifiedStr)
			skill.LastVerified = &parsedTime
		}
		if createdAtStr != "" {
			skill.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}
		if updatedAtStr != "" {
			skill.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		}

		skills = append(skills, skill)
	}

	return skills, nil
}

// GetExpiredSkills retrieves all expired skills for a user
func (r *SkillRepository) GetExpiredSkills(userID int64) ([]*UserSkillDB, error) {
	query := `
		SELECT us.id, us.user_id, us.skill_id, us.proficiency, us.years_experience, us.last_verified_at, us.created_at, us.updated_at, 
		       s.name, s.category
		FROM user_skills us
		JOIN skills s ON us.skill_id = s.id
		WHERE us.user_id = ? 
		AND (us.last_verified_at IS NULL OR datetime(us.last_verified_at) < datetime('now', '-1 year'))
		ORDER BY s.name`

	rows, err := r.db.Conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query expired user skills: %w", err)
	}
	defer rows.Close()

	var skills []*UserSkillDB
	for rows.Next() {
		skill := &UserSkillDB{}
		var lastVerifiedStr, createdAtStr, updatedAtStr string
		if err := rows.Scan(
			&skill.ID, &skill.UserID, &skill.SkillID, &skill.Proficiency, &skill.Years,
			&lastVerifiedStr, &createdAtStr, &updatedAtStr, &skill.SkillName, &skill.Category,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user skill: %w", err)
		}

		// Parse timestamps
		if lastVerifiedStr != "" {
			parsedTime, _ := time.Parse("2006-01-02 15:04:05", lastVerifiedStr)
			skill.LastVerified = &parsedTime
		}
		if createdAtStr != "" {
			skill.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}
		if updatedAtStr != "" {
			skill.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		}

		skills = append(skills, skill)
	}

	return skills, nil
}

// UpsertWorkExperience upserts work experience for a user
func (r *SkillRepository) UpsertWorkExperience(userID int64, exp models.WorkExperience) error {
	// Check if this experience already exists
	var existingID int64
	err := r.db.Conn.QueryRow(
		`SELECT id FROM work_experience WHERE user_id = ? AND company = ? AND position = ? AND start_date = ?`,
		userID, exp.Company, exp.Position, exp.StartDate.Format("2006-01-02"),
	).Scan(&existingID)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing work experience: %w", err)
	}

	// Format dates
	startDate := exp.StartDate.Format("2006-01-02")
	var endDate sql.NullString
	if exp.EndDate != nil {
		endDate.String = exp.EndDate.Format("2006-01-02")
		endDate.Valid = true
	}

	if existingID > 0 {
		// Update existing
		query := `UPDATE work_experience SET position = ?, end_date = ?, current = ?, description = ?, updated_at = ? WHERE id = ?`
		_, err = r.db.Conn.Exec(query, exp.Position, endDate, exp.Current, exp.Description, time.Now(), existingID)
	} else {
		// Insert new
		query := `INSERT INTO work_experience (user_id, company, position, start_date, end_date, current, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = r.db.Conn.Exec(query, userID, exp.Company, exp.Position, startDate, endDate, exp.Current, exp.Description, time.Now(), time.Now())
	}

	if err != nil {
		return fmt.Errorf("failed to upsert work experience: %w", err)
	}

	return nil
}

// UpsertEducation upserts education for a user
func (r *SkillRepository) UpsertEducation(userID int64, edu models.Education) error {
	// Check if this education already exists
	var existingID int64
	err := r.db.Conn.QueryRow(
		`SELECT id FROM education WHERE user_id = ? AND institution = ? AND degree = ? AND start_date = ?`,
		userID, edu.Institution, edu.Degree, edu.StartDate.Format("2006-01-02"),
	).Scan(&existingID)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing education: %w", err)
	}

	// Format dates
	startDate := edu.StartDate.Format("2006-01-02")
	var endDate sql.NullString
	if edu.EndDate != nil {
		endDate.String = edu.EndDate.Format("2006-01-02")
		endDate.Valid = true
	}

	if existingID > 0 {
		// Update existing
		query := `UPDATE education SET degree = ?, field_of_study = ?, end_date = ?, gpa = ?, updated_at = ? WHERE id = ?`
		_, err = r.db.Conn.Exec(query, edu.Degree, edu.FieldOfStudy, endDate, edu.GPA, time.Now(), existingID)
	} else {
		// Insert new
		query := `INSERT INTO education (user_id, institution, degree, field_of_study, start_date, end_date, gpa, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = r.db.Conn.Exec(query, userID, edu.Institution, edu.Degree, edu.FieldOfStudy, startDate, endDate, edu.GPA, time.Now(), time.Now())
	}

	if err != nil {
		return fmt.Errorf("failed to upsert education: %w", err)
	}

	return nil
}

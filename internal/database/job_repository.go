package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
)

// JobRepository handles job offer-related database operations
type JobRepository struct {
	db *Database
}

// NewJobRepository creates a new job repository
func NewJobRepository(db *Database) *JobRepository {
	return &JobRepository{db: db}
}

// JobOfferDB represents a job offer in the database
type JobOfferDB struct {
	ID          int64
	UserID      int64
	Company     string
	Position    string
	Location    string
	Description string
	Source      string
	RawContent  string
	CreatedAt   time.Time
}

// JobRequirementDB represents a job requirement in the database
type JobRequirementDB struct {
	ID              int64
	JobOfferID      int64
	RequirementType string
	Name            string
	Category        string
	IsRequired      bool
	MinYears        *int
	CreatedAt       time.Time
}

// UserResponseDB represents a user's response to a job requirement
type UserResponseDB struct {
	ID               int64
	UserID          int64
	JobRequirementID int64
	HasSkill        bool
	Proficiency     string
	Years           *int
	Notes           string
	RespondedAt     time.Time
	ExpiresAt       time.Time
	CreatedAt       time.Time
	RequirementName string
	RequirementType string
}

// StoreJobOffer stores a job offer in the database
func (r *JobRepository) StoreJobOffer(userID int64, jobOffer *models.JobOffer) (*JobOfferDB, error) {
	query := `INSERT INTO job_offers (user_id, company, position, location, description, source, raw_content, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Conn.Exec(
		query,
		userID,
		jobOffer.Company.Name,
		jobOffer.Position.Title,
		jobOffer.Position.Location,
		jobOffer.Description,
		jobOffer.Source,
		jobOffer.RawContent,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert job offer: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &JobOfferDB{
		ID:          id,
		UserID:      userID,
		Company:     jobOffer.Company.Name,
		Position:    jobOffer.Position.Title,
		Location:    jobOffer.Position.Location,
		Description: jobOffer.Description,
		Source:      jobOffer.Source,
		RawContent:  jobOffer.RawContent,
		CreatedAt:   time.Now(),
	}, nil
}

// GetJobOffer retrieves a job offer by ID
func (r *JobRepository) GetJobOffer(id int64) (*JobOfferDB, error) {
	query := `SELECT id, user_id, company, position, location, description, source, raw_content, created_at FROM job_offers WHERE id = ?`
	row := r.db.Conn.QueryRow(query, id)

	job := &JobOfferDB{}
	var createdAtStr string
	err := row.Scan(
		&job.ID, &job.UserID, &job.Company, &job.Position, &job.Location,
		&job.Description, &job.Source, &job.RawContent, &createdAtStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan job offer: %w", err)
	}

	if createdAtStr != "" {
		job.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	}

	return job, nil
}

// ListJobOffers lists all job offers for a user
func (r *JobRepository) ListJobOffers(userID int64) ([]*JobOfferDB, error) {
	query := `SELECT id, user_id, company, position, location, description, source, raw_content, created_at FROM job_offers WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := r.db.Conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query job offers: %w", err)
	}
	defer rows.Close()

	var jobs []*JobOfferDB
	for rows.Next() {
		job := &JobOfferDB{}
		var createdAtStr string
		if err := rows.Scan(
			&job.ID, &job.UserID, &job.Company, &job.Position, &job.Location,
			&job.Description, &job.Source, &job.RawContent, &createdAtStr,
		); err != nil {
			return nil, fmt.Errorf("failed to scan job offer: %w", err)
		}

		if createdAtStr != "" {
			job.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// StoreJobRequirements stores job requirements in the database
func (r *JobRepository) StoreJobRequirements(jobOfferID int64, requirements models.Requirements) ([]*JobRequirementDB, error) {
	var storedRequirements []*JobRequirementDB

	// Store required skills
	for _, skill := range requirements.RequiredSkills {
		req, err := r.storeRequirement(jobOfferID, "skill", skill.Name, skill.Category, true, skill.Years)
		if err != nil {
			return nil, err
		}
		storedRequirements = append(storedRequirements, req)
	}

	// Store preferred skills
	for _, skill := range requirements.PreferredSkills {
		req, err := r.storeRequirement(jobOfferID, "skill", skill.Name, skill.Category, false, skill.Years)
		if err != nil {
			return nil, err
		}
		storedRequirements = append(storedRequirements, req)
	}

	// Store experience requirements
	if requirements.Experience.TotalYears != nil {
		req, err := r.storeRequirement(jobOfferID, "experience", 
			fmt.Sprintf("%d years total experience", *requirements.Experience.TotalYears), 
			"", true, requirements.Experience.TotalYears)
		if err != nil {
			return nil, err
		}
		storedRequirements = append(storedRequirements, req)
	}

	if requirements.Experience.InField != nil {
		req, err := r.storeRequirement(jobOfferID, "experience",
			fmt.Sprintf("%d years in field", *requirements.Experience.InField),
			"", true, requirements.Experience.InField)
		if err != nil {
			return nil, err
		}
		storedRequirements = append(storedRequirements, req)
	}

	if requirements.Experience.Management {
		req, err := r.storeRequirement(jobOfferID, "experience", "Management experience", "", true, nil)
		if err != nil {
			return nil, err
		}
		storedRequirements = append(storedRequirements, req)
	}

	// Store education requirements
	for _, edu := range requirements.Education {
		req, err := r.storeRequirement(jobOfferID, "education", 
			fmt.Sprintf("%s in %s", edu.Degree, edu.FieldOfStudy), 
			"Education", edu.Required, nil)
		if err != nil {
			return nil, err
		}
		storedRequirements = append(storedRequirements, req)
	}

	// Store certifications
	for _, cert := range requirements.Certifications {
		req, err := r.storeRequirement(jobOfferID, "certification", cert, "", true, nil)
		if err != nil {
			return nil, err
		}
		storedRequirements = append(storedRequirements, req)
	}

	// Store languages
	for _, lang := range requirements.Languages {
		req, err := r.storeRequirement(jobOfferID, "language", 
			fmt.Sprintf("%s (%s)", lang.Language, lang.Proficiency), 
			"Language", lang.Required, nil)
		if err != nil {
			return nil, err
		}
		storedRequirements = append(storedRequirements, req)
	}

	// Store keywords as skills
	for _, keyword := range requirements.Keywords {
		// Skip if already stored
		alreadyStored := false
		for _, req := range storedRequirements {
			if strings.EqualFold(req.Name, keyword) {
				alreadyStored = true
				break
			}
		}
		if !alreadyStored {
			req, err := r.storeRequirement(jobOfferID, "skill", keyword, "", false, nil)
			if err != nil {
				return nil, err
			}
			storedRequirements = append(storedRequirements, req)
		}
	}

	return storedRequirements, nil
}

// storeRequirement stores a single requirement
func (r *JobRepository) storeRequirement(jobOfferID int64, reqType, name, category string, isRequired bool, minYears *int) (*JobRequirementDB, error) {
	query := `INSERT INTO job_requirements (job_offer_id, requirement_type, name, category, is_required, min_years, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Conn.Exec(query, jobOfferID, reqType, name, category, isRequired, minYears, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to insert job requirement: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &JobRequirementDB{
		ID:              id,
		JobOfferID:      jobOfferID,
		RequirementType: reqType,
		Name:            name,
		Category:        category,
		IsRequired:      isRequired,
		MinYears:        minYears,
		CreatedAt:       time.Now(),
	}, nil
}

// GetJobRequirements retrieves all requirements for a job offer
func (r *JobRepository) GetJobRequirements(jobOfferID int64) ([]*JobRequirementDB, error) {
	query := `SELECT id, job_offer_id, requirement_type, name, category, is_required, min_years, created_at FROM job_requirements WHERE job_offer_id = ? ORDER BY requirement_type, name`
	rows, err := r.db.Conn.Query(query, jobOfferID)
	if err != nil {
		return nil, fmt.Errorf("failed to query job requirements: %w", err)
	}
	defer rows.Close()

	var requirements []*JobRequirementDB
	for rows.Next() {
		req := &JobRequirementDB{}
		var createdAtStr string
		var minYears sql.NullInt64
		if err := rows.Scan(
			&req.ID, &req.JobOfferID, &req.RequirementType, &req.Name, &req.Category,
			&req.IsRequired, &minYears, &createdAtStr,
		); err != nil {
			return nil, fmt.Errorf("failed to scan job requirement: %w", err)
		}

		if minYears.Valid {
			years := int(minYears.Int64)
			req.MinYears = &years
		}

		if createdAtStr != "" {
			req.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}

		requirements = append(requirements, req)
	}

	return requirements, nil
}

// StoreUserResponse stores a user's response to a job requirement
func (r *JobRepository) StoreUserResponse(userID int64, jobRequirementID int64, hasSkill bool, proficiency, yearsStr, notes string) (*UserResponseDB, error) {
	// Check if response already exists
	var existingID int64
	err := r.db.Conn.QueryRow(
		`SELECT id FROM user_responses WHERE user_id = ? AND job_requirement_id = ?`,
		userID, jobRequirementID,
	).Scan(&existingID)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing user response: %w", err)
	}

	// Parse years
	var years *int
	if yearsStr != "" {
		var y int
		if _, err := fmt.Sscanf(yearsStr, "%d", &y); err == nil {
			years = &y
		}
	}

	// Get expiration date (1 year from now)
	expiresAt := GetExpirationDate()

	var responseID int64
	if existingID > 0 {
		// Update existing response
		query := `UPDATE user_responses SET has_skill = ?, proficiency = ?, years_experience = ?, notes = ?, responded_at = ?, expires_at = ? WHERE id = ?`
		_, err := r.db.Conn.Exec(query, hasSkill, proficiency, years, notes, time.Now(), expiresAt, existingID)
		if err != nil {
			return nil, fmt.Errorf("failed to update user response: %w", err)
		}
		responseID = existingID
	} else {
		// Insert new response
		query := `INSERT INTO user_responses (user_id, job_requirement_id, has_skill, proficiency, years_experience, notes, responded_at, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := r.db.Conn.Exec(query, userID, jobRequirementID, hasSkill, proficiency, years, notes, time.Now(), expiresAt, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to insert user response: %w", err)
		}
		responseID, err = result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get last insert id: %w", err)
		}
	}

	// Get requirement info for the response
	var req JobRequirementDB
	var createdAtStr string
	err = r.db.Conn.QueryRow(
		`SELECT jr.id, jr.job_offer_id, jr.requirement_type, jr.name, jr.category, jr.is_required, jr.created_at FROM job_requirements jr WHERE jr.id = ?`,
		jobRequirementID,
	).Scan(&req.ID, &req.JobOfferID, &req.RequirementType, &req.Name, &req.Category, &req.IsRequired, &createdAtStr)

	if err != nil {
		return nil, fmt.Errorf("failed to get requirement info: %w", err)
	}

	return &UserResponseDB{
		ID:               responseID,
		UserID:          userID,
		JobRequirementID: jobRequirementID,
		HasSkill:        hasSkill,
		Proficiency:     proficiency,
		Years:           years,
		Notes:           notes,
		RespondedAt:     time.Now(),
		ExpiresAt:       expiresAt,
		CreatedAt:       time.Now(),
		RequirementName: req.Name,
		RequirementType: req.RequirementType,
	}, nil
}

// GetUserResponses retrieves all responses for a user
func (r *JobRepository) GetUserResponses(userID int64) ([]*UserResponseDB, error) {
	query := `
		SELECT ur.id, ur.user_id, ur.job_requirement_id, ur.has_skill, ur.proficiency, ur.years_experience, ur.notes, ur.responded_at, ur.expires_at, ur.created_at,
		       jr.name, jr.requirement_type
		FROM user_responses ur
		JOIN job_requirements jr ON ur.job_requirement_id = jr.id
		WHERE ur.user_id = ?
		ORDER BY ur.responded_at DESC`

	rows, err := r.db.Conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user responses: %w", err)
	}
	defer rows.Close()

	var responses []*UserResponseDB
	for rows.Next() {
		response := &UserResponseDB{}
		var respondedAtStr, expiresAtStr, createdAtStr string
		var years sql.NullInt64
		if err := rows.Scan(
			&response.ID, &response.UserID, &response.JobRequirementID, &response.HasSkill,
			&response.Proficiency, &years, &response.Notes, &respondedAtStr,
			&expiresAtStr, &createdAtStr, &response.RequirementName, &response.RequirementType,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user response: %w", err)
		}

		// Parse years
		if years.Valid {
			y := int(years.Int64)
			response.Years = &y
		}

		// Parse timestamps
		if respondedAtStr != "" {
			response.RespondedAt, _ = time.Parse("2006-01-02 15:04:05", respondedAtStr)
		}
		if expiresAtStr != "" {
			response.ExpiresAt, _ = time.Parse("2006-01-02 15:04:05", expiresAtStr)
		}
		if createdAtStr != "" {
			response.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// GetExpiredResponses retrieves all expired responses for a user
func (r *JobRepository) GetExpiredResponses(userID int64) ([]*UserResponseDB, error) {
	query := `
		SELECT ur.id, ur.user_id, ur.job_requirement_id, ur.has_skill, ur.proficiency, ur.years_experience, ur.notes, ur.responded_at, ur.expires_at, ur.created_at,
		       jr.name, jr.requirement_type
		FROM user_responses ur
		JOIN job_requirements jr ON ur.job_requirement_id = jr.id
		WHERE ur.user_id = ? 
		AND datetime(ur.expires_at) < datetime('now')
		ORDER BY ur.responded_at DESC`

	rows, err := r.db.Conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query expired user responses: %w", err)
	}
	defer rows.Close()

	var responses []*UserResponseDB
	for rows.Next() {
		response := &UserResponseDB{}
		var respondedAtStr, expiresAtStr, createdAtStr string
		var years sql.NullInt64
		if err := rows.Scan(
			&response.ID, &response.UserID, &response.JobRequirementID, &response.HasSkill,
			&response.Proficiency, &years, &response.Notes, &respondedAtStr,
			&expiresAtStr, &createdAtStr, &response.RequirementName, &response.RequirementType,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user response: %w", err)
		}

		// Parse years
		if years.Valid {
			y := int(years.Int64)
			response.Years = &y
		}

		// Parse timestamps
		if respondedAtStr != "" {
			response.RespondedAt, _ = time.Parse("2006-01-02 15:04:05", respondedAtStr)
		}
		if expiresAtStr != "" {
			response.ExpiresAt, _ = time.Parse("2006-01-02 15:04:05", expiresAtStr)
		}
		if createdAtStr != "" {
			response.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// GetJobOfferWithRequirements retrieves a job offer with all its requirements
func (r *JobRepository) GetJobOfferWithRequirements(jobOfferID int64) (*JobOfferDB, []*JobRequirementDB, error) {
	job, err := r.GetJobOffer(jobOfferID)
	if err != nil {
		return nil, nil, err
	}
	if job == nil {
		return nil, nil, fmt.Errorf("job offer not found")
	}

	requirements, err := r.GetJobRequirements(jobOfferID)
	if err != nil {
		return nil, nil, err
	}

	return job, requirements, nil
}

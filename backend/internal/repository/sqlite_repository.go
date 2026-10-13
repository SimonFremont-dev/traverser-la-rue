package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SimonFremont-dev/cv-builder/backend/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteRepository handles SQLite database operations for personal data
type SQLiteRepository struct {
	db     *sql.DB
	dbPath string
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
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

	repo := &SQLiteRepository{
		db:     db,
		dbPath: dbPath,
	}

	// Initialize schema
	if err := repo.InitSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

// InitSchema initializes the SQLite database schema
func (r *SQLiteRepository) InitSchema() error {
	schema := `
	-- Profile table
	CREATE TABLE IF NOT EXISTS profiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Profile skills (links to standardized skills in PostgreSQL)
	CREATE TABLE IF NOT EXISTS profile_skills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_id INTEGER NOT NULL,
		skill_id INTEGER NOT NULL,
		skill_name TEXT NOT NULL,
		skill_source_id INTEGER NOT NULL,
		skill_type_id INTEGER NOT NULL,
		proficiency TEXT CHECK(proficiency IN ('Beginner', 'Intermediate', 'Advanced', 'Expert')),
		years_experience INTEGER DEFAULT 0,
		last_verified TIMESTAMP,
		notes TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
		UNIQUE(profile_id, skill_id)
	);

	-- Experience table
	CREATE TABLE IF NOT EXISTS experience (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_id INTEGER NOT NULL,
		company TEXT NOT NULL,
		position TEXT NOT NULL,
		start_date TEXT NOT NULL,
		end_date TEXT,
		current BOOLEAN DEFAULT FALSE,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
	);

	-- Education table
	CREATE TABLE IF NOT EXISTS education (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_id INTEGER NOT NULL,
		institution TEXT NOT NULL,
		degree TEXT NOT NULL,
		field_of_study TEXT,
		start_date TEXT,
		end_date TEXT,
		gpa REAL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
	);

	-- Job offers table (local cache)
	CREATE TABLE IF NOT EXISTS job_offers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		company TEXT NOT NULL,
		description TEXT,
		requirements TEXT,
		location TEXT,
		url TEXT,
		source TEXT,
		raw_content TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Keywords extracted from job offers
	CREATE TABLE IF NOT EXISTS keywords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_offer_id INTEGER NOT NULL,
		keyword TEXT NOT NULL,
		category TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (job_offer_id) REFERENCES job_offers(id) ON DELETE CASCADE,
		UNIQUE(job_offer_id, keyword)
	);

	-- Resumes table
	CREATE TABLE IF NOT EXISTS resumes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_id INTEGER NOT NULL,
		job_offer_id INTEGER,
		content TEXT NOT NULL,
		format TEXT NOT NULL DEFAULT 'text',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
		FOREIGN KEY (job_offer_id) REFERENCES job_offers(id) ON DELETE SET NULL
	);

	-- Custom skills table (user-defined skills not in standardized dictionaries)
	CREATE TABLE IF NOT EXISTS custom_skills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_id INTEGER NOT NULL,
		name TEXT NOT NULL COLLATE NOCASE,
		category TEXT,
		proficiency TEXT CHECK(proficiency IN ('Beginner', 'Intermediate', 'Advanced', 'Expert')),
		years_experience INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
		UNIQUE(profile_id, name COLLATE NOCASE)
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_profiles_email ON profiles(email);
	CREATE INDEX IF NOT EXISTS idx_profile_skills_profile ON profile_skills(profile_id);
	CREATE INDEX IF NOT EXISTS idx_profile_skills_skill ON profile_skills(skill_id);
	CREATE INDEX IF NOT EXISTS idx_experience_profile ON experience(profile_id);
	CREATE INDEX IF NOT EXISTS idx_education_profile ON education(profile_id);
	CREATE INDEX IF NOT EXISTS idx_keywords_job ON keywords(job_offer_id);
	CREATE INDEX IF NOT EXISTS idx_keywords_keyword ON keywords(keyword);
	CREATE INDEX IF NOT EXISTS idx_resumes_profile ON resumes(profile_id);
	CREATE INDEX IF NOT EXISTS idx_custom_skills_profile ON custom_skills(profile_id);
	CREATE INDEX IF NOT EXISTS idx_custom_skills_name ON custom_skills(LOWER(name));
	`

	_, err := r.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// GetDBPath returns the database path
func (r *SQLiteRepository) GetDBPath() string {
	return r.dbPath
}

// --- Profile operations ---

// CreateProfile creates a new profile
func (r *SQLiteRepository) CreateProfile(profile *models.Profile) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert profile
	result, err := tx.Exec(
		"INSERT INTO profiles (name, email) VALUES (?, ?)",
		profile.Name, profile.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to insert profile: %w", err)
	}

	profileID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get profile ID: %w", err)
	}
	profile.ID = int(profileID)

	// Insert skills
	for _, skill := range profile.Skills {
		_, err := tx.Exec(
			"INSERT INTO profile_skills (profile_id, skill_id, skill_name, skill_source_id, skill_type_id, proficiency, years_experience, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			profileID, skill.SkillID, skill.SkillName, skill.SourceID, skill.TypeID, skill.Proficiency, skill.YearsExperience, skill.Notes,
		)
		if err != nil {
			return fmt.Errorf("failed to insert skill: %w", err)
		}
	}

	// Insert experience
	for _, exp := range profile.Experience {
		_, err := tx.Exec(
			"INSERT INTO experience (profile_id, company, position, start_date, end_date, current, description) VALUES (?, ?, ?, ?, ?, ?, ?)",
			profileID, exp.Company, exp.Position, exp.StartDate, exp.EndDate, exp.Current, exp.Description,
		)
		if err != nil {
			return fmt.Errorf("failed to insert experience: %w", err)
		}
	}

	// Insert education
	for _, edu := range profile.Education {
		_, err := tx.Exec(
			"INSERT INTO education (profile_id, institution, degree, field_of_study, start_date, end_date, gpa) VALUES (?, ?, ?, ?, ?, ?, ?)",
			profileID, edu.Institution, edu.Degree, edu.FieldOfStudy, edu.StartDate, edu.EndDate, edu.GPA,
		)
		if err != nil {
			return fmt.Errorf("failed to insert education: %w", err)
		}
	}

	// Insert custom skills
	for _, skill := range profile.CustomSkills {
		_, err := tx.Exec(
			"INSERT INTO custom_skills (profile_id, name, category, proficiency, years_experience) VALUES (?, ?, ?, ?, ?)",
			profileID, skill.Name, skill.Category, skill.Proficiency, skill.YearsExperience,
		)
		if err != nil {
			return fmt.Errorf("failed to insert custom skill: %w", err)
		}
	}

	return tx.Commit()
}

// GetProfile retrieves a profile by ID
func (r *SQLiteRepository) GetProfile(id int) (*models.Profile, error) {
	profile := &models.Profile{}

	// Get profile
	err := r.db.QueryRow(
		"SELECT id, name, email FROM profiles WHERE id = ?", id,
	).Scan(&profile.ID, &profile.Name, &profile.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Get skills
	rows, err := r.db.Query(
		"SELECT id, skill_id, skill_name, skill_source_id, skill_type_id, proficiency, years_experience, last_verified, notes FROM profile_skills WHERE profile_id = ?", id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		skill := models.ProfileSkill{}
		var lastVerified sql.NullString
		if err := rows.Scan(&skill.ID, &skill.SkillID, &skill.SkillName, &skill.SourceID, &skill.TypeID, &skill.Proficiency, &skill.YearsExperience, &lastVerified, &skill.Notes); err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}
		if lastVerified.Valid {
			skill.LastVerified = lastVerified.String
		}
		profile.Skills = append(profile.Skills, skill)
	}

	// Get experience
	rows, err = r.db.Query(
		"SELECT id, company, position, start_date, end_date, current, description FROM experience WHERE profile_id = ?", id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get experience: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		exp := models.ProfileExperience{}
		if err := rows.Scan(&exp.ID, &exp.Company, &exp.Position, &exp.StartDate, &exp.EndDate, &exp.Current, &exp.Description); err != nil {
			return nil, fmt.Errorf("failed to scan experience: %w", err)
		}
		profile.Experience = append(profile.Experience, exp)
	}

	// Get education
	rows, err = r.db.Query(
		"SELECT id, institution, degree, field_of_study, start_date, end_date, gpa FROM education WHERE profile_id = ?", id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get education: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		edu := models.ProfileEducation{}
		if err := rows.Scan(&edu.ID, &edu.Institution, &edu.Degree, &edu.FieldOfStudy, &edu.StartDate, &edu.EndDate, &edu.GPA); err != nil {
			return nil, fmt.Errorf("failed to scan education: %w", err)
		}
		profile.Education = append(profile.Education, edu)
	}

	// Get custom skills
	rows, err = r.db.Query(
		"SELECT id, name, category, proficiency, years_experience, created_at FROM custom_skills WHERE profile_id = ?", id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom skills: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		skill := models.CustomSkill{}
		if err := rows.Scan(&skill.ID, &skill.Name, &skill.Category, &skill.Proficiency, &skill.YearsExperience, &skill.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan custom skill: %w", err)
		}
		profile.CustomSkills = append(profile.CustomSkills, skill)
	}

	return profile, nil
}

// UpdateProfile updates a profile
func (r *SQLiteRepository) UpdateProfile(profile *models.Profile) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update profile
	_, err = tx.Exec(
		"UPDATE profiles SET name = ?, email = ? WHERE id = ?",
		profile.Name, profile.Email, profile.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	// Delete existing profile_skills and re-insert
	_, err = tx.Exec("DELETE FROM profile_skills WHERE profile_id = ?", profile.ID)
	if err != nil {
		return fmt.Errorf("failed to delete profile_skills: %w", err)
	}

	for _, skill := range profile.Skills {
		_, err := tx.Exec(
			"INSERT INTO profile_skills (profile_id, skill_id, skill_name, skill_source_id, skill_type_id, proficiency, years_experience, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			profile.ID, skill.SkillID, skill.SkillName, skill.SourceID, skill.TypeID, skill.Proficiency, skill.YearsExperience, skill.Notes,
		)
		if err != nil {
			return fmt.Errorf("failed to insert skill: %w", err)
		}
	}

	// Delete existing experience and re-insert
	_, err = tx.Exec("DELETE FROM experience WHERE profile_id = ?", profile.ID)
	if err != nil {
		return fmt.Errorf("failed to delete experience: %w", err)
	}

	for _, exp := range profile.Experience {
		_, err := tx.Exec(
			"INSERT INTO experience (profile_id, company, position, start_date, end_date, current, description) VALUES (?, ?, ?, ?, ?, ?, ?)",
			profile.ID, exp.Company, exp.Position, exp.StartDate, exp.EndDate, exp.Current, exp.Description,
		)
		if err != nil {
			return fmt.Errorf("failed to insert experience: %w", err)
		}
	}

	// Delete existing education and re-insert
	_, err = tx.Exec("DELETE FROM education WHERE profile_id = ?", profile.ID)
	if err != nil {
		return fmt.Errorf("failed to delete education: %w", err)
	}

	for _, edu := range profile.Education {
		_, err := tx.Exec(
			"INSERT INTO education (profile_id, institution, degree, field_of_study, start_date, end_date, gpa) VALUES (?, ?, ?, ?, ?, ?, ?)",
			profile.ID, edu.Institution, edu.Degree, edu.FieldOfStudy, edu.StartDate, edu.EndDate, edu.GPA,
		)
		if err != nil {
			return fmt.Errorf("failed to insert education: %w", err)
		}
	}

	// Delete existing custom skills and re-insert
	_, err = tx.Exec("DELETE FROM custom_skills WHERE profile_id = ?", profile.ID)
	if err != nil {
		return fmt.Errorf("failed to delete custom skills: %w", err)
	}

	for _, skill := range profile.CustomSkills {
		_, err := tx.Exec(
			"INSERT INTO custom_skills (profile_id, name, category, proficiency, years_experience) VALUES (?, ?, ?, ?, ?)",
			profile.ID, skill.Name, skill.Category, skill.Proficiency, skill.YearsExperience,
		)
		if err != nil {
			return fmt.Errorf("failed to insert custom skill: %w", err)
		}
	}

	return tx.Commit()
}

// AddSkill adds a skill to a profile
func (r *SQLiteRepository) AddSkill(profileID int, skill *models.ProfileSkill) error {
	_, err := r.db.Exec(
		"INSERT INTO profile_skills (profile_id, skill_id, skill_name, skill_source_id, skill_type_id, proficiency, years_experience, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		profileID, skill.SkillID, skill.SkillName, skill.SourceID, skill.TypeID, skill.Proficiency, skill.YearsExperience, skill.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to add skill: %w", err)
	}
	return nil
}

// GetProfileSkills returns all skills for a profile
func (r *SQLiteRepository) GetProfileSkills(profileID int) ([]models.ProfileSkill, error) {
	var skills []models.ProfileSkill

	rows, err := r.db.Query(
		"SELECT id, skill_id, skill_name, skill_source_id, skill_type_id, proficiency, years_experience, last_verified, notes FROM profile_skills WHERE profile_id = ?", profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		skill := models.ProfileSkill{}
		var lastVerified sql.NullString
		if err := rows.Scan(&skill.ID, &skill.SkillID, &skill.SkillName, &skill.SourceID, &skill.TypeID, &skill.Proficiency, &skill.YearsExperience, &lastVerified, &skill.Notes); err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}
		if lastVerified.Valid {
			skill.LastVerified = lastVerified.String
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

// CheckSkillExists checks if a profile has a specific skill by skill_id
func (r *SQLiteRepository) CheckSkillExists(profileID int, skillID int) (bool, error) {
	var exists int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM profile_skills WHERE profile_id = ? AND skill_id = ?",
		profileID, skillID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check skill: %w", err)
	}
	return exists > 0, nil
}

// SaveJobOffer saves a job offer
func (r *SQLiteRepository) SaveJobOffer(offer *models.JobOffer) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert job offer
	result, err := tx.Exec(
		"INSERT INTO job_offers (title, company, description, requirements, location, url, source, raw_content) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		offer.Title, offer.Company, offer.Description, offer.Requirements, offer.Location, offer.URL, offer.Source, offer.RawContent,
	)
	if err != nil {
		return fmt.Errorf("failed to insert job offer: %w", err)
	}

	offerID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get job offer ID: %w", err)
	}
	offer.ID = int(offerID)

	return tx.Commit()
}

// SaveKeywords saves extracted keywords for a job offer
func (r *SQLiteRepository) SaveKeywords(jobOfferID int, keywords []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Clear existing keywords for this job offer
	_, err = tx.Exec("DELETE FROM keywords WHERE job_offer_id = ?", jobOfferID)
	if err != nil {
		return fmt.Errorf("failed to delete keywords: %w", err)
	}

	// Insert new keywords
	for _, keyword := range keywords {
		_, err := tx.Exec(
			"INSERT INTO keywords (job_offer_id, keyword) VALUES (?, ?)",
			jobOfferID, keyword,
		)
		if err != nil {
			return fmt.Errorf("failed to insert keyword: %w", err)
		}
	}

	return tx.Commit()
}

// GetKeywords returns keywords for a job offer
func (r *SQLiteRepository) GetKeywords(jobOfferID int) ([]string, error) {
	var keywords []string

	rows, err := r.db.Query(
		"SELECT keyword FROM keywords WHERE job_offer_id = ?", jobOfferID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get keywords: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			return nil, fmt.Errorf("failed to scan keyword: %w", err)
		}
		keywords = append(keywords, keyword)
	}

	return keywords, nil
}

// SaveResume saves a generated resume
func (r *SQLiteRepository) SaveResume(resume *models.Resume) error {
	_, err := r.db.Exec(
		"INSERT INTO resumes (profile_id, job_offer_id, content, format) VALUES (?, ?, ?, ?)",
		resume.ProfileID, resume.JobOfferID, resume.Content, resume.Format,
	)
	if err != nil {
		return fmt.Errorf("failed to save resume: %w", err)
	}
	return nil
}

// GetAllProfiles returns all profiles
func (r *SQLiteRepository) GetAllProfiles() ([]*models.Profile, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM profiles")
	if err != nil {
		return nil, fmt.Errorf("failed to get profiles: %w", err)
	}
	defer rows.Close()

	var profiles []*models.Profile
	for rows.Next() {
		profile := &models.Profile{}
		if err := rows.Scan(&profile.ID, &profile.Name, &profile.Email); err != nil {
			return nil, fmt.Errorf("failed to scan profile: %w", err)
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// CheckCustomSkillExists checks if a custom skill exists for a profile
func (r *SQLiteRepository) CheckCustomSkillExists(profileID int, skillName string) (bool, error) {
	var exists int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM custom_skills WHERE profile_id = ? AND LOWER(name) = LOWER(?)",
		profileID, skillName,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check custom skill: %w", err)
	}
	return exists > 0, nil
}

// AddCustomSkill adds a custom skill to a profile
func (r *SQLiteRepository) AddCustomSkill(profileID int, skill *models.CustomSkill) error {
	_, err := r.db.Exec(
		"INSERT INTO custom_skills (profile_id, name, category, proficiency, years_experience) VALUES (?, ?, ?, ?, ?)",
		profileID, skill.Name, skill.Category, skill.Proficiency, skill.YearsExperience,
	)
	if err != nil {
		return fmt.Errorf("failed to add custom skill: %w", err)
	}
	return nil
}

// GetProfileCustomSkills returns all custom skills for a profile
func (r *SQLiteRepository) GetProfileCustomSkills(profileID int) ([]models.CustomSkill, error) {
	var skills []models.CustomSkill

	rows, err := r.db.Query(
		"SELECT id, name, category, proficiency, years_experience, created_at FROM custom_skills WHERE profile_id = ?", profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom skills: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		skill := models.CustomSkill{}
		if err := rows.Scan(&skill.ID, &skill.Name, &skill.Category, &skill.Proficiency, &skill.YearsExperience, &skill.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan custom skill: %w", err)
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

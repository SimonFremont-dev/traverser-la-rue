package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/SimonFremont-dev/cv-builder/backend/internal/models"
)

// SkillSQLRepository handles skill dictionary operations in SQL database
type SkillSQLRepository struct {
	db *sql.DB
}

// NewSkillSQLRepository creates a new skill SQL repository
func NewSkillSQLRepository(db *sql.DB) *SkillSQLRepository {
	return &SkillSQLRepository{db: db}
}

// GetAllSources returns all skill sources
func (r *SkillSQLRepository) GetAllSources() ([]*models.Source, error) {
	var sources []*models.Source

	rows, err := r.db.Query(
		`SELECT id, name, description, is_official, created_at, updated_at 
		 FROM source ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get sources: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		s := &models.Source{}
		if err := rows.Scan(
			&s.ID, &s.Name, &s.Description, &s.IsOfficial, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan source: %w", err)
		}
		sources = append(sources, s)
	}

	return sources, nil
}

// GetSourceByName returns a source by name
func (r *SkillSQLRepository) GetSourceByName(name string) (*models.Source, error) {
	s := &models.Source{}
	err := r.db.QueryRow(
		`SELECT id, name, description, is_official, created_at, updated_at 
		 FROM source WHERE LOWER(name) = LOWER($1)`,
		name,
	).Scan(
		&s.ID, &s.Name, &s.Description, &s.IsOfficial, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("source not found: %s", name)
		}
		return nil, fmt.Errorf("failed to get source: %w", err)
	}
	return s, nil
}

// GetOfficialSources returns all official sources
func (r *SkillSQLRepository) GetOfficialSources() ([]*models.Source, error) {
	var sources []*models.Source

	rows, err := r.db.Query(
		`SELECT id, name, description, is_official, created_at, updated_at 
		 FROM source WHERE is_official = TRUE ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get official sources: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		s := &models.Source{}
		if err := rows.Scan(
			&s.ID, &s.Name, &s.Description, &s.IsOfficial, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan source: %w", err)
		}
		sources = append(sources, s)
	}

	return sources, nil
}

// GetAllTypes returns all skill types
func (r *SkillSQLRepository) GetAllTypes() ([]*models.Type, error) {
	var types []*models.Type

	rows, err := r.db.Query(
		`SELECT id, name, description, created_at 
		 FROM type ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get types: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := &models.Type{}
		if err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &t.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan type: %w", err)
		}
		types = append(types, t)
	}

	return types, nil
}

// GetTypeByName returns a type by name
func (r *SkillSQLRepository) GetTypeByName(name string) (*models.Type, error) {
	t := &models.Type{}
	err := r.db.QueryRow(
		`SELECT id, name, description, created_at 
		 FROM type WHERE LOWER(name) = LOWER($1)`,
		name,
	).Scan(
		&t.ID, &t.Name, &t.Description, &t.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("type not found: %s", name)
		}
		return nil, fmt.Errorf("failed to get type: %w", err)
	}
	return t, nil
}

// GetSkillByID returns a skill by ID
func (r *SkillSQLRepository) GetSkillByID(id int) (*models.Skill, error) {
	s := &models.Skill{}
	err := r.db.QueryRow(
		`SELECT id, source_id, type_id, name, external_id, description, created_at, updated_at 
		 FROM skill WHERE id = $1`,
		id,
	).Scan(
		&s.ID, &s.SourceID, &s.TypeID, &s.Name, &s.ExternalID, &s.Description, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("skill not found")
		}
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}
	return s, nil
}

// GetSkillByNameAndSource returns a skill by name and source
func (r *SkillSQLRepository) GetSkillByNameAndSource(name string, sourceName string) (*models.Skill, error) {
	s := &models.Skill{}
	err := r.db.QueryRow(
		`SELECT s.id, s.source_id, s.type_id, s.name, s.external_id, s.description, s.created_at, s.updated_at 
		 FROM skill s
		 JOIN source src ON s.source_id = src.id
		 WHERE LOWER(s.name) = LOWER($1) AND LOWER(src.name) = LOWER($2)`,
		name, sourceName,
	).Scan(
		&s.ID, &s.SourceID, &s.TypeID, &s.Name, &s.ExternalID, &s.Description, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("skill not found: %s from %s", name, sourceName)
		}
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}
	return s, nil
}

// SearchSkills searches for skills by name
func (r *SkillSQLRepository) SearchSkills(filter *models.SkillFilter) ([]*models.SkillWithDetails, error) {
	var skills []*models.SkillWithDetails
	
	// Build query
	query := `
		SELECT s.id, s.source_id, s.type_id, s.name, s.external_id, s.description, s.created_at, s.updated_at,
		       src.name as source_name, src.description as source_description, src.is_official,
		       t.name as type_name, t.description as type_description
		FROM skill s
		JOIN source src ON s.source_id = src.id
		JOIN type t ON s.type_id = t.id
	`
	
	var conditions []string
	var args []interface{}
	argNum := 1
	
	// Add source filter
	if len(filter.SourceIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("s.source_id = ANY($%d)", argNum))
		args = append(args, "{" + intSliceToString(filter.SourceIDs) + "}")
		argNum++
	} else if filter.IsOfficial != nil {
		if *filter.IsOfficial {
			conditions = append(conditions, "src.is_official = TRUE")
		} else {
			conditions = append(conditions, "src.is_official = FALSE")
		}
	}
	
	// Add type filter
	if len(filter.TypeIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("s.type_id = ANY($%d)", argNum))
		args = append(args, "{" + intSliceToString(filter.TypeIDs) + "}")
		argNum++
	}
	
	// Add search filter
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(s.name) LIKE LOWER($%d) OR LOWER(s.description) LIKE LOWER($%d)", argNum, argNum))
		args = append(args, "%"+filter.Search+"%")
		argNum++
	}
	
	// Add conditions to query
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	// Add ordering and limit
	query += " ORDER BY s.name"
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argNum)
		args = append(args, filter.Limit)
	}
	
	// Execute query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search skills: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		s := &models.SkillWithDetails{}
		if err := rows.Scan(
			&s.ID, &s.SourceID, &s.TypeID, &s.Name, &s.ExternalID, &s.Description,
			&s.CreatedAt, &s.UpdatedAt,
			&s.Source.Name, &s.Source.Description, &s.Source.IsOfficial,
			&s.Type.Name, &s.Type.Description,
		); err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}
		skills = append(skills, s)
	}

	return skills, nil
}

// GetAllSkills returns all skills
func (r *SkillSQLRepository) GetAllSkills() ([]*models.SkillWithDetails, error) {
	return r.SearchSkills(&models.SkillFilter{})
}

// GetOfficialSkills returns all skills from official sources
func (r *SkillSQLRepository) GetOfficialSkills() ([]*models.SkillWithDetails, error) {
	official := true
	return r.SearchSkills(&models.SkillFilter{
		IsOfficial: &official,
	})
}

// CreateSkill creates a new skill
func (r *SkillSQLRepository) CreateSkill(skill *models.Skill) error {
	_, err := r.db.Exec(
		`INSERT INTO skill (source_id, type_id, name, external_id, description) 
		 VALUES ($1, $2, $3, $4, $5)`,
		skill.SourceID, skill.TypeID, skill.Name, skill.ExternalID, skill.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to create skill: %w", err)
	}
	return nil
}

// CreateSource creates a new source
func (r *SkillSQLRepository) CreateSource(source *models.Source) error {
	_, err := r.db.Exec(
		`INSERT INTO source (name, description, is_official) 
		 VALUES ($1, $2, $3)`,
		source.Name, source.Description, source.IsOfficial,
	)
	if err != nil {
		return fmt.Errorf("failed to create source: %w", err)
	}
	return nil
}

// CreateType creates a new type
func (r *SkillSQLRepository) CreateType(typ *models.Type) error {
	_, err := r.db.Exec(
		`INSERT INTO type (name, description) 
		 VALUES ($1, $2)`,
		typ.Name, typ.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to create type: %w", err)
	}
	return nil
}

// BulkInsertSkills inserts multiple skills efficiently
func (r *SkillSQLRepository) BulkInsertSkills(skills []*models.Skill) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, skill := range skills {
		_, err := tx.Exec(
			`INSERT INTO skill (source_id, type_id, name, external_id, description) 
			 VALUES ($1, $2, $3, $4, $5) 
			 ON CONFLICT (source_id, LOWER(name)) DO NOTHING`,
			skill.SourceID, skill.TypeID, skill.Name, skill.ExternalID, skill.Description,
		)
		if err != nil {
			return fmt.Errorf("failed to insert skill: %w", err)
		}
	}

	return tx.Commit()
}

// FindSkillByKeyword finds skills matching a keyword
func (r *SkillSQLRepository) FindSkillByKeyword(keyword string) ([]*models.StandardizedSkill, error) {
	var results []*models.StandardizedSkill

	// First, try exact match on skill name
	query := `
		SELECT s.id, s.source_id, s.type_id, s.name, s.external_id, s.description, 
		       src.name as source_name, src.is_official,
		       t.name as type_name
		FROM skill s
		JOIN source src ON s.source_id = src.id
		JOIN type t ON s.type_id = t.id
		WHERE LOWER(s.name) = LOWER($1)
		ORDER BY src.is_official DESC, s.name
	`
	
	rows, err := r.db.Query(query, keyword)
	if err != nil {
		return nil, fmt.Errorf("failed to query exact matches: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		s := &models.Skill{}
		var sourceName string
		var isOfficial bool
		var typeName string
		
		if err := rows.Scan(
			&s.ID, &s.SourceID, &s.TypeID, &s.Name, &s.ExternalID, &s.Description,
			&sourceName, &isOfficial, &typeName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}
		
		// Set transient fields
		s.SourceName = sourceName
		s.TypeName = typeName
		s.IsOfficial = isOfficial
		
		results = append(results, &models.StandardizedSkill{
			Keyword:    keyword,
			Skill:      *s,
			Confidence: 1.0,
		})
	}

	// If no exact match, try partial match
	if len(results) == 0 {
		query := `
			SELECT s.id, s.source_id, s.type_id, s.name, s.external_id, s.description,
			       src.name as source_name, src.is_official,
			       t.name as type_name
			FROM skill s
			JOIN source src ON s.source_id = src.id
			JOIN type t ON s.type_id = t.id
			WHERE LOWER(s.name) LIKE LOWER($1) OR LOWER(s.description) LIKE LOWER($1)
			ORDER BY src.is_official DESC, LEVENSHTEIN(LOWER(s.name), LOWER($1)) LIMIT 10
		`
		
		rows, err := r.db.Query(query, "%"+keyword+"%", keyword)
		if err != nil {
			return nil, fmt.Errorf("failed to query partial matches: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			s := &models.Skill{}
			var sourceName string
			var isOfficial bool
			var typeName string
			
			if err := rows.Scan(
				&s.ID, &s.SourceID, &s.TypeID, &s.Name, &s.ExternalID, &s.Description,
				&sourceName, &isOfficial, &typeName,
			); err != nil {
				return nil, fmt.Errorf("failed to scan skill: %w", err)
			}
			
			s.SourceName = sourceName
			s.TypeName = typeName
			s.IsOfficial = isOfficial
			
			// Calculate confidence based on match quality
			confidence := 0.8
			if strings.EqualFold(s.Name, keyword) {
				confidence = 1.0
			} else if strings.Contains(strings.ToLower(s.Name), strings.ToLower(keyword)) {
				confidence = 0.9
			}
			
			results = append(results, &models.StandardizedSkill{
				Keyword:    keyword,
				Skill:      *s,
				Confidence: confidence,
			})
		}
	}

	return results, nil
}

// intSliceToString converts a slice of ints to a PostgreSQL array string
func intSliceToString(ints []int) string {
	if len(ints) == 0 {
		return ""
	}
	parts := make([]string, len(ints))
	for i, v := range ints {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, ",")
}

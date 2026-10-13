package database

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SimonFremont-dev/cv-builder/internal/models"
	"github.com/SimonFremont-dev/cv-builder/internal/parser"
)

// Service provides high-level database operations for the CV Builder workflow
type Service struct {
	db           *Database
	userRepo    *UserRepository
	skillRepo   *SkillRepository
	jobRepo     *JobRepository
}

// NewService creates a new database service
func NewService(db *Database) *Service {
	return &Service{
		db:        db,
		userRepo:   NewUserRepository(db),
		skillRepo:  NewSkillRepository(db),
		jobRepo:    NewJobRepository(db),
	}
}

// InitializeUserFromResume initializes a user profile from a resume file
func (s *Service) InitializeUserFromResume(resumePath, userName, userEmail string) (*UserDB, error) {
	// Parse the resume
	resumeParser := parser.NewResumeParser()
	resume, err := resumeParser.ParseResume(resumePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resume: %w", err)
	}

	// Create or get user
	user, err := s.userRepo.GetUserByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = s.userRepo.CreateUser(userName, userEmail)
		if err != nil {
			return nil, err
		}
	} else {
		// Update user info
		if err := s.userRepo.UpdateUser(user.ID, userName, userEmail); err != nil {
			return nil, err
		}
	}

	// Import resume data
	if err := s.userRepo.ImportResume(user.ID, resume); err != nil {
		return nil, fmt.Errorf("failed to import resume: %w", err)
	}

	return user, nil
}

// ProcessJobOffer processes a job offer and extracts requirements
func (s *Service) ProcessJobOffer(userID int64, jobOfferPath string) (*JobOfferDB, []*JobRequirementDB, error) {
	// Parse the job offer
	jobParser := parser.NewJobParser()
	jobOffer, err := jobParser.ParseJobOffer(jobOfferPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse job offer: %w", err)
	}

	// Store the job offer
	jobDB, err := s.jobRepo.StoreJobOffer(userID, jobOffer)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to store job offer: %w", err)
	}

	// Store the requirements
	requirementsDB, err := s.jobRepo.StoreJobRequirements(jobDB.ID, jobOffer.Requirements)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to store job requirements: %w", err)
	}

	return jobDB, requirementsDB, nil
}

// AskUserAboutRequirements interactively asks the user about job requirements
func (s *Service) AskUserAboutRequirements(userID int64, jobOfferID int64) ([]*UserResponseDB, error) {
	// Get the job offer with requirements
	job, requirements, err := s.jobRepo.GetJobOfferWithRequirements(jobOfferID)
	if err != nil {
		return nil, err
	}

	// Get user's existing skills
	userSkills, err := s.skillRepo.GetUserSkills(userID)
	if err != nil {
		return nil, err
	}

	// Create a map of user's skills for quick lookup
	userSkillMap := make(map[string]*UserSkillDB)
	for _, skill := range userSkills {
		userSkillMap[strings.ToLower(skill.SkillName)] = skill
	}

	// Get user's existing responses (to check for expired ones)
	expiredResponses, err := s.jobRepo.GetExpiredResponses(userID)
	if err != nil {
		return nil, err
	}

	// Create a map of expired response requirement IDs
	expiredResponseMap := make(map[int64]bool)
	for _, resp := range expiredResponses {
		expiredResponseMap[resp.JobRequirementID] = true
	}

	var responses []*UserResponseDB
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Printf("=== Job Offer: %s - %s ===\n", job.Company, job.Position)
	fmt.Println()
	fmt.Println("Please answer the following questions about the job requirements.")
	fmt.Println("After 1 year, you will be asked again to verify your skills.")
	fmt.Println()

	// Ask about each requirement
	for _, req := range requirements {
		// Skip if user already has this skill and it's not expired
		skillNameLower := strings.ToLower(req.Name)
		userSkill, hasSkill := userSkillMap[skillNameLower]

		// Check if there's an expired response for this requirement
		_, needsReAsk := expiredResponseMap[req.ID]

		if hasSkill && !needsReAsk && !IsSkillExpired(userSkill.LastVerified) {
			// User already has this skill and it's still valid
			fmt.Printf("✓ %s (Already have: %s, %d years)\n", req.Name, userSkill.Proficiency, userSkill.Years)
			continue
		}

		// Ask the user
		fmt.Printf("\n%s: %s\n", req.RequirementType, req.Name)
		if req.MinYears != nil {
			fmt.Printf("  (Required: %d+ years)\n", *req.MinYears)
		}

		var hasSkillInput string
		for {
			fmt.Print("  Do you have this skill/experience? (y/n): ")
			input, _ := reader.ReadString('\n')
			hasSkillInput = strings.TrimSpace(strings.ToLower(input))
			if hasSkillInput == "y" || hasSkillInput == "n" {
				break
			}
			fmt.Println("  Please enter 'y' or 'n'")
		}

		userHasSkill := hasSkillInput == "y"

		var proficiency string
		var years *int
		var notes string

		if userHasSkill {
			// Ask for proficiency
			fmt.Print("  Proficiency (Beginner/Intermediate/Advanced/Expert): ")
			proficiencyInput, _ := reader.ReadString('\n')
			proficiency = strings.TrimSpace(proficiencyInput)

			// Ask for years of experience
			fmt.Print("  Years of experience: ")
			yearsStr, _ := reader.ReadString('\n')
			yearsStr = strings.TrimSpace(yearsStr)
			if yearsStr != "" {
				var y int
				if _, err := fmt.Sscanf(yearsStr, "%d", &y); err == nil {
					years = &y
				}
			}

			// Ask for notes
			fmt.Print("  Additional notes (optional): ")
			notes, _ = reader.ReadString('\n')
			notes = strings.TrimSpace(notes)
		}

		// Store the response
		response, err := s.jobRepo.StoreUserResponse(userID, req.ID, userHasSkill, proficiency, "", notes)
		if err != nil {
			return nil, fmt.Errorf("failed to store user response: %w", err)
		}

		responses = append(responses, response)

		// If user has the skill, also add it to their user_skills
		if userHasSkill {
			// Create a temporary skill model
			skillModel := models.Skill{
				Name:        req.Name,
				Category:    req.Category,
				Proficiency: proficiency,
				Years:       0,
			}
			if years != nil {
				skillModel.Years = *years
			}

			if err := s.skillRepo.UpsertSkillAndUserSkill(userID, skillModel); err != nil {
				// Log but don't fail
				fmt.Printf("  Warning: failed to update user skill: %v\n", err)
			}
		}
	}

	return responses, nil
}

// GenerateUpdatedResume generates an updated resume based on user's profile and job requirements
func (s *Service) GenerateUpdatedResume(userID int64, jobOfferID int64) (*models.Resume, error) {
	// Get the user
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Get the job offer (optional)
	var requirements []*JobRequirementDB
	if jobOfferID > 0 {
		_, requirements, err = s.jobRepo.GetJobOfferWithRequirements(jobOfferID)
		if err != nil {
			return nil, err
		}
	}

	// Get user's current resume
	var resume models.Resume

	// Get user's skills
	userSkills, err := s.skillRepo.GetUserSkills(userID)
	if err != nil {
		return nil, err
	}

	// Get user's work experience
	// (Would need to add a method to get work experience)

	// Get user's education
	// (Would need to add a method to get education)

	// Build the resume
	resume.Metadata = models.ResumeMetadata{
		ID:          fmt.Sprintf("resume-%d", user.ID),
		Title:       fmt.Sprintf("%s's Resume", user.Name),
		Version:     "1.0.0",
		LastUpdated: time.Now(),
		Template:    "default",
	}

	resume.Contact = models.ContactInfo{
		Name:  user.Name,
		Email: user.Email,
	}

	// Add skills from user_skills
	for _, skill := range userSkills {
		resume.Skills = append(resume.Skills, models.Skill{
			Name:        skill.SkillName,
			Category:    skill.Category,
			Proficiency: skill.Proficiency,
			Years:       skill.Years,
		})
	}

	// Add skills from job requirements that user has
	for range requirements {
		// Check if user has this skill (via responses)
		// (Would need to query user_responses)
		// For now, just add all skills from requirements that user has responded positively to
	}

	return &resume, nil
}

// GetMissingRequirements gets requirements that the user hasn't confirmed they have
func (s *Service) GetMissingRequirements(userID int64, jobOfferID int64) ([]*JobRequirementDB, error) {
	// Get all requirements for the job offer
	_, requirements, err := s.jobRepo.GetJobOfferWithRequirements(jobOfferID)
	if err != nil {
		return nil, err
	}

	// Get user's existing responses for this job offer
	userResponses, err := s.jobRepo.GetUserResponses(userID)
	if err != nil {
		return nil, err
	}

	// Create a map of user's positive responses
	userHasMap := make(map[int64]bool)
	for _, resp := range userResponses {
		if resp.HasSkill {
			userHasMap[resp.JobRequirementID] = true
		}
	}

	// Get user's skills that are still valid
	userSkills, err := s.skillRepo.GetUserSkills(userID)
	if err != nil {
		return nil, err
	}

	// Create a map of user's valid skills
	userSkillMap := make(map[string]bool)
	for _, skill := range userSkills {
		if !IsSkillExpired(skill.LastVerified) {
			userSkillMap[strings.ToLower(skill.SkillName)] = true
		}
	}

	// Filter requirements
	var missing []*JobRequirementDB
	for _, req := range requirements {
		// Check if user has this requirement via responses or skills
		hasRequirement := userHasMap[req.ID]
		hasSkill := userSkillMap[strings.ToLower(req.Name)]

		if !hasRequirement && !hasSkill {
			missing = append(missing, req)
		}
	}

	return missing, nil
}

// GetJobOfferWithRequirements retrieves a job offer with all its requirements
func (s *Service) GetJobOfferWithRequirements(jobOfferID int64) (*JobOfferDB, []*JobRequirementDB, error) {
	return s.jobRepo.GetJobOfferWithRequirements(jobOfferID)
}

// GetExpiredSkillsAndResponses returns all expired skills and responses for a user
func (s *Service) GetExpiredSkillsAndResponses(userID int64) ([]*UserSkillDB, []*UserResponseDB, error) {
	expiredSkills, err := s.skillRepo.GetExpiredSkills(userID)
	if err != nil {
		return nil, nil, err
	}

	expiredResponses, err := s.jobRepo.GetExpiredResponses(userID)
	if err != nil {
		return nil, nil, err
	}

	return expiredSkills, expiredResponses, nil
}

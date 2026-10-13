package cli

import (
	"fmt"
	"log"

	"github.com/SimonFremont-dev/cv-builder/internal/database"
	"github.com/SimonFremont-dev/cv-builder/internal/generator"
	"github.com/spf13/cobra"
)

// RegisterDBCommands registers all database-related commands
func RegisterDBCommands(rootCmd *cobra.Command) {
	// Database management commands
	rootCmd.AddCommand(newDBInitCommand())
	rootCmd.AddCommand(newDBSetupCommand())

	// User profile commands
	rootCmd.AddCommand(newDBImportResumeCommand())
	rootCmd.AddCommand(newDBListUsersCommand())

	// Job offer workflow commands
	rootCmd.AddCommand(newDBProcessOfferCommand())
	rootCmd.AddCommand(newDBAskAboutRequirementsCommand())
	rootCmd.AddCommand(newDBListRequirementsCommand())
	rootCmd.AddCommand(newDBListExpiredCommand())
	rootCmd.AddCommand(newDBGenerateResumeCommand())
}

// newDBInitCommand creates the db_init subcommand
func newDBInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db_init",
		Short: "Initialize the SQLite database",
		Long: `Initialize the SQLite database with all required tables.

This command creates the database file and sets up all the tables needed
for storing user profiles, resumes, skills, and job offers.

Examples:
  # Initialize database with default path
  resume-adapter db-init

  # Initialize database with custom path
  resume-adapter db-init --path /path/to/database.db`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			dbPath, _ := cmd.Flags().GetString("path")

			if err := runDBInit(verbose, dbPath); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("path", "p", "", "path to database file (default: data/cv_builder.db)")

	return cmd
}

// runDBInit executes the db-init command
func runDBInit(verbose bool, dbPath string) error {
	if dbPath == "" {
		dbPath = database.DefaultDBPath
	}

	if verbose {
		fmt.Printf("Initializing database at: %s\n", dbPath)
	}

	db, err := database.NewDatabase(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	fmt.Printf("Database initialized successfully at: %s\n", dbPath)

	return nil
}

// newDBSetupCommand creates the db-setup subcommand
func newDBSetupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db-setup",
		Short: "Setup the database and import initial data",
		Long: `Setup the database, import a resume, and optionally process a job offer.

This is a convenience command that combines db-init, import-resume, and optionally
process-offer in one step.

Examples:
  # Setup with resume only
  resume-adapter db-setup --resume input/resume/resume.json --name "John Doe" --email john@example.com

  # Setup with resume and job offer
  resume-adapter db-setup --resume input/resume/resume.json --name "John Doe" --email john@example.com --offer input/job-offers/offer.txt`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			resumePath, _ := cmd.Flags().GetString("resume")
			name, _ := cmd.Flags().GetString("name")
			email, _ := cmd.Flags().GetString("email")
			offerPath, _ := cmd.Flags().GetString("offer")

			if err := runDBSetup(verbose, resumePath, name, email, offerPath); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("resume", "r", "", "path to resume file (required)")
	cmd.Flags().StringP("name", "n", "", "user name")
	cmd.Flags().StringP("email", "e", "", "user email")
	cmd.Flags().StringP("offer", "o", "", "path to job offer file (optional)")
	cmd.MarkFlagRequired("resume")

	return cmd
}

// runDBSetup executes the db-setup command
func runDBSetup(verbose bool, resumePath, name, email, offerPath string) error {
	if name == "" {
		name = "User"
	}
	if email == "" {
		email = "user@example.com"
	}

	// Initialize database
	db, err := database.InitDefaultDB()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Import resume
	service := database.NewService(db)
	user, err := service.InitializeUserFromResume(resumePath, name, email)
	if err != nil {
		return fmt.Errorf("failed to import resume: %w", err)
	}

	if verbose {
		fmt.Printf("User created/updated: %s (%s)\n", user.Name, user.Email)
	}

	// Process job offer if provided
	if offerPath != "" {
		job, requirements, err := service.ProcessJobOffer(user.ID, offerPath)
		if err != nil {
			return fmt.Errorf("failed to process job offer: %w", err)
		}

		if verbose {
			fmt.Printf("Job offer processed: %s - %s\n", job.Company, job.Position)
			fmt.Printf("  Found %d requirements\n", len(requirements))
		}

		// Ask user about requirements
		fmt.Println()
		fmt.Println("Now, let's go through the job requirements...")
		_, err = service.AskUserAboutRequirements(user.ID, job.ID)
		if err != nil {
			return fmt.Errorf("failed to ask user about requirements: %w", err)
		}

		fmt.Println()
		fmt.Println("✓ Setup complete! Your profile has been created and job requirements have been recorded.")
	} else {
		fmt.Println("✓ Setup complete! Your profile has been created from your resume.")
	}

	return nil
}

// newDBImportResumeCommand creates the import-resume subcommand
func newDBImportResumeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-resume",
		Short: "Import a resume into the database",
		Long: `Import a resume file to initialize or update a user's profile in the database.

This command parses a resume file and stores all the information (skills, 
experience, education) in the database for the specified user.

Examples:
  # Import resume for a user
  resume-adapter import-resume --resume input/resume/resume.json --name "John Doe" --email john@example.com

  # Import resume with database path
  resume-adapter import-resume --resume input/resume/resume.json --name "John Doe" --email john@example.com --db /path/to/db`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			resumePath, _ := cmd.Flags().GetString("resume")
			name, _ := cmd.Flags().GetString("name")
			email, _ := cmd.Flags().GetString("email")
			dbPath, _ := cmd.Flags().GetString("db")

			if err := runImportResume(verbose, resumePath, name, email, dbPath); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("resume", "r", "", "path to resume file (required)")
	cmd.Flags().StringP("name", "n", "", "user name")
	cmd.Flags().StringP("email", "e", "", "user email")
	cmd.Flags().StringP("db", "d", "", "path to database file")
	cmd.MarkFlagRequired("resume")

	return cmd
}

// runImportResume executes the import-resume command
func runImportResume(verbose bool, resumePath, name, email, dbPath string) error {
	// Initialize database
	var db *database.Database
	var err error

	if dbPath != "" {
		db, err = database.NewDatabase(dbPath)
	} else {
		db, err = database.InitDefaultDB()
	}
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Import resume
	service := database.NewService(db)
	user, err := service.InitializeUserFromResume(resumePath, name, email)
	if err != nil {
		return fmt.Errorf("failed to import resume: %w", err)
	}

	fmt.Printf("✓ Resume imported successfully for user: %s (%s)\n", user.Name, user.Email)

	return nil
}

// newDBProcessOfferCommand creates the process-offer subcommand
func newDBProcessOfferCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process-offer",
		Short: "Process a job offer and extract requirements",
		Long: `Process a job offer file and extract all requirements (skills, experience, education).

This command parses a job offer and stores it in the database along with all
the requirements extracted from it. The requirements can then be reviewed and
answered by the user.

Examples:
  # Process a job offer for a user
  resume-adapter process-offer --email john@example.com --offer input/job-offers/offer.txt

  # Process with database path
  resume-adapter process-offer --email john@example.com --offer input/job-offers/offer.txt --db /path/to/db`,
		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool("verbose")
			email, _ := cmd.Flags().GetString("email")
			offerPath, _ := cmd.Flags().GetString("offer")
			dbPath, _ := cmd.Flags().GetString("db")

			if err := runProcessOffer(verbose, email, offerPath, dbPath); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("email", "e", "", "user email (required)")
	cmd.Flags().StringP("offer", "o", "", "path to job offer file (required)")
	cmd.Flags().StringP("db", "d", "", "path to database file")
	cmd.MarkFlagRequired("email")
	cmd.MarkFlagRequired("offer")

	return cmd
}

// runProcessOffer executes the process-offer command
func runProcessOffer(verbose bool, email, offerPath, dbPath string) error {
	// Initialize database
	var db *database.Database
	var err error

	if dbPath != "" {
		db, err = database.NewDatabase(dbPath)
	} else {
		db, err = database.InitDefaultDB()
	}
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Get user
	userRepo := database.NewUserRepository(db)
	user, err := userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user with email %s not found", email)
	}

	// Process job offer
	service := database.NewService(db)
	job, requirements, err := service.ProcessJobOffer(user.ID, offerPath)
	if err != nil {
		return fmt.Errorf("failed to process job offer: %w", err)
	}

	fmt.Printf("✓ Job offer processed: %s - %s\n", job.Company, job.Position)
	fmt.Printf("  Found %d requirements\n", len(requirements))

	if verbose {
		fmt.Println("\nRequirements:")
		for _, req := range requirements {
			fmt.Printf("  - [%s] %s", req.RequirementType, req.Name)
			if req.IsRequired {
				fmt.Print(" (Required)")
			}
			if req.MinYears != nil {
				fmt.Printf(" (%d+ years)", *req.MinYears)
			}
			fmt.Println()
		}
	}

	return nil
}

// newDBAskAboutRequirementsCommand creates the ask-requirements subcommand
func newDBAskAboutRequirementsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ask-requirements",
		Short: "Interactively ask user about job requirements",
		Long: `Interactively ask the user whether they have the skills and experience
required by a job offer.

This command will:
1. List all requirements from the job offer
2. Ask the user if they have each skill/experience
3. If yes, ask for proficiency and years of experience
4. Store all answers in the database
5. After 1 year, the user will be asked again to verify

Examples:
  # Ask about requirements for a job offer
  resume-adapter ask-requirements --email john@example.com --job-id 1

  # Ask with database path
  resume-adapter ask-requirements --email john@example.com --job-id 1 --db /path/to/db`,
		Run: func(cmd *cobra.Command, args []string) {
			email, _ := cmd.Flags().GetString("email")
			jobID, _ := cmd.Flags().GetInt64("job-id")
			dbPath, _ := cmd.Flags().GetString("db")

			if err := runAskAboutRequirements(email, jobID, dbPath); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("email", "e", "", "user email (required)")
	cmd.Flags().Int64P("job-id", "j", 0, "job offer ID (required)")
	cmd.Flags().StringP("db", "d", "", "path to database file")
	cmd.MarkFlagRequired("email")
	cmd.MarkFlagRequired("job-id")

	return cmd
}

// runAskAboutRequirements executes the ask-requirements command
func runAskAboutRequirements(email string, jobID int64, dbPath string) error {
	// Initialize database
	var db *database.Database
	var err error

	if dbPath != "" {
		db, err = database.NewDatabase(dbPath)
	} else {
		db, err = database.InitDefaultDB()
	}
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Get user
	userRepo := database.NewUserRepository(db)
	user, err := userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user with email %s not found", email)
	}

	// Ask user about requirements
	service := database.NewService(db)
	responses, err := service.AskUserAboutRequirements(user.ID, jobID)
	if err != nil {
		return fmt.Errorf("failed to ask user about requirements: %w", err)
	}

	fmt.Printf("\n✓ Answered %d requirements\n", len(responses))

	return nil
}

// newDBListRequirementsCommand creates the list-requirements subcommand
func newDBListRequirementsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-requirements",
		Short: "List requirements for a job offer",
		Long: `List all requirements extracted from a job offer.

This command shows all the skills, experience, and education requirements
that were extracted from a job offer and stored in the database.

Examples:
  # List requirements for a job offer
  resume-adapter list-requirements --job-id 1

  # List with user responses
  resume-adapter list-requirements --job-id 1 --email john@example.com`,
		Run: func(cmd *cobra.Command, args []string) {
			jobID, _ := cmd.Flags().GetInt64("job-id")
			email, _ := cmd.Flags().GetString("email")
			dbPath, _ := cmd.Flags().GetString("db")
			missingOnly, _ := cmd.Flags().GetBool("missing")

			if err := runListRequirements(jobID, email, dbPath, missingOnly); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().Int64P("job-id", "j", 0, "job offer ID (required)")
	cmd.Flags().StringP("email", "e", "", "user email (optional, to show user responses)")
	cmd.Flags().StringP("db", "d", "", "path to database file")
	cmd.Flags().BoolP("missing", "m", false, "show only missing requirements")
	cmd.MarkFlagRequired("job-id")

	return cmd
}

// runListRequirements executes the list-requirements command
func runListRequirements(jobID int64, email, dbPath string, missingOnly bool) error {
	// Initialize database
	var db *database.Database
	var err error

	if dbPath != "" {
		db, err = database.NewDatabase(dbPath)
	} else {
		db, err = database.InitDefaultDB()
	}
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	service := database.NewService(db)
	jobRepo := database.NewJobRepository(db)

	// Get job offer info
	job, requirements, err := service.GetJobOfferWithRequirements(jobID)
	if err != nil {
		return err
	}

	fmt.Printf("=== Job Offer: %s - %s ===\n\n", job.Company, job.Position)

	if email != "" {
		// Get user responses
		userRepo := database.NewUserRepository(db)
		user, err := userRepo.GetUserByEmail(email)
		if err != nil {
			return err
		}
		if user == nil {
			return fmt.Errorf("user with email %s not found", email)
		}

		userResponses, err := jobRepo.GetUserResponses(user.ID)
		if err != nil {
			return err
		}

		// Create a map of user responses
		userResponseMap := make(map[int64]*database.UserResponseDB)
		for _, resp := range userResponses {
			userResponseMap[resp.JobRequirementID] = resp
		}

		// List requirements with user responses
		for _, req := range requirements {
			if missingOnly {
				resp, hasResp := userResponseMap[req.ID]
				if hasResp && !resp.HasSkill {
					fmt.Printf("✗ %s: %s\n", req.RequirementType, req.Name)
				}
			} else {
				resp, hasResp := userResponseMap[req.ID]
				status := "?"
				if hasResp {
					if resp.HasSkill {
						status = "✓"
					} else {
						status = "✗"
					}
				}

				fmt.Printf("%s [%s] %s", status, req.RequirementType, req.Name)
				if req.IsRequired {
					fmt.Print(" (Required)")
				}
				if req.MinYears != nil {
					fmt.Printf(" (%d+ years)", *req.MinYears)
				}
				fmt.Println()
			}
		}
	} else {
		// List all requirements
		for _, req := range requirements {
			fmt.Printf("  - [%s] %s", req.RequirementType, req.Name)
			if req.IsRequired {
				fmt.Print(" (Required)")
			}
			if req.MinYears != nil {
				fmt.Printf(" (%d+ years)", *req.MinYears)
			}
			fmt.Println()
		}
	}

	return nil
}

// newDBListExpiredCommand creates the list-expired subcommand
func newDBListExpiredCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-expired",
		Short: "List expired skills and responses",
		Long: `List all skills and job requirement responses that have expired and need re-verification.

After 1 year of not verifying a skill, the user will be asked again whether they
have it. This command shows all expired items that need to be re-verified.

Examples:
  # List expired items for a user
  resume-adapter list-expired --email john@example.com`,
		Run: func(cmd *cobra.Command, args []string) {
			email, _ := cmd.Flags().GetString("email")
			dbPath, _ := cmd.Flags().GetString("db")

			if err := runListExpired(email, dbPath); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("email", "e", "", "user email (required)")
	cmd.Flags().StringP("db", "d", "", "path to database file")
	cmd.MarkFlagRequired("email")

	return cmd
}

// runListExpired executes the list-expired command
func runListExpired(email, dbPath string) error {
	// Initialize database
	var db *database.Database
	var err error

	if dbPath != "" {
		db, err = database.NewDatabase(dbPath)
	} else {
		db, err = database.InitDefaultDB()
	}
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Get user
	userRepo := database.NewUserRepository(db)
	user, err := userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user with email %s not found", email)
	}

	// Get expired items
	service := database.NewService(db)
	expiredSkills, expiredResponses, err := service.GetExpiredSkillsAndResponses(user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("=== Expired Items for %s ===\n\n", user.Name)

	if len(expiredSkills) > 0 {
		fmt.Println("Expired Skills:")
		for _, skill := range expiredSkills {
			fmt.Printf("  - %s (Last verified: %v)\n", skill.SkillName, skill.LastVerified)
		}
		fmt.Println()
	}

	if len(expiredResponses) > 0 {
		fmt.Println("Expired Job Requirement Responses:")
		for _, resp := range expiredResponses {
			fmt.Printf("  - %s: %s (Expired: %v)\n", resp.RequirementType, resp.RequirementName, resp.ExpiresAt)
		}
		fmt.Println()
	}

	if len(expiredSkills) == 0 && len(expiredResponses) == 0 {
		fmt.Println("No expired items found. All skills and responses are still valid.")
	} else {
		fmt.Printf("Total: %d expired skills, %d expired responses\n", len(expiredSkills), len(expiredResponses))
	}

	return nil
}

// newDBGenerateResumeCommand creates the generate-resume subcommand
func newDBGenerateResumeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-resume",
		Short: "Generate a resume from database profile",
		Long: `Generate a resume using the user's profile data stored in the database.

This command creates a resume by combining the user's stored skills, experience,
and education with any additional skills they've confirmed from job offers.

Examples:
  # Generate resume for a user
  resume-adapter generate-resume --email john@example.com --output output/resume.md

  # Generate with specific format
  resume-adapter generate-resume --email john@example.com --output output/resume.pdf --format pdf`,
		Run: func(cmd *cobra.Command, args []string) {
			email, _ := cmd.Flags().GetString("email")
			output, _ := cmd.Flags().GetString("output")
			format, _ := cmd.Flags().GetString("format")
			jobID, _ := cmd.Flags().GetInt64("job-id")
			dbPath, _ := cmd.Flags().GetString("db")

			if err := runGenerateResume(email, output, format, jobID, dbPath); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("email", "e", "", "user email (required)")
	cmd.Flags().StringP("output", "o", "", "output file path (required)")
	cmd.Flags().StringP("format", "f", "markdown", "output format: json, yaml, markdown, html, text, pdf, docx")
	cmd.Flags().Int64P("job-id", "j", 0, "job offer ID (optional, to adapt resume to this job)")
	cmd.Flags().StringP("db", "d", "", "path to database file")
	cmd.MarkFlagRequired("email")
	cmd.MarkFlagRequired("output")

	return cmd
}

// runGenerateResume executes the generate-resume command
func runGenerateResume(email, output, format string, jobID int64, dbPath string) error {
	// Initialize database
	var db *database.Database
	var err error

	if dbPath != "" {
		db, err = database.NewDatabase(dbPath)
	} else {
		db, err = database.InitDefaultDB()
	}
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// Get user
	userRepo := database.NewUserRepository(db)
	user, err := userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user with email %s not found", email)
	}

	// Generate resume
	service := database.NewService(db)
	resume, err := service.GenerateUpdatedResume(user.ID, jobID)
	if err != nil {
		return fmt.Errorf("failed to generate resume: %w", err)
	}

	// Save resume to file
	gen := generator.NewGenerator()
	if err := gen.GenerateResume(resume, format, output); err != nil {
		return fmt.Errorf("failed to generate output: %w", err)
	}

	fmt.Printf("✓ Resume generated: %s\n", output)

	return nil
}

// newDBListUsersCommand creates the list-users subcommand
func newDBListUsersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-users",
		Short: "List all users in the database",
		Long: `List all users registered in the database.

Examples:
  # List all users
  resume-adapter list-users

  # List with database path
  resume-adapter list-users --db /path/to/db`,
		Run: func(cmd *cobra.Command, args []string) {
			dbPath, _ := cmd.Flags().GetString("db")

			if err := runListUsers(dbPath); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("db", "d", "", "path to database file")

	return cmd
}

// runListUsers executes the list-users command
func runListUsers(dbPath string) error {
	// Initialize database
	var db *database.Database
	var err error

	if dbPath != "" {
		db, err = database.NewDatabase(dbPath)
	} else {
		db, err = database.InitDefaultDB()
	}
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	userRepo := database.NewUserRepository(db)
	users, err := userRepo.ListUsers()
	if err != nil {
		return err
	}

	fmt.Println("Users:")
	for _, user := range users {
		fmt.Printf("  - %s (%s) - ID: %d, Created: %v\n", user.Name, user.Email, user.ID, user.CreatedAt)
	}

	fmt.Printf("\nTotal: %d users\n", len(users))

	return nil
}

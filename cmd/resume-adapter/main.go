package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SimonFremont-dev/cv-builder/internal/cli"
)

func main() {
	// Check if any arguments were provided (other than the program name)
	if len(os.Args) < 2 {
		// No command provided, show interactive menu
		showInteractiveMenu()
		return
	}

	// Check if the first argument is help or version
	if len(os.Args) == 2 && (os.Args[1] == "help" || os.Args[1] == "--help" || os.Args[1] == "-h") {
		// Show help
		if err := cli.NewRootCommand().Execute(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return
	}

	// Try to execute the command
	if err := cli.NewRootCommand().Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// showInteractiveMenu displays an interactive menu for selecting commands
func showInteractiveMenu() {
	for {
		fmt.Println()
		fmt.Println("===========================================")
		fmt.Println("  CV Builder - Resume Adapter")
		fmt.Println("===========================================")
		fmt.Println()
		fmt.Println("Database Workflow:")
		fmt.Println("  1.  db-init              - Initialize the SQLite database")
		fmt.Println("  2.  db-setup            - Setup database and import resume")
		fmt.Println("  3.  import-resume       - Import a resume into the database")
		fmt.Println("  4.  process-offer       - Process a job offer")
		fmt.Println("  5.  ask-requirements    - Ask about job requirements")
		fmt.Println("  6.  list-requirements   - List job requirements")
		fmt.Println("  7.  list-expired        - List expired skills")
		fmt.Println("  8.  generate-resume     - Generate resume from DB")
		fmt.Println("  9.  list-users          - List all users")
		fmt.Println()
		fmt.Println("Original Commands:")
		fmt.Println("  10. match               - Match resume vs job offer")
		fmt.Println("  11. adapt               - Create adapted resume")
		fmt.Println("  12. analyze             - Analyze resume vs job")
		fmt.Println("  13. convert             - Convert resume format")
		fmt.Println("  14. generate            - Generate resume template")
		fmt.Println("  15. generate-all        - Generate all formats")
		fmt.Println("  16. adapt-config        - Adapt using config file")
		fmt.Println()
		fmt.Println("  0.  help                - Show full help")
		fmt.Println("  q.  quit/exit           - Exit")
		fmt.Println()
		fmt.Print("Enter command number or name: ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Handle special cases
		if input == "" {
			continue
		}

		if input == "quit" || input == "exit" || input == "q" {
			fmt.Println("Goodbye!")
			return
		}

		if input == "help" || input == "0" {
			rootCmd := cli.NewRootCommand()
			rootCmd.SetArgs([]string{"--help"})
			_ = rootCmd.Execute()
			continue
		}

		// Map of command numbers to command names
		commandMap := map[string]string{
			"1":  "db-init",
			"2":  "db-setup",
			"3":  "import-resume",
			"4":  "process-offer",
			"5":  "ask-requirements",
			"6":  "list-requirements",
			"7":  "list-expired",
			"8":  "generate-resume",
			"9":  "list-users",
			"10": "match",
			"11": "adapt",
			"12": "analyze",
			"13": "convert",
			"14": "generate",
			"15": "generate-all",
			"16": "adapt-config",
		}

		var cmdName string

		// Check if input is a number
		if num, ok := commandMap[input]; ok {
			cmdName = num
		} else {
			cmdName = input
		}

		// Run the command
		rootCmd := cli.NewRootCommand()
		rootCmd.SetArgs([]string{cmdName})
		if err := rootCmd.Execute(); err != nil {
			log.Println("Error:", err)
		}
	}
}

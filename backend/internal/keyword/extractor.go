package keyword

import (
	"regexp"
	"strings"
)

// Extractor handles keyword extraction from text
type Extractor struct {
	categories map[string][]string
	minLength  int
	maxLength  int
}

// NewExtractor creates a new keyword extractor
func NewExtractor(categories map[string][]string) *Extractor {
	return &Extractor{
		categories: categories,
		minLength:  3,
		maxLength:  30,
	}
}

// ExtractKeywords extracts keywords from text
func (e *Extractor) ExtractKeywords(text string) []string {
	// Normalize text
	text = strings.ToLower(text)
	
	// Remove special characters but keep alphanumeric and spaces
	regex := regexp.MustCompile(`[^a-zA-Z0-9\s-]`)
	text = regex.ReplaceAllString(text, " ")
	
	// Split into words
	words := strings.Fields(text)
	
	// Filter by length and collect unique
	seen := make(map[string]bool)
	var keywords []string
	
	for _, word := range words {
		word = strings.TrimSpace(word)
		// Skip if too short or too long
		if len(word) < e.minLength || len(word) > e.maxLength {
			continue
		}
		
		// Skip common words
		if isCommonWord(word) {
			continue
		}
		
		if !seen[word] {
			seen[word] = true
			keywords = append(keywords, word)
		}
	}
	
	return keywords
}

// ExtractTechnicalKeywords extracts technical keywords (skills, tools, etc.)
func (e *Extractor) ExtractTechnicalKeywords(text string) []string {
	keywords := e.ExtractKeywords(text)
	
	// Known technical keywords database
	technicalKeywords := getTechnicalKeywords()
	
	var technical []string
	for _, kw := range keywords {
		if containsString(technicalKeywords, kw) {
			technical = append(technical, kw)
		}
	}
	
	return technical
}

// containsString checks if a string exists in a slice
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// isCommonWord checks if a word is a common non-technical word
func isCommonWord(word string) bool {
	commonWords := []string{
		"the", "and", "a", "an", "in", "on", "at", "to", "for", "of", "with",
		"is", "are", "was", "were", "be", "been", "being", "have", "has", "had",
		"do", "does", "did", "will", "would", "should", "could", "may", "might",
		"as", "by", "from", "this", "that", "these", "those", "it", "its", "they",
		"we", "you", "he", "she", "them", "their", "our", "your", "his", "her",
		"what", "which", "who", "when", "where", "why", "how", "all", "any",
		"both", "each", "few", "more", "most", "other", "some", "such", "no",
		"nor", "not", "only", "own", "same", "so", "than", "too", "very",
	}
	
	return containsString(commonWords, strings.ToLower(word))
}

// getTechnicalKeywords returns a list of known technical keywords
func getTechnicalKeywords() []string {
	return []string{
		// Programming languages
		"go", "golang", "python", "java", "javascript", "typescript", "js", "ts",
		"c", "c++", "cpp", "csharp", "c#", "ruby", "php", "rust", "kotlin",
		"swift", "scala", "r", "perl", "bash", "shell",
		
		// Tools
		"docker", "kubernetes", "k8s", "helm", "terraform", "ansible", "puppet",
		"chef", "git", "svn", "mercurial", "github", "gitlab", "bitbucket",
		"jenkins", "travis", "circleci", "github actions", "gitlab ci",
		"argocd", "flux", "tekton",
		
		// Cloud
		"aws", "azure", "gcp", "google cloud", "amazon web services",
		"ec2", "s3", "lambda", "eks", "ecs", "rds", "dynamodb",
		"vpc", "iam", "cloudformation", "cloudwatch",
		
		// Databases
		"postgres", "postgresql", "mysql", "mariadb", "mongodb", "redis",
		"sqlite", "cassandra", "elasticsearch", "neo4j",
		
		// Frameworks
		"react", "angular", "vue", "svelte", "nextjs", "nuxt", "express",
		"django", "flask", "spring", "spring boot", "laravel", "ruby on rails",
		"gin", "echo", "fiber", "gorilla",
		
		// DevOps
		"ci", "cd", "cicd", "pipeline", "deployment", "monitoring", "logging",
		"prometheus", "grafana", "loki", "jaeger", "zipkin",
		"nagios", "datadog", "new relic",
		
		// Network
		"tcp", "udp", "http", "https", "dns", "dhcp", "ssl", "tls",
		"vpn", "firewall", "load balancer", "proxy", "cdn",
		"rest", "graphql", "grpc", "websocket", "api",
		
		// Security
		"authentication", "authorization", "jwt", "oauth", "saml", "ldap",
		"encryption", "hashing", "ssl", "tls", "certificate",
		
		// Methodologies
		"agile", "scrum", "kanban", "devops", "sre", "tdr", "pair programming",
		"code review", "testing", "unit test", "integration test",
		
		// Skills
		"orchestration", "containerization", "microservices", "monolith",
		"scalability", "performance", "security", "reliability",
		"troubleshooting", "debugging", "optimization", "refactoring",
		"reviewing", "mentoring", "documentation",
		
		// Operating Systems
		"linux", "unix", "windows", "macos", "ubuntu", "centos", "debian",
		"fedora", "rhel", "arch",
	}
}

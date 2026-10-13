package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"strconv"
	"time"

	"github.com/SimonFremont-dev/cv-builder/backend/internal/models"
	"github.com/SimonFremont-dev/cv-builder/backend/internal/service"
	"github.com/gorilla/mux"
)

// Handlers holds API handlers
type Handlers struct {
	service *service.Service
}

// NewHandlers creates new API handlers
func NewHandlers(service *service.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

// SetupRoutes sets up API routes
func (h *Handlers) SetupRoutes(r *mux.Router) {
	// Profile routes
	r.HandleFunc("/api/profiles", h.listProfiles).Methods("GET")
	r.HandleFunc("/api/profiles", h.createProfile).Methods("POST")
	r.HandleFunc("/api/profiles/{id}", h.getProfile).Methods("GET")
	r.HandleFunc("/api/profiles/{id}", h.updateProfile).Methods("PUT")

	// Skill routes
	r.HandleFunc("/api/profiles/{id}/skills", h.listSkills).Methods("GET")
	r.HandleFunc("/api/profiles/{id}/skills", h.addSkill).Methods("POST")

	// Keyword extraction routes
	r.HandleFunc("/api/extract-keywords", h.extractKeywords).Methods("POST")

	// Matching routes
	r.HandleFunc("/api/match", h.matchProfile).Methods("POST")

	// Resume generation routes
	r.HandleFunc("/api/resume/generate", h.generateResume).Methods("POST")

	// Job offer processing
	r.HandleFunc("/api/job-offers/process", h.processJobOffer).Methods("POST")

	// Skills endpoint
	r.HandleFunc("/api/skills", h.getAllSkills).Methods("GET")

	// URL fetch endpoint
	r.HandleFunc("/api/fetch-url", h.fetchURL).Methods("GET")

	// Health check
	r.HandleFunc("/api/health", h.healthCheck).Methods("GET")

	// Serve static files (frontend)
	// Use FileServer directly for static directory
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.PathPrefix("/static/").Handler(staticHandler)

	// Root path serves index.html
	r.HandleFunc("/", h.serveIndex).Methods("GET")
}

// listProfiles lists all profiles
func (h *Handlers) listProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.service.GetAllProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profiles)
}

// createProfile creates a new profile
func (h *Handlers) createProfile(w http.ResponseWriter, r *http.Request) {
	var profile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create profile using SQLite repository
	err := h.service.CreateProfile(&profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// getProfile gets a profile by ID
func (h *Handlers) getProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	profile, err := h.service.MatchProfileToKeywords(id, []string{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// updateProfile updates a profile
func (h *Handlers) updateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	var profile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile.ID = id

	// Update profile using SQLite repository
	if err := h.service.UpdateProfile(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// listSkills lists skills for a profile
func (h *Handlers) listSkills(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	skills, err := h.service.GetProfileSkills(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skills)
}

// addSkill adds a skill to a profile
func (h *Handlers) addSkill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	var skillUpdate models.SkillUpdate
	if err := json.NewDecoder(r.Body).Decode(&skillUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add skill
	if err := h.service.AddMissingSkill(id, skillUpdate.Name, &skillUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated skills list
	skills, err := h.service.GetProfileSkills(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skills)
}

// extractKeywords extracts keywords from text
func (h *Handlers) extractKeywords(w http.ResponseWriter, r *http.Request) {
	var req models.KeywordExtractionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	keywords, err := h.service.ExtractKeywords(req.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.KeywordExtractionResponse{
		Keywords: keywords,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// matchProfile matches a profile against keywords
func (h *Handlers) matchProfile(w http.ResponseWriter, r *http.Request) {
	var req models.MatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matchResult, err := h.service.MatchProfileToKeywords(req.ProfileID, req.Keywords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matchResult)
}

// generateResume generates a resume
func (h *Handlers) generateResume(w http.ResponseWriter, r *http.Request) {
	var req models.ResumeGenerationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get keywords if provided, otherwise use empty slice
	var keywords []string
	if req.Keywords != nil {
		keywords = req.Keywords
	} else if req.JobOfferID > 0 {
		// For now, we don't have SQL repo in service, so we'll skip this
		// This will be enhanced when we integrate SQL repository
	}

	resumeContent, err := h.service.GenerateResume(req.ProfileID, keywords, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save resume
	resume := &models.Resume{
		ProfileID:  req.ProfileID,
		JobOfferID: req.JobOfferID,
		Content:    resumeContent,
		Format:     "text",
	}
	if err := h.service.SaveResume(resume); err != nil {
		log.Printf("Warning: failed to save resume: %v", err)
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, resumeContent)
}

// processJobOffer processes a job offer and generates a resume
func (h *Handlers) processJobOffer(w http.ResponseWriter, r *http.Request) {
	var jobOffer models.JobOffer
	if err := json.NewDecoder(r.Body).Decode(&jobOffer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get profile ID from query parameter
	profileID, err := strconv.Atoi(r.URL.Query().Get("profile_id"))
	if err != nil || profileID == 0 {
		http.Error(w, "Missing or invalid profile_id query parameter", http.StatusBadRequest)
		return
	}

	matchResult, resumeContent, err := h.service.ProcessJobOffer(profileID, &jobOffer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return both match result and resume content
	response := map[string]interface{}{
		"match":   matchResult,
		"resume":  resumeContent,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// healthCheck returns health status
func (h *Handlers) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
		"message": "Backend is running",
	})
}

// serveIndex serves the main index.html file
func (h *Handlers) serveIndex(w http.ResponseWriter, r *http.Request) {
	// Try to serve index.html from static directory
	http.ServeFile(w, r, "static/index.html")
}

// fetchURL fetches content from a URL
func (h *Handlers) fetchURL(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	// Check if it's a LinkedIn URL - they block scraping
	if strings.Contains(url, "linkedin.com/jobs") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "LinkedIn URLs cannot be fetched automatically due to their anti-scraping protection",
			"hint":    "Please copy the job description text directly and paste it into the text area",
			"url":     url,
		})
		return
	}

	// Create HTTP client with timeout and redirect following
	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 5 redirects
			if len(via) >= 5 {
				return fmt.Errorf("too many redirects")
			}
			// Check if redirecting to a blocked domain
			if strings.Contains(req.URL.Host, "linkedin.com") && len(via) > 0 {
				return fmt.Errorf("linkedin redirect detected")
			}
			return nil
		},
	}

	// Create request with headers to mimic a browser
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; CV-Builder/1.0; +http://localhost)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	// Fetch the URL
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch URL: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("HTTP %d: %s - %s", resp.StatusCode, resp.Status, string(body)), resp.StatusCode)
		return
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the content
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"content":   string(body),
		"url":      url,
		"status":   resp.Status,
		"length":   len(body),
	})
}

// getAllSkills returns all skills from the database
func (h *Handlers) getAllSkills(w http.ResponseWriter, r *http.Request) {
	skills, err := h.service.GetAllSkills()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"skills": skills,
		"count":  len(skills),
	})
}

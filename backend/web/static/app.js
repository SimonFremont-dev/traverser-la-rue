// CV Builder Frontend - TypeScript
// Enhanced with better loading states and error handling
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
// Global state
let currentJobOfferText = '';
let extractedKeywords = [];
let matchedSkills = [];
let unmatchedKeywords = [];
let currentProfileId = null;
let dbSkills = [];
// DOM Elements
const loadSection = document.getElementById('load-section');
const keywordsSection = document.getElementById('keywords-section');
const skillsSection = document.getElementById('skills-section');
const resumeSection = document.getElementById('resume-section');
const loadingOverlay = document.getElementById('loading-overlay');
const loadingMessage = document.getElementById('loading-message');
const errorMessage = document.getElementById('error-message');
const createProfileModal = document.getElementById('create-profile-modal');
const newProfileNameInput = document.getElementById('new-profile-name');
const newProfileEmailInput = document.getElementById('new-profile-email');
const offerUrlInput = document.getElementById('offer-url');
const offerFileInput = document.getElementById('offer-file');
const offerTextInput = document.getElementById('offer-text');
const keywordsList = document.getElementById('keywords-list');
const matchedSkillsContainer = document.getElementById('matched-skills');
const unmatchedKeywordsContainer = document.getElementById('unmatched-keywords');
const resumeContent = document.getElementById('resume-content');
const profileSelect = document.getElementById('profile-select');
// API Configuration
const API_BASE = '/api';
// Show loading state with optional message
function showLoading(show = true, message = 'Processing...') {
    if (show) {
        loadingOverlay.classList.remove('hidden');
        if (loadingMessage) {
            loadingMessage.textContent = message;
        }
    }
    else {
        loadingOverlay.classList.add('hidden');
    }
}
// Modal functions
function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.remove('hidden');
    }
}
function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.add('hidden');
        // Clear inputs when closing create profile modal
        if (modalId === 'create-profile-modal') {
            newProfileNameInput.value = '';
            newProfileEmailInput.value = '';
        }
    }
}
// Show error message
function showError(message, duration = 5000) {
    console.log('Error:', message);
    errorMessage.textContent = message;
    errorMessage.classList.remove('hidden');
    setTimeout(() => {
        errorMessage.classList.add('hidden');
    }, duration);
}
// Show success/info message
function showInfo(message, duration = 3000) {
    errorMessage.textContent = message;
    errorMessage.classList.remove('hidden');
    errorMessage.style.backgroundColor = '#2ecc71';
    setTimeout(() => {
        errorMessage.classList.add('hidden');
        errorMessage.style.backgroundColor = '';
    }, duration);
}
// Hide all sections
function hideAllSections() {
    loadSection.classList.add('hidden');
    keywordsSection.classList.add('hidden');
    skillsSection.classList.add('hidden');
    resumeSection.classList.add('hidden');
}
// Load profiles from API
function loadProfiles() {
    return __awaiter(this, void 0, void 0, function* () {
        showLoading(true, 'Loading profiles...');
        try {
            const response = yield fetch(`${API_BASE}/profiles`);
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            const profiles = yield response.json();
            // Clear existing options
            profileSelect.innerHTML = '';
            if (profiles.length === 0) {
                const option = document.createElement('option');
                option.value = '';
                option.textContent = 'No profiles found - Click "Create New Profile" above';
                profileSelect.appendChild(option);
            }
            else {
                profiles.forEach((profile) => {
                    const option = document.createElement('option');
                    option.value = profile.id;
                    option.textContent = `${profile.name} (${profile.email})`;
                    profileSelect.appendChild(option);
                });
            }
            showLoading(false);
        }
        catch (error) {
            console.error('Error loading profiles:', error);
            showLoading(false);
            showError(`Failed to load profiles: ${error}`);
        }
    });
}
// Create a new profile
function createProfile() {
    return __awaiter(this, void 0, void 0, function* () {
        const name = newProfileNameInput.value.trim();
        const email = newProfileEmailInput.value.trim();
        if (!name) {
            showError('Please enter your name');
            return;
        }
        if (!email) {
            showError('Please enter your email');
            return;
        }
        if (!email.includes('@')) {
            showError('Please enter a valid email address');
            return;
        }
        showLoading(true, 'Creating profile...');
        closeModal('create-profile-modal');
        try {
            const response = yield fetch(`${API_BASE}/profiles`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: name,
                    email: email
                })
            });
            if (!response.ok) {
                const errorText = yield response.text();
                throw new Error(errorText || `HTTP ${response.status}: ${response.statusText}`);
            }
            const profile = yield response.json();
            showLoading(false);
            showInfo(`Profile "${name}" created successfully!`);
            // Reload profiles to update the dropdown
            yield loadProfiles();
            // If this is the first profile, select it automatically
            if (profile.id) {
                currentProfileId = profile.id;
                // Select this profile in the dropdown
                for (let i = 0; i < profileSelect.options.length; i++) {
                    if (profileSelect.options[i].value === String(profile.id)) {
                        profileSelect.selectedIndex = i;
                        break;
                    }
                }
            }
        }
        catch (error) {
            console.error('Error creating profile:', error);
            showLoading(false);
            const errorMsg = error instanceof Error ? error.message : String(error);
            showError(`Failed to create profile: ${errorMsg}`);
        }
    });
}
// Load skills from database
function loadDatabaseSkills() {
    return __awaiter(this, void 0, void 0, function* () {
        showLoading(true, 'Loading skills database...');
        try {
            const response = yield fetch(`${API_BASE}/skills`);
            if (!response.ok) {
                // If endpoint fails, return empty array
                if (response.status === 404 || response.status === 500) {
                    console.log('Skills endpoint not available');
                    showLoading(false);
                    return [];
                }
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            const data = yield response.json();
            dbSkills = data.skills || [];
            console.log('Loaded', dbSkills.length, 'skills from database');
            showLoading(false);
            return dbSkills;
        }
        catch (error) {
            console.error('Error loading database skills:', error);
            showLoading(false);
            showError(`Failed to load skills: ${error}`);
            return [];
        }
    });
}
// Match keywords with database skills
function matchKeywordsWithSkills(keywords) {
    const matched = [];
    const unmatched = [];
    if (dbSkills.length === 0) {
        return { matched: [], unmatched: keywords };
    }
    const skillNames = dbSkills.map((skill) => skill.name.toLowerCase());
    for (const keyword of keywords) {
        const keywordLower = keyword.toLowerCase();
        if (skillNames.some(name => name.includes(keywordLower) || keywordLower.includes(name))) {
            matched.push(keyword);
        }
        else {
            unmatched.push(keyword);
        }
    }
    return { matched, unmatched };
}
// Display database skills
function displayDatabaseSkills() {
    // Remove existing database skills display if any
    const existing = document.getElementById('database-skills');
    if (existing) {
        existing.remove();
    }
    const container = document.createElement('div');
    container.className = 'skills-grid';
    container.id = 'database-skills';
    const title = document.createElement('h3');
    title.textContent = `Skills from Database (${dbSkills.length} total):`;
    title.style.gridColumn = '1 / -1';
    title.style.marginBottom = '15px';
    container.appendChild(title);
    if (dbSkills.length === 0) {
        const msg = document.createElement('p');
        msg.style.color = '#7f8c8d';
        msg.textContent = 'No skills in database. Configure PostgreSQL with ESCO/O*NET data via migrations.';
        container.appendChild(msg);
    }
    else {
        dbSkills.forEach((skill) => {
            var _a, _b;
            const card = document.createElement('div');
            card.className = 'skill-card';
            card.innerHTML = `
                <h3><span class="match-indicator matched"></span>${skill.name}</h3>
                <div class="skill-info">
                    <span>Source: ${((_a = skill.source) === null || _a === void 0 ? void 0 : _a.name) || 'Unknown'}</span>
                    <span>Type: ${((_b = skill.type) === null || _b === void 0 ? void 0 : _b.name) || 'Unknown'}</span>
                </div>
            `;
            container.appendChild(card);
        });
    }
    // Insert before matched skills
    const matchedContainer = document.getElementById('matched-skills');
    if (matchedContainer && matchedContainer.parentNode) {
        matchedContainer.parentNode.insertBefore(container, matchedContainer);
    }
}
// Load from URL
function loadFromURL() {
    return __awaiter(this, void 0, void 0, function* () {
        const url = offerUrlInput.value.trim();
        if (!url) {
            showError('Please enter a URL');
            return;
        }
        // Check if it's a LinkedIn URL
        if (url.includes('linkedin.com/jobs')) {
            showError('LinkedIn URLs are protected by anti-scraping. Please copy the job description text and paste it directly below.', 8000);
            return;
        }
        showLoading(true, 'Fetching URL content...');
        try {
            console.log('Fetching URL:', url);
            const response = yield fetch(`/api/fetch-url?url=${encodeURIComponent(url)}`);
            if (!response.ok) {
                const errorText = yield response.text();
                console.error('Fetch URL error:', errorText);
                throw new Error(errorText || `HTTP ${response.status}: ${response.statusText}`);
            }
            const data = yield response.json();
            console.log('URL fetch response:', data);
            if (data.content && data.content.trim()) {
                offerTextInput.value = data.content;
                showLoading(false);
                showInfo('URL loaded! Click "Process Job Offer" to continue.');
            }
            else if (data.error) {
                throw new Error(data.error);
            }
            else {
                throw new Error('No content received from URL');
            }
        }
        catch (error) {
            console.error('Error loading from URL:', error);
            showLoading(false);
            const errorMsg = error instanceof Error ? error.message : String(error);
            showError(`Failed to load URL: ${errorMsg}. Please try pasting the text directly.`);
        }
    });
}
// Load from file
function loadFromFile() {
    var _a;
    const file = (_a = offerFileInput.files) === null || _a === void 0 ? void 0 : _a[0];
    if (!file) {
        showError('Please select a file');
        return;
    }
    showLoading(true, 'Reading file...');
    const reader = new FileReader();
    reader.onload = (e) => {
        var _a;
        showLoading(false);
        offerTextInput.value = (_a = e.target) === null || _a === void 0 ? void 0 : _a.result;
        showInfo('File loaded! Click "Process Job Offer" to continue.');
    };
    reader.onerror = () => {
        showLoading(false);
        showError('Failed to read file');
    };
    reader.readAsText(file);
}
// Process text
function processText() {
    return __awaiter(this, void 0, void 0, function* () {
        const text = offerTextInput.value.trim();
        if (!text) {
            showError('Please provide job offer text');
            return;
        }
        currentJobOfferText = text;
        showLoading(true, 'Extracting keywords from job offer...');
        console.log('Processing text of length:', text.length);
        try {
            // Extract keywords from the text
            const response = yield fetch(`${API_BASE}/extract-keywords`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ text })
            });
            console.log('Extract keywords response status:', response.status);
            if (!response.ok) {
                const errorText = yield response.text();
                console.error('Extract keywords error:', errorText);
                throw new Error(errorText || `HTTP ${response.status}: ${response.statusText}`);
            }
            const data = yield response.json();
            console.log('Extracted keywords:', data);
            extractedKeywords = data.keywords || [];
            if (extractedKeywords.length === 0) {
                showLoading(false);
                showError('No keywords were extracted from the text. Please try a different job description.');
                return;
            }
            // Load database skills
            showLoading(true, 'Loading skills database...');
            yield loadDatabaseSkills();
            // Display keywords
            displayKeywords();
            // Move to keywords section
            hideAllSections();
            keywordsSection.classList.remove('hidden');
            showLoading(false);
            showInfo(`${extractedKeywords.length} keywords extracted!`);
        }
        catch (error) {
            console.error('Error processing text:', error);
            showLoading(false);
            const errorMsg = error instanceof Error ? error.message : String(error);
            showError(`Failed to process job offer: ${errorMsg}`);
        }
    });
}
// Display extracted keywords
function displayKeywords() {
    keywordsList.innerHTML = '';
    if (extractedKeywords.length === 0) {
        keywordsList.innerHTML = '<p style="color: #7f8c8d;">No keywords extracted</p>';
        return;
    }
    extractedKeywords.forEach(keyword => {
        const tag = document.createElement('span');
        tag.className = 'tag';
        tag.textContent = keyword;
        keywordsList.appendChild(tag);
    });
}
// Move to skills matching
function nextToSkills() {
    return __awaiter(this, void 0, void 0, function* () {
        showLoading(true, 'Loading profiles and skills...');
        try {
            // Load profiles
            yield loadProfiles();
            // Move to skills section
            hideAllSections();
            skillsSection.classList.remove('hidden');
            // Display database skills
            displayDatabaseSkills();
            // Match keywords against skills
            if (extractedKeywords.length > 0) {
                const matched = matchKeywordsWithSkills(extractedKeywords);
                console.log('Matched keywords:', matched.matched);
                console.log('Unmatched keywords:', matched.unmatched);
                displayMatchedKeywords(matched.matched, matched.unmatched);
            }
            showLoading(false);
        }
        catch (error) {
            console.error('Error in nextToSkills:', error);
            showLoading(false);
            const errorMsg = error instanceof Error ? error.message : String(error);
            showError(`Failed to load skills: ${errorMsg}`);
        }
    });
}
// Display matched and unmatched keywords
function displayMatchedKeywords(matched, unmatched) {
    // Remove existing matched keywords display if any
    const existingMatched = document.getElementById('matched-keywords');
    if (existingMatched) {
        existingMatched.remove();
    }
    // Display matched keywords
    if (matched.length > 0) {
        const matchedContainer = document.createElement('div');
        matchedContainer.className = 'skills-grid';
        matchedContainer.id = 'matched-keywords';
        const title = document.createElement('h3');
        title.textContent = `Matched Keywords (${matched.length}):`;
        title.style.color = '#2ecc71';
        title.style.marginBottom = '15px';
        title.style.gridColumn = '1 / -1';
        matchedContainer.appendChild(title);
        matched.forEach(keyword => {
            const card = document.createElement('div');
            card.className = 'skill-card';
            card.innerHTML = `
                <h3><span class="match-indicator matched"></span>${keyword}</h3>
                <p style="color: #2ecc71; font-size: 14px;">✓ Found in database</p>
            `;
            matchedContainer.appendChild(card);
        });
        // Insert before unmatched
        const unmatchedContainer = document.getElementById('unmatched-keywords');
        if (unmatchedContainer && unmatchedContainer.parentNode) {
            unmatchedContainer.parentNode.insertBefore(matchedContainer, unmatchedContainer);
        }
        else {
            skillsSection.insertBefore(matchedContainer, unmatchedKeywordsContainer);
        }
    }
    // Display unmatched keywords
    unmatchedKeywordsContainer.innerHTML = '';
    if (unmatched.length > 0) {
        unmatchedKeywordsContainer.classList.remove('hidden');
        const title = document.createElement('h3');
        title.textContent = `Unmatched Keywords (${unmatched.length}) - Click to add:`;
        title.style.color = '#e74c3c';
        title.style.marginBottom = '15px';
        title.style.gridColumn = '1 / -1';
        unmatchedKeywordsContainer.appendChild(title);
        unmatched.forEach(keyword => {
            const card = document.createElement('div');
            card.className = 'skill-card';
            const safeKeyword = keyword.replace(/'/g, "\\'");
            card.innerHTML = `
                <h3><span class="match-indicator unmatched"></span>${keyword}</h3>
                <button onclick="addAsCustomSkill('${safeKeyword}')" class="btn btn-success" style="padding: 6px 12px; font-size: 14px;">
                    Add to Profile
                </button>
            `;
            unmatchedKeywordsContainer.appendChild(card);
        });
    }
    else {
        unmatchedKeywordsContainer.classList.add('hidden');
    }
}
// Load profile skills
function loadProfileSkills() {
    return __awaiter(this, void 0, void 0, function* () {
        const profileId = profileSelect.value;
        if (!profileId) {
            matchedSkillsContainer.innerHTML = '<p style="color: #7f8c8d;">Please select a profile</p>';
            return;
        }
        currentProfileId = parseInt(profileId);
        showLoading(true, 'Loading profile skills...');
        try {
            // Get profile skills
            const response = yield fetch(`${API_BASE}/profiles/${profileId}/skills`);
            if (!response.ok) {
                const errorText = yield response.text();
                throw new Error(errorText || `HTTP ${response.status}: ${response.statusText}`);
            }
            const skills = yield response.json();
            matchedSkills = skills || [];
            // Display matched skills
            displayMatchedSkills();
            // Match keywords against skills
            if (extractedKeywords.length > 0) {
                const matched = matchKeywordsWithSkills(extractedKeywords);
                displayMatchedKeywords(matched.matched, matched.unmatched);
            }
            showLoading(false);
        }
        catch (error) {
            console.error('Error loading profile skills:', error);
            showLoading(false);
            const errorMsg = error instanceof Error ? error.message : String(error);
            showError(`Failed to load profile skills: ${errorMsg}`);
        }
    });
}
// Display matched skills
function displayMatchedSkills() {
    matchedSkillsContainer.innerHTML = '';
    if (matchedSkills.length === 0) {
        matchedSkillsContainer.innerHTML = '<p style="color: #7f8c8d;">No custom skills found for this profile</p>';
        return;
    }
    matchedSkills.forEach(skill => {
        const card = document.createElement('div');
        card.className = 'skill-card';
        card.innerHTML = `
            <h3><span class="match-indicator matched"></span>${skill.skill_name || skill.name || 'Unnamed Skill'}</h3>
            <div class="skill-info">
                <span>Proficiency: ${skill.proficiency || 'Not specified'}</span>
                <span>Experience: ${skill.years_experience || 0} years</span>
            </div>
        `;
        matchedSkillsContainer.appendChild(card);
    });
}
// Add keyword as custom skill
function addAsCustomSkill(keyword) {
    return __awaiter(this, void 0, void 0, function* () {
        if (!currentProfileId) {
            showError('Please select a profile first');
            return;
        }
        showLoading(true, `Adding "${keyword}" to profile...`);
        try {
            const response = yield fetch(`${API_BASE}/profiles/${currentProfileId}/skills`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: keyword,
                    category: 'Other',
                    proficiency: 'Intermediate',
                    years: 1
                })
            });
            if (!response.ok) {
                const errorText = yield response.text();
                throw new Error(errorText || `HTTP ${response.status}: ${response.statusText}`);
            }
            const skills = yield response.json();
            showLoading(false);
            showInfo(`"${keyword}" added to your profile!`);
            // Reload profile skills
            yield loadProfileSkills();
        }
        catch (error) {
            console.error('Error adding skill:', error);
            showLoading(false);
            const errorMsg = error instanceof Error ? error.message : String(error);
            showError(`Failed to add skill: ${errorMsg}`);
        }
    });
}
// Generate resume
function generateResume() {
    return __awaiter(this, void 0, void 0, function* () {
        if (!currentProfileId) {
            showError('Please select a profile first');
            return;
        }
        showLoading(true, 'Generating your resume...');
        try {
            const response = yield fetch(`${API_BASE}/resume/generate`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    profile_id: currentProfileId,
                    keywords: extractedKeywords,
                    include_unmatched: true
                })
            });
            if (!response.ok) {
                const errorText = yield response.text();
                throw new Error(errorText || `HTTP ${response.status}: ${response.statusText}`);
            }
            const resumeText = yield response.text();
            // Display resume
            resumeContent.textContent = resumeText;
            // Move to resume section
            hideAllSections();
            resumeSection.classList.remove('hidden');
            showLoading(false);
            showInfo('Resume generated! You can download or copy it.');
        }
        catch (error) {
            console.error('Error generating resume:', error);
            showLoading(false);
            const errorMsg = error instanceof Error ? error.message : String(error);
            showError(`Failed to generate resume: ${errorMsg}`);
        }
    });
}
// Download resume
function downloadResume() {
    const content = resumeContent.textContent;
    if (!content) {
        showError('No resume to download');
        return;
    }
    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `resume-${currentProfileId || 'unknown'}-${Date.now()}.txt`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    showInfo('Resume downloaded!');
}
// Copy to clipboard
function copyToClipboard() {
    return __awaiter(this, void 0, void 0, function* () {
        const content = resumeContent.textContent;
        if (!content) {
            showError('No resume to copy');
            return;
        }
        try {
            yield navigator.clipboard.writeText(content);
            showInfo('Resume copied to clipboard!');
        }
        catch (error) {
            console.error('Failed to copy:', error);
            showError('Failed to copy to clipboard');
        }
    });
}
// Start over
function startOver() {
    currentJobOfferText = '';
    extractedKeywords = [];
    matchedSkills = [];
    unmatchedKeywords = [];
    currentProfileId = null;
    dbSkills = [];
    offerUrlInput.value = '';
    offerFileInput.value = '';
    offerTextInput.value = '';
    // Clear all dynamic content
    const databaseSkills = document.getElementById('database-skills');
    if (databaseSkills)
        databaseSkills.remove();
    const matchedKeywords = document.getElementById('matched-keywords');
    if (matchedKeywords)
        matchedKeywords.remove();
    hideAllSections();
    loadSection.classList.remove('hidden');
    showInfo('Ready for a new job offer!');
}
// Initialize
function init() {
    return __awaiter(this, void 0, void 0, function* () {
        showLoading(true, 'Initializing...');
        try {
            // Load profiles on page load
            yield loadProfiles();
            // Try to load database skills
            yield loadDatabaseSkills();
            showLoading(false);
            console.log('Initialization complete');
        }
        catch (error) {
            console.error('Error during initialization:', error);
            showLoading(false);
            // Don't show error on init - app can still work without database
        }
    });
}
// Run on page load
document.addEventListener('DOMContentLoaded', init);
console.log('CV Builder frontend loaded');

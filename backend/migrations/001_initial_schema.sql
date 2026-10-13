-- CV Builder Database Schema
-- Version: 001
-- Description: Initial schema with ESCO and O*NET skill dictionaries

-- ============================================
-- SKILL DICTIONARIES
-- ============================================

-- Skill categories (e.g., programming, tools, frameworks)
CREATE TABLE IF NOT EXISTS skill_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id INTEGER REFERENCES skill_categories(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ESCO Skill taxonomy
CREATE TABLE IF NOT EXISTS esco_skills (
    id SERIAL PRIMARY KEY,
    esco_id VARCHAR(50) NOT NULL UNIQUE,           -- e.g., "1.1.1.1.1"
    uri VARCHAR(255) NOT NULL UNIQUE,             -- ESCO URI
    concept_type VARCHAR(50) NOT NULL,           -- skill, knowledge, competence
    preferred_label VARCHAR(255) NOT NULL,        -- Primary name
    description TEXT,
    alt_labels TEXT[],                            -- Alternative names
    category_id INTEGER REFERENCES skill_categories(id),
    parent_esco_id VARCHAR(50),                  -- Parent in ESCO hierarchy
    level INTEGER,                               -- Hierarchy level
    relevance_score FLOAT DEFAULT 1.0,           -- Relevance weight (0-1)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- O*NET Skill taxonomy
CREATE TABLE IF NOT EXISTS onet_skills (
    id SERIAL PRIMARY KEY,
    onet_code VARCHAR(50) NOT NULL UNIQUE,       -- O*NET-SOC code
    element_id VARCHAR(50),                      -- Specific element ID
    content_model VARCHAR(50),                   -- Skills, Knowledge, Abilities
    name VARCHAR(255) NOT NULL,                  -- Skill name
    description TEXT,
    category VARCHAR(100),                        -- O*NET category
    category_id INTEGER REFERENCES skill_categories(id),
    example TEXT,
    relevance_score FLOAT DEFAULT 1.0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Skill synonyms (alternative names for skills)
CREATE TABLE IF NOT EXISTS skill_synonyms (
    id SERIAL PRIMARY KEY,
    skill_id INTEGER NOT NULL,                  -- References either esco_skills or onet_skills
    skill_type VARCHAR(10) NOT NULL,            -- 'esco' or 'onet'
    synonym VARCHAR(255) NOT NULL,
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(skill_type, skill_id, synonym)
);

-- Map keywords to standardized skills
CREATE TABLE IF NOT EXISTS keyword_skill_map (
    id SERIAL PRIMARY KEY,
    keyword VARCHAR(255) NOT NULL,
    skill_type VARCHAR(10) NOT NULL,            -- 'esco' or 'onet'
    skill_id INTEGER NOT NULL,                 -- ID in esco_skills or onet_skills
    confidence FLOAT DEFAULT 1.0,             -- Confidence score (0-1)
    source VARCHAR(50),                         -- How the mapping was created
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(skill_type, skill_id, keyword)
);

-- ============================================
-- PROFILES & USER DATA
-- ============================================

-- User profiles
CREATE TABLE IF NOT EXISTS profiles (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Profile experience
CREATE TABLE IF NOT EXISTS profile_experience (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    company VARCHAR(255) NOT NULL,
    position VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    current BOOLEAN DEFAULT FALSE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Profile education
CREATE TABLE IF NOT EXISTS profile_education (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    institution VARCHAR(255) NOT NULL,
    degree VARCHAR(255) NOT NULL,
    field_of_study VARCHAR(255),
    start_date DATE,
    end_date DATE,
    gpa NUMERIC(3,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Link profiles to standardized skills (ESCO/O*NET)
CREATE TABLE IF NOT EXISTS profile_skills (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    skill_type VARCHAR(10) NOT NULL,            -- 'esco' or 'onet'
    skill_id INTEGER NOT NULL,                 -- ID in esco_skills or onet_skills
    proficiency VARCHAR(20) CHECK (proficiency IN ('Beginner', 'Intermediate', 'Advanced', 'Expert')),
    years_experience INTEGER DEFAULT 0,
    last_verified TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(profile_id, skill_type, skill_id)
);

-- Profile custom skills (not in dictionaries)
CREATE TABLE IF NOT EXISTS profile_custom_skills (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    proficiency VARCHAR(20) CHECK (proficiency IN ('Beginner', 'Intermediate', 'Advanced', 'Expert')),
    years_experience INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(profile_id, LOWER(name))
);

-- ============================================
-- JOB OFFERS & MATCHING
-- ============================================

-- Job offers
CREATE TABLE IF NOT EXISTS job_offers (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    company VARCHAR(255),
    location VARCHAR(255),
    description TEXT,
    requirements TEXT,
    url VARCHAR(500),
    source VARCHAR(50),                          -- e.g., "linkedin", "indeed", "manual"
    raw_content TEXT,                           -- Original scraped content
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Keywords extracted from job offers
CREATE TABLE IF NOT EXISTS job_keywords (
    id SERIAL PRIMARY KEY,
    job_offer_id INTEGER NOT NULL REFERENCES job_offers(id) ON DELETE CASCADE,
    keyword VARCHAR(255) NOT NULL,
    context TEXT,                               -- Surrounding text
    frequency INTEGER DEFAULT 1,                -- How many times it appears
    relevance FLOAT DEFAULT 1.0,               -- Relevance score
    category VARCHAR(100),                      -- Detected category
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Job requirements (structured)
CREATE TABLE IF NOT EXISTS job_requirements (
    id SERIAL PRIMARY KEY,
    job_offer_id INTEGER NOT NULL REFERENCES job_offers(id) ON DELETE CASCADE,
    requirement_type VARCHAR(50) CHECK (requirement_type IN ('skill', 'experience', 'education', 'certification', 'language')),
    name VARCHAR(255) NOT NULL,
    min_years INTEGER,
    is_required BOOLEAN DEFAULT TRUE,
    weight FLOAT DEFAULT 1.0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- MATCHING RESULTS
-- ============================================

-- Match results between profiles and job offers
CREATE TABLE IF NOT EXISTS match_results (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    profile_id INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    job_offer_id INTEGER NOT NULL REFERENCES job_offers(id) ON DELETE CASCADE,
    overall_score FLOAT NOT NULL,              -- 0-100
    match_grade VARCHAR(20),                    -- A, B, C, D, F
    total_matches INTEGER DEFAULT 0,
    total_requirements INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Individual skill matches
CREATE TABLE IF NOT EXISTS skill_matches (
    id SERIAL PRIMARY KEY,
    match_result_id INTEGER NOT NULL REFERENCES match_results(id) ON DELETE CASCADE,
    profile_skill_type VARCHAR(10),            -- 'esco', 'onet', or 'custom'
    profile_skill_id INTEGER,                 -- ID in profile_skills or profile_custom_skills
    job_requirement_id INTEGER REFERENCES job_requirements(id) ON DELETE SET NULL,
    match_score FLOAT NOT NULL,              -- 0-100
    match_type VARCHAR(50),                   -- exact, partial, fuzzy, etc.
    confidence FLOAT DEFAULT 1.0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Missing skills/requirements
CREATE TABLE IF NOT EXISTS missing_requirements (
    id SERIAL PRIMARY KEY,
    match_result_id INTEGER NOT NULL REFERENCES match_results(id) ON DELETE CASCADE,
    job_requirement_id INTEGER REFERENCES job_requirements(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    requirement_type VARCHAR(50),
    priority VARCHAR(20),                     -- High, Medium, Low
    asked_user BOOLEAN DEFAULT FALSE,          -- Did we ask the user about this?
    user_response VARCHAR(20),                 -- yes, no, maybe
    user_response_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- GENERATED RESUMES
-- ============================================

-- Generated resumes
CREATE TABLE IF NOT EXISTS resumes (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    profile_id INTEGER NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    job_offer_id INTEGER REFERENCES job_offers(id) ON DELETE SET NULL,
    match_result_id INTEGER REFERENCES match_results(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    format VARCHAR(20) NOT NULL,              -- text, markdown, json, html, pdf, docx
    filename VARCHAR(255),
    storage_path VARCHAR(500),                -- Where it's stored
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- INDEXES
-- ============================================

-- Skill dictionary indexes
CREATE INDEX IF NOT EXISTS idx_esco_skills_preferred_label ON esco_skills(preferred_label);
CREATE INDEX IF NOT EXISTS idx_esco_skills_esco_id ON esco_skills(esco_id);
CREATE INDEX IF NOT EXISTS idx_esco_skills_category ON esco_skills(category_id);

CREATE INDEX IF NOT EXISTS idx_onet_skills_name ON onet_skills(name);
CREATE INDEX IF NOT EXISTS idx_onet_skills_onet_code ON onet_skills(onet_code);
CREATE INDEX IF NOT EXISTS idx_onet_skills_category ON onet_skills(category_id);

CREATE INDEX IF NOT EXISTS idx_skill_synonyms_synonym ON skill_synonyms(synonym);
CREATE INDEX IF NOT EXISTS idx_keyword_skill_map_keyword ON keyword_skill_map(keyword);

-- Profile indexes
CREATE INDEX IF NOT EXISTS idx_profiles_email ON profiles(email);
CREATE INDEX IF NOT EXISTS idx_profiles_uuid ON profiles(uuid);

CREATE INDEX IF NOT EXISTS idx_profile_experience_profile ON profile_experience(profile_id);
CREATE INDEX IF NOT EXISTS idx_profile_education_profile ON profile_education(profile_id);
CREATE INDEX IF NOT EXISTS idx_profile_skills_profile ON profile_skills(profile_id);
CREATE INDEX IF NOT EXISTS idx_profile_skills_type_id ON profile_skills(skill_type, skill_id);

-- Job offer indexes
CREATE INDEX IF NOT EXISTS idx_job_offers_uuid ON job_offers(uuid);
CREATE INDEX IF NOT EXISTS idx_job_keywords_job ON job_keywords(job_offer_id);
CREATE INDEX IF NOT EXISTS idx_job_keywords_keyword ON job_keywords(keyword);
CREATE INDEX IF NOT EXISTS idx_job_requirements_job ON job_requirements(job_offer_id);

-- Matching indexes
CREATE INDEX IF NOT EXISTS idx_match_results_profile_job ON match_results(profile_id, job_offer_id);
CREATE INDEX IF NOT EXISTS idx_skill_matches_result ON skill_matches(match_result_id);
CREATE INDEX IF NOT EXISTS idx_missing_requirements_result ON missing_requirements(match_result_id);

-- Resume indexes
CREATE INDEX IF NOT EXISTS idx_resumes_profile ON resumes(profile_id);
CREATE INDEX IF NOT EXISTS idx_resumes_job ON resumes(job_offer_id);

-- ============================================
-- TRIGGERS
-- ============================================

-- Update updated_at timestamp on profile changes
CREATE OR REPLACE FUNCTION update_profile_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_profile_timestamp
BEFORE UPDATE ON profiles
FOR EACH ROW
EXECUTE FUNCTION update_profile_timestamp();

-- ============================================
-- VIEWS
-- ============================================

-- View for all skills (combined ESCO and O*NET)
CREATE OR REPLACE VIEW all_skills AS
SELECT 
    'esco' AS skill_type,
    id,
    esco_id AS external_id,
    preferred_label AS name,
    description,
    category_id,
    relevance_score
FROM esco_skills
UNION ALL
SELECT 
    'onet' AS skill_type,
    id,
    onet_code AS external_id,
    name,
    description,
    category_id,
    relevance_score
FROM onet_skills;

-- View for profile skills with standardized skill details
CREATE OR REPLACE VIEW profile_skills_with_details AS
SELECT 
    ps.id,
    ps.profile_id,
    ps.skill_type,
    ps.skill_id,
    ps.proficiency,
    ps.years_experience,
    ps.last_verified,
    ps.notes,
    ps.created_at,
    COALESCE(
        (SELECT preferred_label FROM esco_skills WHERE id = ps.skill_id AND ps.skill_type = 'esco'),
        (SELECT name FROM onet_skills WHERE id = ps.skill_id AND ps.skill_type = 'onet')
    ) AS skill_name,
    COALESCE(
        (SELECT description FROM esco_skills WHERE id = ps.skill_id AND ps.skill_type = 'esco'),
        (SELECT description FROM onet_skills WHERE id = ps.skill_id AND ps.skill_type = 'onet')
    ) AS skill_description
FROM profile_skills ps;

-- View for keyword to skill mapping
CREATE OR REPLACE VIEW keyword_to_skill_mapping AS
SELECT 
    ksm.keyword,
    ksm.skill_type,
    ksm.skill_id,
    ksm.confidence,
    COALESCE(
        (SELECT preferred_label FROM esco_skills WHERE id = ksm.skill_id AND ksm.skill_type = 'esco'),
        (SELECT name FROM onet_skills WHERE id = ksm.skill_id AND ksm.skill_type = 'onet')
    ) AS skill_name
FROM keyword_skill_map ksm;

-- ============================================
-- INITIAL DATA
-- ============================================

-- Insert basic skill categories
INSERT INTO skill_categories (name, description) VALUES 
    ('Programming', 'Programming languages and development'),
    ('Frameworks', 'Development frameworks and libraries'),
    ('Tools', 'Development and operations tools'),
    ('Databases', 'Database technologies'),
    ('Cloud', 'Cloud platforms and services'),
    ('DevOps', 'DevOps practices and methodologies'),
    ('Network', 'Networking technologies'),
    ('Security', 'Security practices and technologies'),
    ('Methodologies', 'Development methodologies'),
    ('Testing', 'Testing frameworks and practices')
ON CONFLICT (name) DO NOTHING;

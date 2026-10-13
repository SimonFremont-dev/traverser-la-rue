-- +atlas:sum "Create skill database schema with sources and types"

-- ============================================
-- SOURCE TABLE
-- ============================================
-- Defines where skills come from (ESCO, O*NET, personal)

CREATE TABLE IF NOT EXISTS source (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_official BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- TYPE TABLE  
-- ============================================
-- Defines skill types (skill, tool, framework, etc.)

CREATE TABLE IF NOT EXISTS type (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- SKILL TABLE
-- ============================================
-- Main skill table with source, type, and name

CREATE TABLE IF NOT EXISTS skill (
    id SERIAL PRIMARY KEY,
    source_id INTEGER NOT NULL REFERENCES source(id) ON DELETE CASCADE,
    type_id INTEGER NOT NULL REFERENCES type(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    external_id VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_id, external_id),
    UNIQUE(source_id, LOWER(name))
);

-- ============================================
-- INDEXES
-- ============================================

CREATE INDEX IF NOT EXISTS idx_skill_name ON skill(LOWER(name));
CREATE INDEX IF NOT EXISTS idx_skill_source ON skill(source_id);
CREATE INDEX IF NOT EXISTS idx_skill_type ON skill(type_id);
CREATE INDEX IF NOT EXISTS idx_skill_source_type ON skill(source_id, type_id);

-- ============================================
-- SCHEMA MIGRATIONS TRACKING
-- ============================================

CREATE TABLE IF NOT EXISTS atlas_schema_revisions (
    version VARCHAR(255) PRIMARY KEY,
    description TEXT,
    executed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

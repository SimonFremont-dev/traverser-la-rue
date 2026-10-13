-- Sample ESCO and O*NET data for CV Builder
-- Version: 002
-- Description: Sample skill dictionaries for testing

-- ============================================
-- ESCO SKILLS - Sample Data
-- ============================================

-- Programming Languages
INSERT INTO esco_skills (esco_id, uri, concept_type, preferred_label, description, category_id, parent_esco_id, level, relevance_score) 
VALUES 
    ('1.1.1.1.1', 'http://data.europa.eu/esco/skill/1.1.1.1.1', 'skill', 'Go', 'Go programming language', 1, NULL, 1, 1.0),
    ('1.1.1.1.2', 'http://data.europa.eu/esco/skill/1.1.1.1.2', 'skill', 'Python', 'Python programming language', 1, NULL, 1, 1.0),
    ('1.1.1.1.3', 'http://data.europa.eu/esco/skill/1.1.1.1.3', 'skill', 'Java', 'Java programming language', 1, NULL, 1, 1.0),
    ('1.1.1.1.4', 'http://data.europa.eu/esco/skill/1.1.1.1.4', 'skill', 'JavaScript', 'JavaScript programming language', 1, NULL, 1, 1.0),
    ('1.1.1.1.5', 'http://data.europa.eu/esco/skill/1.1.1.1.5', 'skill', 'TypeScript', 'TypeScript programming language', 1, NULL, 1, 1.0),
    ('1.1.1.1.6', 'http://data.europa.eu/esco/skill/1.1.1.1.6', 'skill', 'C', 'C programming language', 1, NULL, 1, 0.9),
    ('1.1.1.1.7', 'http://data.europa.eu/esco/skill/1.1.1.1.7', 'skill', 'C++', 'C++ programming language', 1, NULL, 1, 0.9),
    ('1.1.1.1.8', 'http://data.europa.eu/esco/skill/1.1.1.1.8', 'skill', 'C#', 'C# programming language', 1, NULL, 1, 0.8),
    ('1.1.1.1.9', 'http://data.europa.eu/esco/skill/1.1.1.1.9', 'skill', 'Rust', 'Rust programming language', 1, NULL, 1, 0.9),
    ('1.1.1.1.10', 'http://data.europa.eu/esco/skill/1.1.1.1.10', 'skill', 'Kotlin', 'Kotlin programming language', 1, NULL, 1, 0.8)
ON CONFLICT (esco_id) DO NOTHING;

-- Tools and Technologies
INSERT INTO esco_skills (esco_id, uri, concept_type, preferred_label, description, category_id, parent_esco_id, level, relevance_score) 
VALUES 
    ('2.1.1.1.1', 'http://data.europa.eu/esco/skill/2.1.1.1.1', 'skill', 'Docker', 'Docker container platform', 2, NULL, 1, 1.0),
    ('2.1.1.1.2', 'http://data.europa.eu/esco/skill/2.1.1.1.2', 'skill', 'Kubernetes', 'Kubernetes orchestration', 2, NULL, 1, 1.0),
    ('2.1.1.1.3', 'http://data.europa.eu/esco/skill/2.1.1.1.3', 'skill', 'Git', 'Git version control', 2, NULL, 1, 1.0),
    ('2.1.1.1.4', 'http://data.europa.eu/esco/skill/2.1.1.1.4', 'skill', 'Terraform', 'Infrastructure as code tool', 2, NULL, 1, 0.9),
    ('2.1.1.1.5', 'http://data.europa.eu/esco/skill/2.1.1.1.5', 'skill', 'Ansible', 'Configuration management tool', 2, NULL, 1, 0.9),
    ('2.1.1.1.6', 'http://data.europa.eu/esco/skill/2.1.1.1.6', 'skill', 'Jenkins', 'Continuous integration tool', 2, NULL, 1, 0.9),
    ('2.1.1.1.7', 'http://data.europa.eu/esco/skill/2.1.1.1.7', 'skill', 'Helm', 'Kubernetes package manager', 2, NULL, 1, 0.8),
    ('2.1.1.1.8', 'http://data.europa.eu/esco/skill/2.1.1.1.8', 'skill', 'Prometheus', 'Monitoring system', 2, NULL, 1, 0.8),
    ('2.1.1.1.9', 'http://data.europa.eu/esco/skill/2.1.1.1.9', 'skill', 'Grafana', 'Visualization platform', 2, NULL, 1, 0.8),
    ('2.1.1.1.10', 'http://data.europa.eu/esco/skill/2.1.1.1.10', 'skill', 'ELK Stack', 'Elasticsearch, Logstash, Kibana', 2, NULL, 1, 0.8)
ON CONFLICT (esco_id) DO NOTHING;

-- Databases
INSERT INTO esco_skills (esco_id, uri, concept_type, preferred_label, description, category_id, parent_esco_id, level, relevance_score) 
VALUES 
    ('3.1.1.1.1', 'http://data.europa.eu/esco/skill/3.1.1.1.1', 'skill', 'PostgreSQL', 'PostgreSQL database', 3, NULL, 1, 1.0),
    ('3.1.1.1.2', 'http://data.europa.eu/esco/skill/3.1.1.1.2', 'skill', 'MySQL', 'MySQL database', 3, NULL, 1, 1.0),
    ('3.1.1.1.3', 'http://data.europa.eu/esco/skill/3.1.1.1.3', 'skill', 'MongoDB', 'MongoDB NoSQL database', 3, NULL, 1, 1.0),
    ('3.1.1.1.4', 'http://data.europa.eu/esco/skill/3.1.1.1.4', 'skill', 'Redis', 'Redis in-memory database', 3, NULL, 1, 0.9),
    ('3.1.1.1.5', 'http://data.europa.eu/esco/skill/3.1.1.1.5', 'skill', 'SQLite', 'SQLite embedded database', 3, NULL, 1, 0.8),
    ('3.1.1.1.6', 'http://data.europa.eu/esco/skill/3.1.1.1.6', 'skill', 'Cassandra', 'Cassandra distributed database', 3, NULL, 1, 0.8)
ON CONFLICT (esco_id) DO NOTHING;

-- Cloud Platforms
INSERT INTO esco_skills (esco_id, uri, concept_type, preferred_label, description, category_id, parent_esco_id, level, relevance_score) 
VALUES 
    ('4.1.1.1.1', 'http://data.europa.eu/esco/skill/4.1.1.1.1', 'skill', 'AWS', 'Amazon Web Services', 4, NULL, 1, 1.0),
    ('4.1.1.1.2', 'http://data.europa.eu/esco/skill/4.1.1.1.2', 'skill', 'Azure', 'Microsoft Azure', 4, NULL, 1, 1.0),
    ('4.1.1.1.3', 'http://data.europa.eu/esco/skill/4.1.1.1.3', 'skill', 'GCP', 'Google Cloud Platform', 4, NULL, 1, 1.0),
    ('4.1.1.1.4', 'http://data.europa.eu/esco/skill/4.1.1.1.4', 'skill', 'EC2', 'Amazon Elastic Compute Cloud', 4, '4.1.1.1.1', 2, 0.9),
    ('4.1.1.1.5', 'http://data.europa.eu/esco/skill/4.1.1.1.5', 'skill', 'S3', 'Amazon Simple Storage Service', 4, '4.1.1.1.1', 2, 0.9),
    ('4.1.1.1.6', 'http://data.europa.eu/esco/skill/4.1.1.1.6', 'skill', 'Lambda', 'AWS Lambda serverless', 4, '4.1.1.1.1', 2, 0.9),
    ('4.1.1.1.7', 'http://data.europa.eu/esco/skill/4.1.1.1.7', 'skill', 'EKS', 'Amazon Elastic Kubernetes Service', 4, '4.1.1.1.1', 2, 0.8)
ON CONFLICT (esco_id) DO NOTHING;

-- Frameworks
INSERT INTO esco_skills (esco_id, uri, concept_type, preferred_label, description, category_id, parent_esco_id, level, relevance_score) 
VALUES 
    ('5.1.1.1.1', 'http://data.europa.eu/esco/skill/5.1.1.1.1', 'skill', 'React', 'React JavaScript framework', 5, NULL, 1, 1.0),
    ('5.1.1.1.2', 'http://data.europa.eu/esco/skill/5.1.1.1.2', 'skill', 'Angular', 'Angular JavaScript framework', 5, NULL, 1, 0.9),
    ('5.1.1.1.3', 'http://data.europa.eu/esco/skill/5.1.1.1.3', 'skill', 'Vue.js', 'Vue JavaScript framework', 5, NULL, 1, 0.9),
    ('5.1.1.1.4', 'http://data.europa.eu/esco/skill/5.1.1.1.4', 'skill', 'Spring Boot', 'Spring Boot Java framework', 5, NULL, 1, 0.9),
    ('5.1.1.1.5', 'http://data.europa.eu/esco/skill/5.1.1.1.5', 'skill', 'Django', 'Django Python framework', 5, NULL, 1, 0.9),
    ('5.1.1.1.6', 'http://data.europa.eu/esco/skill/5.1.1.1.6', 'skill', 'Flask', 'Flask Python framework', 5, NULL, 1, 0.8)
ON CONFLICT (esco_id) DO NOTHING;

-- ============================================
-- O*NET SKILLS - Sample Data
-- ============================================

-- Skills from O*NET (using simplified codes)
INSERT INTO onet_skills (onet_code, element_id, content_model, name, description, category, category_id, example, relevance_score) 
VALUES 
    -- Programming
    ('2.S.1.1', '1.A.1', 'Skills', 'Programming', 'Writing computer programs', 'Technical', 1, 'Developing software applications', 1.0),
    ('2.S.1.2', '1.A.2', 'Skills', 'Software Development', 'Developing software applications', 'Technical', 1, 'Building applications using programming languages', 1.0),
    ('2.S.1.3', '1.A.3', 'Skills', 'Debugging', 'Finding and fixing errors in code', 'Technical', 1, 'Identifying and resolving software bugs', 0.9),
    
    -- Tools
    ('2.S.2.1', '1.A.4', 'Skills', 'Containerization', 'Creating and managing containers', 'Technical', 2, 'Using Docker to package applications', 1.0),
    ('2.S.2.2', '1.A.5', 'Skills', 'Orchestration', 'Managing container orchestration', 'Technical', 2, 'Using Kubernetes to deploy applications', 1.0),
    ('2.S.2.3', '1.A.6', 'Skills', 'Version Control', 'Managing code versions', 'Technical', 2, 'Using Git for source control', 1.0),
    
    -- Databases
    ('2.S.3.1', '1.A.7', 'Skills', 'Database Design', 'Designing database structures', 'Technical', 3, 'Creating database schemas', 1.0),
    ('2.S.3.2', '1.A.8', 'Skills', 'SQL', 'Structured Query Language', 'Technical', 3, 'Writing database queries', 1.0),
    ('2.S.3.3', '1.A.9', 'Skills', 'NoSQL', 'Non-relational database management', 'Technical', 3, 'Using MongoDB or similar databases', 0.9),
    
    -- Cloud
    ('2.S.4.1', '1.A.10', 'Skills', 'Cloud Computing', 'Using cloud platforms', 'Technical', 4, 'Deploying to AWS, Azure, or GCP', 1.0),
    ('2.S.4.2', '1.A.11', 'Skills', 'Serverless Architecture', 'Building serverless applications', 'Technical', 4, 'Using AWS Lambda or similar', 0.9)
ON CONFLICT (onet_code, element_id) DO NOTHING;

-- ============================================
-- KEYWORD TO SKILL MAPPINGS
-- ============================================

-- Map common keyword variations to standardized skills
INSERT INTO keyword_skill_map (keyword, skill_type, skill_id, confidence, source) 
VALUES 
    -- Go variations
    ('go', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.1'), 1.0, 'manual'),
    ('golang', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.1'), 1.0, 'manual'),
    ('go programming', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.1'), 0.9, 'manual'),
    
    -- Python variations
    ('python', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.2'), 1.0, 'manual'),
    ('python3', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.2'), 0.9, 'manual'),
    ('py', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.2'), 0.8, 'manual'),
    
    -- JavaScript variations
    ('javascript', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.4'), 1.0, 'manual'),
    ('js', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.4'), 0.9, 'manual'),
    ('typescript', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.5'), 1.0, 'manual'),
    ('ts', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.5'), 0.9, 'manual'),
    
    -- Docker variations
    ('docker', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.1'), 1.0, 'manual'),
    ('docker container', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.1'), 0.9, 'manual'),
    ('dockerfile', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.1'), 0.8, 'manual'),
    
    -- Kubernetes variations
    ('kubernetes', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.2'), 1.0, 'manual'),
    ('k8s', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.2'), 1.0, 'manual'),
    ('kubernetes orchestration', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.2'), 0.9, 'manual'),
    ('kube', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.2'), 0.8, 'manual'),
    
    -- Git variations
    ('git', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.3'), 1.0, 'manual'),
    ('github', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.3'), 0.9, 'manual'),
    ('gitlab', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.3'), 0.9, 'manual'),
    
    -- Database variations
    ('postgres', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '3.1.1.1.1'), 0.9, 'manual'),
    ('postgresql', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '3.1.1.1.1'), 1.0, 'manual'),
    ('mysql', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '3.1.1.1.2'), 1.0, 'manual'),
    ('mongodb', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '3.1.1.1.3'), 1.0, 'manual'),
    
    -- AWS variations
    ('aws', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '4.1.1.1.1'), 1.0, 'manual'),
    ('amazon web services', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '4.1.1.1.1'), 1.0, 'manual'),
    ('amazon', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '4.1.1.1.1'), 0.8, 'manual'),
    
    -- Azure variations
    ('azure', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '4.1.1.1.2'), 1.0, 'manual'),
    ('microsoft azure', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '4.1.1.1.2'), 1.0, 'manual'),
    
    -- GCP variations
    ('gcp', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '4.1.1.1.3'), 0.9, 'manual'),
    ('google cloud', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '4.1.1.1.3'), 1.0, 'manual'),
    ('google cloud platform', 'esco', (SELECT id FROM esco_skills WHERE esco_id = '4.1.1.1.3'), 1.0, 'manual')
ON CONFLICT (skill_type, skill_id, keyword) DO NOTHING;

-- ============================================
-- SKILL SYNONYMS
-- ============================================

-- Add synonyms for skills
INSERT INTO skill_synonyms (skill_type, skill_id, synonym, language) 
VALUES 
    ('esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.1'), 'Golang', 'en'),
    ('esco', (SELECT id FROM esco_skills WHERE esco_id = '1.1.1.1.4'), 'JS', 'en'),
    ('esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.2'), 'K8s', 'en'),
    ('esco', (SELECT id FROM esco_skills WHERE esco_id = '2.1.1.1.2'), 'Kubernetes', 'en')
ON CONFLICT (skill_type, skill_id, synonym) DO NOTHING;

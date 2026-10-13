-- +atlas:sum "Insert sample sources, types, and skills"

-- ============================================
-- SOURCES
-- ============================================

-- Official sources (used for matching)
INSERT INTO source (name, description, is_official) VALUES 
    ('ESCO', 'European Skills/Competences, Qualifications and Occupations taxonomy', TRUE),
    ('O*NET', 'US Occupational Information Network taxonomy', TRUE)
ON CONFLICT (name) DO NOTHING;

-- Personal/unofficial sources
INSERT INTO source (name, description, is_official) VALUES 
    ('Personal', 'User-defined personal skills', FALSE)
ON CONFLICT (name) DO NOTHING;

-- ============================================
-- TYPES
-- ============================================

INSERT INTO type (name, description) VALUES 
    ('Programming Language', 'Programming and scripting languages'),
    ('Tool', 'Development and operations tools'),
    ('Framework', 'Development frameworks and libraries'),
    ('Database', 'Database technologies'),
    ('Cloud', 'Cloud platforms and services'),
    ('DevOps', 'DevOps practices and methodologies'),
    ('Network', 'Networking technologies and protocols'),
    ('Security', 'Security practices and technologies'),
    ('Methodology', 'Development methodologies'),
    ('Testing', 'Testing frameworks and practices'),
    ('Operating System', 'Operating systems'),
    ('Other', 'Other skills and technologies')
ON CONFLICT (name) DO NOTHING;

-- ============================================
-- SKILLS FROM ESCO (Official)
-- ============================================

-- Programming Languages
INSERT INTO skill (source_id, type_id, name, external_id, description) 
VALUES 
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Programming Language'), 'Go', '1.1.1.1.1', 'Go programming language'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Programming Language'), 'Golang', '1.1.1.1.1-alt', 'Go programming language (alternative name)'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Programming Language'), 'Python', '1.1.1.1.2', 'Python programming language'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Programming Language'), 'Java', '1.1.1.1.3', 'Java programming language'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Programming Language'), 'JavaScript', '1.1.1.1.4', 'JavaScript programming language'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Programming Language'), 'TypeScript', '1.1.1.1.5', 'TypeScript programming language'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Programming Language'), 'C', '1.1.1.1.6', 'C programming language'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Programming Language'), 'C++', '1.1.1.1.7', 'C++ programming language'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Programming Language'), 'C#', '1.1.1.1.8', 'C# programming language'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Programming Language'), 'Rust', '1.1.1.1.9', 'Rust programming language'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Programming Language'), 'Kotlin', '1.1.1.1.10', 'Kotlin programming language')
ON CONFLICT (source_id, LOWER(name)) DO NOTHING;

-- Tools
INSERT INTO skill (source_id, type_id, name, external_id, description) 
VALUES 
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Tool'), 'Docker', '2.1.1.1.1', 'Docker container platform'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Tool'), 'Kubernetes', '2.1.1.1.2', 'Kubernetes orchestration system'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Tool'), 'k8s', '2.1.1.1.2-alt', 'Kubernetes (short name)'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Tool'), 'Git', '2.1.1.1.3', 'Git version control system'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Tool'), 'Terraform', '2.1.1.1.4', 'Infrastructure as code tool'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Tool'), 'Ansible', '2.1.1.1.5', 'Configuration management tool'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Tool'), 'Jenkins', '2.1.1.1.6', 'Continuous integration server'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Tool'), 'Helm', '2.1.1.1.7', 'Kubernetes package manager'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Tool'), 'Prometheus', '2.1.1.1.8', 'Monitoring and alerting toolkit'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Tool'), 'Grafana', '2.1.1.1.9', 'Analytics and monitoring platform')
ON CONFLICT (source_id, LOWER(name)) DO NOTHING;

-- Databases
INSERT INTO skill (source_id, type_id, name, external_id, description) 
VALUES 
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Database'), 'PostgreSQL', '3.1.1.1.1', 'PostgreSQL relational database'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Database'), 'Postgres', '3.1.1.1.1-alt', 'PostgreSQL (short name)'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Database'), 'MySQL', '3.1.1.1.2', 'MySQL relational database'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Database'), 'MongoDB', '3.1.1.1.3', 'MongoDB NoSQL database'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Database'), 'Redis', '3.1.1.1.4', 'Redis in-memory database'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Database'), 'SQLite', '3.1.1.1.5', 'SQLite embedded database'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Database'), 'Cassandra', '3.1.1.1.6', 'Cassandra distributed database')
ON CONFLICT (source_id, LOWER(name)) DO NOTHING;

-- Cloud Platforms
INSERT INTO skill (source_id, type_id, name, external_id, description) 
VALUES 
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Cloud'), 'AWS', '4.1.1.1.1', 'Amazon Web Services'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Cloud'), 'Amazon Web Services', '4.1.1.1.1-alt', 'AWS (full name)'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Cloud'), 'Azure', '4.1.1.1.2', 'Microsoft Azure'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Cloud'), 'GCP', '4.1.1.1.3', 'Google Cloud Platform'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Cloud'), 'Google Cloud Platform', '4.1.1.1.3-alt', 'GCP (full name)'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Cloud'), 'EC2', '4.1.1.1.4', 'Amazon Elastic Compute Cloud'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Cloud'), 'S3', '4.1.1.1.5', 'Amazon Simple Storage Service'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Cloud'), 'Lambda', '4.1.1.1.6', 'AWS Lambda serverless'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Cloud'), 'EKS', '4.1.1.1.7', 'Amazon Elastic Kubernetes Service')
ON CONFLICT (source_id, LOWER(name)) DO NOTHING;

-- Frameworks
INSERT INTO skill (source_id, type_id, name, external_id, description) 
VALUES 
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Framework'), 'React', '5.1.1.1.1', 'React JavaScript framework'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Framework'), 'Angular', '5.1.1.1.2', 'Angular JavaScript framework'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Framework'), 'Vue.js', '5.1.1.1.3', 'Vue JavaScript framework'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Framework'), 'Spring Boot', '5.1.1.1.4', 'Spring Boot Java framework'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Framework'), 'Django', '5.1.1.1.5', 'Django Python framework'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Framework'), 'Flask', '5.1.1.1.6', 'Flask Python framework'),
    ((SELECT id FROM source WHERE name = 'ESCO'), (SELECT id FROM type WHERE name = 'Framework'), 'Gin', '5.1.1.1.7', 'Gin Go web framework')
ON CONFLICT (source_id, LOWER(name)) DO NOTHING;

-- ============================================
-- SKILLS FROM O*NET (Official)
-- ============================================

-- Add O*NET skills with different external IDs
INSERT INTO skill (source_id, type_id, name, external_id, description) 
VALUES 
    ((SELECT id FROM source WHERE name = 'O*NET'), (SELECT id FROM type WHERE name = 'Programming Language'), 'Python', '2.S.1.1', 'Programming using Python'),
    ((SELECT id FROM source WHERE name = 'O*NET'), (SELECT id FROM type WHERE name = 'Tool'), 'Docker', '2.S.2.1', 'Containerization with Docker'),
    ((SELECT id FROM source WHERE name = 'O*NET'), (SELECT id FROM type WHERE name = 'Tool'), 'Kubernetes', '2.S.2.2', 'Orchestration with Kubernetes'),
    ((SELECT id FROM source WHERE name = 'O*NET'), (SELECT id FROM type WHERE name = 'Database'), 'SQL', '2.S.3.1', 'Database querying with SQL')
ON CONFLICT (source_id, LOWER(name)) DO NOTHING;

-- ============================================
-- RECORD MIGRATION
-- ============================================

-- Record that this migration was executed
INSERT INTO atlas_schema_revisions (version, description) 
VALUES ('20261011_000002_insert_sample_data', 'Insert sample sources, types, and skills')
ON CONFLICT (version) DO NOTHING;

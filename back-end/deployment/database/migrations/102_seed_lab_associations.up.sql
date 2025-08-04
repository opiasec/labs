-- Seed Lab Associations
-- This file creates associations between labs and their technologies, languages, and vulnerabilities

-- First, let's create variables for the lab IDs (we'll use the slugs to find them)
-- Lab 1: example-lab
-- Lab 2: copy-n-paste-lab

-- Associate example-lab with technologies
INSERT INTO lab_technologies (lab_id, technology_id)
SELECT l.id, t.id
FROM labs l, technologies t
WHERE l.slug = 'example-lab'
AND t.name IN ('Node.js', 'Express.js', 'Docker', 'PostgreSQL', 'JWT');

-- Associate example-lab with languages
INSERT INTO lab_languages (lab_id, language_id)
SELECT l.id, lang.id
FROM labs l, languages lang
WHERE l.slug = 'example-lab'
AND lang.name IN ('JavaScript', 'HTML', 'CSS');

-- Associate example-lab with vulnerabilities (generic example vulnerabilities)
INSERT INTO lab_vulnerabilities (lab_id, vulnerability_id)
SELECT l.id, v.id
FROM labs l, vulnerabilities v
WHERE l.slug = 'example-lab'
AND v.name IN ('Cross-Site Scripting (XSS)', 'Insecure Direct Object References', 'Security Misconfiguration');

-- Associate copy-n-paste-lab with technologies
INSERT INTO lab_technologies (lab_id, technology_id)
SELECT l.id, t.id
FROM labs l, technologies t
WHERE l.slug = 'copy-n-paste-lab'
AND t.name IN ('Docker', 'MySQL', 'REST API');

-- Associate copy-n-paste-lab with languages
INSERT INTO lab_languages (lab_id, language_id)
SELECT l.id, lang.id
FROM labs l, languages lang
WHERE l.slug = 'copy-n-paste-lab'
AND lang.name IN ('Go', 'HTML', 'CSS', 'JavaScript', 'SQL');

-- Associate copy-n-paste-lab with vulnerabilities
INSERT INTO lab_vulnerabilities (lab_id, vulnerability_id)
SELECT l.id, v.id
FROM labs l, vulnerabilities v
WHERE l.slug = 'copy-n-paste-lab'
AND v.name IN ('SQL Injection', 'Injection', 'Broken Authentication');

-- Additional associations for more comprehensive coverage

-- Add more technologies to copy-n-paste-lab based on its description
INSERT INTO lab_technologies (lab_id, technology_id)
SELECT l.id, t.id
FROM labs l, technologies t
WHERE l.slug = 'copy-n-paste-lab'
AND t.name IN ('Docker')
ON CONFLICT (lab_id, technology_id) DO NOTHING;

-- Add more vulnerabilities related to SQL Injection lab
INSERT INTO lab_vulnerabilities (lab_id, vulnerability_id)
SELECT l.id, v.id
FROM labs l, vulnerabilities v
WHERE l.slug = 'copy-n-paste-lab'
AND v.name IN ('Command Injection', 'Insecure Direct Object References', 'Missing Function Level Access Control')
ON CONFLICT (lab_id, vulnerability_id) DO NOTHING;

-- Remove lab associations
DELETE FROM lab_vulnerabilities WHERE lab_id IN (
    SELECT id FROM labs WHERE slug IN ('example-lab', 'copy-n-paste-lab')
);

DELETE FROM lab_languages WHERE lab_id IN (
    SELECT id FROM labs WHERE slug IN ('example-lab', 'copy-n-paste-lab')
);

DELETE FROM lab_technologies WHERE lab_id IN (
    SELECT id FROM labs WHERE slug IN ('example-lab', 'copy-n-paste-lab')
);

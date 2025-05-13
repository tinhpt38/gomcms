-- Add counter column to attendance_checkins table if it doesn't exist
-- This migration ensures that existing tables get the counter field
-- It's safe to run even if the column already exists (it will skip)

-- MySQL version
SET @counterColumnExists = 0;
SELECT COUNT(*) INTO @counterColumnExists FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'attendance_checkins' AND COLUMN_NAME = 'counter';

SET @alterTableStmt = IF(@counterColumnExists = 0, 
                       'ALTER TABLE attendance_checkins ADD COLUMN counter int DEFAULT 1 COMMENT "Counter";',
                       'SELECT "Counter column already exists" AS message;');

PREPARE stmt FROM @alterTableStmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Initialize counter values for existing records
-- This will set counter=1 for all records that don't have a counter value set
UPDATE attendance_checkins SET counter = 1 WHERE counter IS NULL;

-- SQL script to verify and clean up counter values

-- MySQL version
-- 1. Ensure the counter field exists
SET @counterColumnExists = 0;
SELECT COUNT(*) INTO @counterColumnExists FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'attendance_checkins' AND COLUMN_NAME = 'counter';

SET @alterTableStmt = IF(@counterColumnExists = 0, 
                       'ALTER TABLE attendance_checkins ADD COLUMN counter int DEFAULT 1 COMMENT "Counter";',
                       'SELECT "Counter column already exists" AS message;');

PREPARE stmt FROM @alterTableStmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 2. Fix NULL or zero counter values (set to 1)
UPDATE attendance_checkins SET counter = 1 WHERE counter IS NULL OR counter <= 0;

-- 3. Calculate counters for all records based on unique combinations 
-- This query will find duplicate check-ins and update counter values correctly
CREATE TEMPORARY TABLE temp_counter_calc AS
SELECT 
    attendance_id, 
    partpaticipant_id, 
    condition_id, 
    COUNT(*) as correct_count,
    MIN(id) as first_record_id
FROM 
    attendance_checkins
WHERE 
    condition_id IS NOT NULL
GROUP BY 
    attendance_id, partpaticipant_id, condition_id
HAVING 
    COUNT(*) > 1;

-- 4. Update the counter for the first record of each group
UPDATE attendance_checkins a
JOIN temp_counter_calc t ON a.id = t.first_record_id
SET a.counter = t.correct_count;

-- 5. Delete duplicate records (optional - uncomment if you want to remove duplicates)
-- DELETE a FROM attendance_checkins a
-- JOIN temp_counter_calc t ON a.attendance_id = t.attendance_id 
--     AND a.partpaticipant_id = t.partpaticipant_id 
--     AND a.condition_id = t.condition_id
-- WHERE a.id != t.first_record_id;

-- 6. Clean up
DROP TEMPORARY TABLE IF EXISTS temp_counter_calc;

-- 7. Display statistics about attendance check-ins
SELECT 
    attendance_id, 
    partpaticipant_id, 
    condition_id, 
    counter,
    checkin_date
FROM 
    attendance_checkins
ORDER BY 
    attendance_id, partpaticipant_id, condition_id, checkin_date;

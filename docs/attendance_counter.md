# Attendance Check-in Counter System

## Overview

The attendance check-in counter system tracks how many times a user has checked in for each condition. Instead of creating a new record for each check-in with the same condition, it increments a counter field in the existing record.

## Implementation Details

### Database Schema

The `attendance_checkins` table has a `counter` field that stores the number of times a user has checked in for a specific condition:

```go
type AttendanceCheckIn struct {
    // ... other fields
    Counter int `json:"counter" form:"counter" gorm:"column:counter;comment:;type:int"` //Counter
}
```

### Backend Logic

1. When a user checks in, the system checks if a record with the same conditions already exists:
   - Same participant
   - Same attendance session
   - Same condition

2. If a record exists:
   - The counter field is incremented
   - The location, IP, and timestamp are updated with the latest data
   - The existing record is saved

3. If no record exists:
   - A new record is created with counter = 1

### API Response

The API response includes:
- The counter information in the `checkins` array
- Counter values embedded in the conditions data

### Frontend Display

The frontend displays the counter value:
1. In the main check-in UI, showing how many times a user has checked in for each condition
2. In the success message after a check-in

## Utility Functions

Utility functions in `utils/checkin_counter.go` provide reusable methods for:
- Getting the current counter value
- Incrementing a counter
- Retrieving all counters for a participant in an attendance session

## Database Migration

If needed, migration scripts are available in:
- `server/resource/db/migrations/001_add_counter_column.sql` (Adds the counter column)
- `server/resource/db/migrations/002_verify_counter_values.sql` (Verifies and fixes counter values)

## Future Enhancements

Potential future enhancements could include:
1. Statistics dashboard showing check-in frequency
2. Rate limiting to prevent excessive check-ins
3. Detailed history view for administrators
4. Graphical representation of check-in patterns

package utils

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	"gorm.io/gorm"
)

// GetCheckInCounter retrieves the current counter value for a specific check-in condition
func GetCheckInCounter(participantID uint, attendanceID uint, conditionID uint) int {
	var checkIn checkins.AttendanceCheckIn
	err := global.GVA_DB.Where(&checkins.AttendanceCheckIn{
		PartpaticipantId: &participantID,
		AttendanceId:     &attendanceID,
		ConditionId:      &conditionID,
	}).First(&checkIn).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0
		}
		// If there's a database error, default to 0
		return 0
	}

	return checkIn.Counter
}

// IncrementCheckInCounter increments the counter for an existing check-in record
// Returns the new counter value and any error that occurred
func IncrementCheckInCounter(participantID uint, attendanceID uint, conditionID uint) (int, error) {
	var checkIn checkins.AttendanceCheckIn
	err := global.GVA_DB.Where(&checkins.AttendanceCheckIn{
		PartpaticipantId: &participantID,
		AttendanceId:     &attendanceID,
		ConditionId:      &conditionID,
	}).First(&checkIn).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Record not found, create a new one with counter = 1
			newCheckIn := checkins.AttendanceCheckIn{
				PartpaticipantId: &participantID,
				AttendanceId:     &attendanceID,
				ConditionId:      &conditionID,
				Counter:          1,
			}
			err = global.GVA_DB.Create(&newCheckIn).Error
			if err != nil {
				return 0, err
			}
			return 1, nil
		}
		return 0, err
	}

	// Ensure counter is initialized
	if checkIn.Counter <= 0 {
		checkIn.Counter = 1
	}

	// Increment counter
	checkIn.Counter++
	err = global.GVA_DB.Save(&checkIn).Error
	if err != nil {
		return 0, err
	}

	return checkIn.Counter, nil
}

// GetAllCheckInCounters retrieves all check-in counters for a participant in an attendance session
func GetAllCheckInCounters(participantID uint, attendanceID uint) (map[uint]int, error) {
	var checkIns []checkins.AttendanceCheckIn
	err := global.GVA_DB.Where(&checkins.AttendanceCheckIn{
		PartpaticipantId: &participantID,
		AttendanceId:     &attendanceID,
	}).Find(&checkIns).Error

	if err != nil {
		return nil, err
	}

	counters := make(map[uint]int)
	for _, checkIn := range checkIns {
		if checkIn.ConditionId != nil {
			counters[*checkIn.ConditionId] = checkIn.Counter
		}
	}

	return counters, nil
}

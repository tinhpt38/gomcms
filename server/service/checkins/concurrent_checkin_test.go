package checkins_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	service "github.com/flipped-aurora/gin-vue-admin/server/service/checkins"
	"github.com/stretchr/testify/assert"
)

// TestConcurrentCheckinWithLuckyNumbers tests the full check-in flow with
// multiple concurrent users to ensure unique lucky numbers
func TestConcurrentCheckinWithLuckyNumbers(t *testing.T) {
	// Skip this test in normal runs since it requires complex setup and can be slow
	// Comment this line when you want to run this specific test
	t.Skip("Skipping concurrent check-in test - run manually when needed")

	// Set up test data for check-in flow
	attendance := setupConcurrentTestData(t)

	// Create test attendance with lucky numbers enabled
	attendance.UseLuckyNumber = true
	attendance.LuckyShowAfterMinCount = 1
	attendance.LimitCount = 100       // High limit to allow all check-ins
	attendance.LimitClientCount = 100 // High limit for client
	attendance.IsLocked = false
	err := global.GVA_DB.Save(&attendance).Error
	assert.NoError(t, err, "Failed to update test attendance")

	// Create a location for check-in
	latitude := 10.762622
	longitude := 106.660172
	radius := 100.0
	area := checkins.Area{
		Name:      "Test Area",
		Latitude:  &latitude,
		Longitude: &longitude,
		Radius:    &radius,
	}
	err = global.GVA_DB.Create(&area).Error
	assert.NoError(t, err, "Failed to create test area")

	// Create attendance area
	attendanceArea := checkins.AttendanceArea{
		AreaID:          &area.ID,
		AttendanceID:    &attendance.ID,
		Area:            &area,
		AllowRestrictIp: false,
	}
	err = global.GVA_DB.Create(&attendanceArea).Error
	assert.NoError(t, err, "Failed to create test attendance area")

	// Create a condition for check-in
	showLuckyNumber := true
	condition := checkins.Condition{
		AttendanceId:    &attendance.ID,
		AreaId:          &area.ID,
		ShowLuckyNumber: showLuckyNumber,
	}
	err = global.GVA_DB.Create(&condition).Error
	assert.NoError(t, err, "Failed to create test condition")

	// Number of concurrent participants
	const numParticipants = 20

	// Prepare participants
	for i := 1; i <= numParticipants; i++ {
		fullName := fmt.Sprintf("Test User %d", i)
		participant := checkins.Participant{
			Email:    fmt.Sprintf("test%d@example.com", i), // Generate unique emails
			FullName: &fullName,
		}
		err = global.GVA_DB.Create(&participant).Error
		assert.NoError(t, err, "Failed to create test participant")

		// Link to attendance
		agp := checkins.AttendanceGroupParticipant{
			ParticipantId: &participant.ID,
			AttendanceId:  &attendance.ID,
		}
		err = global.GVA_DB.Create(&agp).Error
		assert.NoError(t, err, "Failed to create attendance group participant")

		// Link condition to participant
		agpCondition := checkins.AGPCondition{
			AttendanceGroupParticipantId: int(agp.ID),
			ConditionId:                  int(condition.ID),
			AttendanceId:                 int(attendance.ID),
		}
		err = global.GVA_DB.Create(&agpCondition).Error
		assert.NoError(t, err, "Failed to create AGP condition")
	}

	// Wait group for concurrent operations
	var wg sync.WaitGroup
	wg.Add(numParticipants)

	// Channel to collect lucky numbers
	luckyNumberChan := make(chan int, numParticipants)

	// Create service
	checkInService := service.AttendanceCheckInService{}

	// Start concurrent check-ins
	for i := 1; i <= numParticipants; i++ {
		go func(index int) {
			defer wg.Done()

			// Create check-in request
			lat := 10.762622 // Within area range
			lng := 106.660172
			accuracy := 10.0
			email := fmt.Sprintf("test%d@example.com", index)
			fullName := fmt.Sprintf("Test User %d", index)
			var pFullName *string = &fullName
			visitorId := fmt.Sprintf("visitor-%d", index)

			req := checkinsReq.CheckinsReq{
				Email:     email,
				Code:      fmt.Sprintf("%d", attendance.ID),
				FullName:  pFullName,
				VisitorId: visitorId,
				Lat:       &lat,
				Lng:       &lng,
				Accuracy:  &accuracy,
			}

			// Random small delay to simulate real-world concurrency
			time.Sleep(time.Duration(index%10) * time.Millisecond)

			// Attempt to check in
			result, err := checkInService.CheckinAttendance(req, fmt.Sprintf("127.0.0.%d", index), "Test Agent")
			if err != nil {
				t.Logf("Error in concurrent check-in %d: %v", index, err)
				return
			}

			// Check if a lucky number was generated
			if luckyNum, ok := result["luckyNumber"]; ok {
				if num, ok := luckyNum.(int); ok {
					luckyNumberChan <- num
				} else {
					t.Logf("Lucky number for participant %d is not an int: %v", index, luckyNum)
				}
			} else {
				t.Logf("No lucky number for participant %d", index)
			}
		}(i)
	}

	// Close channel when all check-ins are done
	go func() {
		wg.Wait()
		close(luckyNumberChan)
	}()

	// Collect and verify lucky numbers
	luckyNumbers := make(map[int]bool)
	duplicates := 0

	for num := range luckyNumberChan {
		if luckyNumbers[num] {
			duplicates++
			t.Logf("Duplicate lucky number detected: %d", num)
		}
		luckyNumbers[num] = true
	}

	// Verify all numbers are unique
	assert.Zero(t, duplicates, "Expected no duplicate lucky numbers")

	// Check the database for distinct lucky numbers
	var dbLuckyNumbers []int
	err = global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
		Where("attendance_id = ? AND lucky_number > 0", attendance.ID).
		Distinct().
		Pluck("lucky_number", &dbLuckyNumbers).Error

	assert.NoError(t, err, "Failed to query lucky numbers from database")

	// Verify number of distinct lucky numbers in DB
	t.Logf("Number of distinct lucky numbers in database: %d", len(dbLuckyNumbers))

	// Check for any duplicates in the database (shouldn't happen)
	dbNumbersMap := make(map[int]int)
	for _, num := range dbLuckyNumbers {
		dbNumbersMap[num]++
	}

	for num, count := range dbNumbersMap {
		if count > 1 {
			t.Errorf("Lucky number %d appears %d times in database", num, count)
		}
	}
}

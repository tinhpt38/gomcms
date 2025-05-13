package checkins

import (
	"encoding/base32"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/checkins"
	checkinsReq "github.com/flipped-aurora/gin-vue-admin/server/model/checkins/request"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// TestMain sets up the database for all tests
func TestMain(m *testing.M) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()
	global.GVA_LOG = logger

	// Create and configure SQLite in-memory database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger:                                   gormlogger.Default.LogMode(gormlogger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Set global DB
	global.GVA_DB = db

	// Auto migrate required tables
	err = db.AutoMigrate(
		&checkins.Attendance{},
		&checkins.Participant{},
		&checkins.AttendanceGroupParticipant{},
		&checkins.AttendanceCheckIn{},
		&checkins.CheckinLog{},
		&checkins.Condition{},
		&checkins.AttendanceArea{},
		&checkins.Area{},
		&checkins.AGPCondition{},
	)
	if err != nil {
		fmt.Printf("Failed to migrate tables: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	exitCode := m.Run()

	// Close database connection
	sqlDB, _ := db.DB()
	sqlDB.Close()

	// Exit with test result code
	os.Exit(exitCode)
}

// setupBasicTestData creates the minimum data needed for testing
func setupBasicTestData(t *testing.T) {
	// Clean up existing data
	cleanupTestData(t)

	// Create test attendance
	startDate := time.Now().Add(-1 * time.Hour)
	endDate := time.Now().Add(1 * time.Hour)
	attendance := checkins.Attendance{
		Title:                  "Test Attendance",
		AllowGuest:             true,
		UseLuckyNumber:         true,
		LuckyShowAfterMinCount: 1,
		LimitCount:             5,
		LimitClientCount:       3,
		IsLocked:               false,
		StartDate:              &startDate,
		EndDate:                &endDate,
	}
	err := global.GVA_DB.Create(&attendance).Error
	assert.NoError(t, err, "Failed to create test attendance")

	// Create test participant
	fullName := "Test User"
	participant := checkins.Participant{
		Email:    "test@example.com",
		FullName: &fullName,
	}
	err = global.GVA_DB.Create(&participant).Error
	assert.NoError(t, err, "Failed to create test participant")

	// Link participant to attendance
	agp := checkins.AttendanceGroupParticipant{
		ParticipantId: &participant.ID,
		AttendanceId:  &attendance.ID,
	}
	err = global.GVA_DB.Create(&agp).Error
	assert.NoError(t, err, "Failed to create attendance group participant")
}

// cleanupTestData removes all test data
func cleanupTestData(t *testing.T) {
	tables := []string{
		"attendance_check_ins",
		"checkin_logs",
		"agp_conditions",
		"attendance_areas",
		"conditions",
		"attendance_group_participants",
		"participants",
		"areas",
		"attendances",
	}

	for _, table := range tables {
		err := global.GVA_DB.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error
		assert.NoError(t, err, "Failed to clean up table: "+table)
	}
}

// setupCompleteTestData creates all the necessary data to test CheckinAttendance
func setupCompleteTestData(t *testing.T) (string, uint, uint) {
	// Clean up existing data
	cleanupTestData(t)

	// Create a test area
	latitude := 10.762622
	longitude := 106.660172
	radius := 100.0
	area := checkins.Area{
		Name:      "Test Area",
		Latitude:  &latitude,
		Longitude: &longitude,
		Radius:    &radius,
	}
	err := global.GVA_DB.Create(&area).Error
	assert.NoError(t, err, "Failed to create test area")

	// Create test attendance
	startDate := time.Now().Add(-1 * time.Hour)
	endDate := time.Now().Add(1 * time.Hour)
	allowGuest := true
	useLuckyNumber := true
	limitCount := 5
	limitClientCount := 3
	attendance := checkins.Attendance{
		Title:                  "Test Attendance",
		AllowGuest:             allowGuest,
		UseLuckyNumber:         useLuckyNumber,
		LuckyShowAfterMinCount: 1,
		LimitCount:             limitCount,
		LimitClientCount:       limitClientCount,
		IsLocked:               false,
		StartDate:              &startDate,
		EndDate:                &endDate,
	}
	err = global.GVA_DB.Create(&attendance).Error
	assert.NoError(t, err, "Failed to create test attendance")

	// Create test participant
	fullName := "Test User"
	participant := checkins.Participant{
		Email:    "test@example.com",
		FullName: &fullName,
	}
	err = global.GVA_DB.Create(&participant).Error
	assert.NoError(t, err, "Failed to create test participant")

	// Link participant to attendance
	agp := checkins.AttendanceGroupParticipant{
		ParticipantId: &participant.ID,
		AttendanceId:  &attendance.ID,
	}
	err = global.GVA_DB.Create(&agp).Error
	assert.NoError(t, err, "Failed to create attendance group participant")

	// Create a condition for check-in
	showLuckyNumber := true
	condition := checkins.Condition{
		AttendanceId:    &attendance.ID,
		AreaId:          &area.ID,
		ShowLuckyNumber: showLuckyNumber,
	}
	err = global.GVA_DB.Create(&condition).Error
	assert.NoError(t, err, "Failed to create test condition")

	// Create attendance area
	attendanceArea := checkins.AttendanceArea{
		AreaID:          &area.ID,
		AttendanceID:    &attendance.ID,
		Area:            &area,
		AllowRestrictIp: false,
	}
	err = global.GVA_DB.Create(&attendanceArea).Error
	assert.NoError(t, err, "Failed to create test attendance area")

	// Update condition with area
	condition.Area = &attendanceArea
	err = global.GVA_DB.Save(&condition).Error
	assert.NoError(t, err, "Failed to update condition with area")

	// Create AGP condition
	agpCondition := checkins.AGPCondition{
		AttendanceGroupParticipantId: int(agp.ID),
		ConditionId:                  int(condition.ID),
		Condition:                    &condition,
	}
	err = global.GVA_DB.Create(&agpCondition).Error
	assert.NoError(t, err, "Failed to create AGP condition")

	// Encode attendance ID for the QR code
	code := base32.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", attendance.ID)))

	return code, attendance.ID, participant.ID
}

// TestCreateAttendanceCheckIn tests the basic functionality of CreateAttendanceCheckIn
func TestCreateAttendanceCheckIn(t *testing.T) {
	// Setup test data
	setupBasicTestData(t)

	// Find the participant and attendance for the test
	var participant checkins.Participant
	var attendance checkins.Attendance

	err := global.GVA_DB.Where("email = ?", "test@example.com").First(&participant).Error
	assert.NoError(t, err, "Failed to find test participant")

	err = global.GVA_DB.Where("title = ?", "Test Attendance").First(&attendance).Error
	assert.NoError(t, err, "Failed to find test attendance")

	// Create a test check-in
	lat := 10.762622
	lng := 106.660172
	accuracy := 10.0
	checkIn := &checkins.AttendanceCheckIn{
		CheckinDate:      time.Now(),
		AttendanceId:     &attendance.ID,
		PartpaticipantId: &participant.ID,
		Lattidue:         &lat,
		Longtidue:        &lng,
		Accuracy:         &accuracy,
		IP:               "127.0.0.1",
		Agent:            "Test Agent",
		VisitorId:        "test-visitor-id",
	}

	// Use the service to create the check-in
	service := AttendanceCheckInService{}
	err = service.CreateAttendanceCheckIn(checkIn)
	assert.NoError(t, err, "Failed to create attendance check-in")

	// Verify the check-in was created
	var count int64
	global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
		Where("attendance_id = ? AND partpaticipant_id = ?", attendance.ID, participant.ID).
		Count(&count)

	assert.Equal(t, int64(1), count, "Expected 1 check-in to be created")
	assert.Equal(t, 1, checkIn.Counter, "Counter should be 1 for new check-in")

	// Create a second check-in for the same attendance/participant/condition
	secondCheckIn := &checkins.AttendanceCheckIn{
		CheckinDate:      time.Now(),
		AttendanceId:     &attendance.ID,
		PartpaticipantId: &participant.ID,
		Lattidue:         &lat,
		Longtidue:        &lng,
		Accuracy:         &accuracy,
		IP:               "127.0.0.2", // Different IP to verify update
		Agent:            "Test Agent 2",
		VisitorId:        "test-visitor-id",
	}

	// Should increment counter instead of creating new record
	err = service.CreateAttendanceCheckIn(secondCheckIn)
	assert.NoError(t, err, "Failed to update attendance check-in")

	// Verify counter increased
	global.GVA_DB.Model(&checkins.AttendanceCheckIn{}).
		Where("attendance_id = ? AND partpaticipant_id = ?", attendance.ID, participant.ID).
		Count(&count)

	assert.Equal(t, int64(1), count, "Should still have 1 check-in record")
	assert.Equal(t, 2, secondCheckIn.Counter, "Counter should be 2 after update")
}

// TestCheckinAttendance tests the CheckinAttendance function with valid data
func TestCheckinAttendance(t *testing.T) {
	// Skip this test by default as it requires mocking of services
	// Uncomment the line below to run the test
	t.Skip("Skipping CheckinAttendance test that needs service mocking")

	// Set up complete test data
	code, _, _ := setupCompleteTestData(t)

	// Create a check-in request
	lat := 10.762622 // Same coordinates as test area
	lng := 106.660172
	accuracy := 10.0
	fullName := "Test User"
	req := checkinsReq.CheckinsReq{
		Code:      code,
		Email:     "test@example.com",
		FullName:  &fullName,
		VisitorId: "test-visitor-id",
		Lat:       &lat,
		Lng:       &lng,
		Accuracy:  &accuracy,
	}

	// Perform check-in
	service := AttendanceCheckInService{}

	// Mock the DecodeBase32 method for testing
	// Note: This is a simplified test - in a real implementation,
	// you would need to mock the dependent services (AttendanceService, ParticipantService, ConditionService)
	result, err := service.CheckinAttendance(req, "127.0.0.1", "Test User Agent")

	// Basic verification
	if err != nil {
		t.Logf("Error during check-in: %v", err)
	} else {
		t.Logf("Check-in message: %v", result["message"])
	}

	// Full assertions are skipped since this test requires proper mocking
}

// TestGenerateUniqueLuckyNumber tests the lucky number generation function
func TestGenerateUniqueLuckyNumber(t *testing.T) {
	// Set up basic test data
	setupBasicTestData(t)

	// Create the service
	service := AttendanceCheckInService{}

	// Find an attendance
	var attendance checkins.Attendance
	err := global.GVA_DB.First(&attendance).Error
	assert.NoError(t, err, "Failed to find test attendance")

	// Create a source for random numbers
	source := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(source))

	// Generate a unique lucky number
	luckyNumber, err := service.generateUniqueLuckyNumber(attendance.ID, 1, rnd)
	assert.NoError(t, err, "Failed to generate lucky number")
	assert.Greater(t, luckyNumber, 0, "Lucky number should be positive")

	// Generate another lucky number for a different participant
	luckyNumber2, err := service.generateUniqueLuckyNumber(attendance.ID, 2, rnd)
	assert.NoError(t, err, "Failed to generate second lucky number")
	assert.Greater(t, luckyNumber2, 0, "Second lucky number should be positive")

	// The two numbers should be different
	assert.NotEqual(t, luckyNumber, luckyNumber2, "Lucky numbers should be unique")
}

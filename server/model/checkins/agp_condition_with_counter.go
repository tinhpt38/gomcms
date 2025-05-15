package checkins

// AGPConditionWithCounter extends the AGPCondition struct with a Counter field for display purposes
type AGPConditionWithCounter struct {
	AGPCondition
	Counter int `json:"counter"`
}

// NewAGPConditionWithCounter creates a new AGPConditionWithCounter from a AGPCondition
func NewAGPConditionWithCounter(condition AGPCondition, counter int) AGPConditionWithCounter {
	return AGPConditionWithCounter{
		AGPCondition: condition,
		Counter:      counter,
	}
}

// ToFrontendMap converts the AGPCondition with counter to a map that matches the frontend's expected format
func (c AGPConditionWithCounter) ToFrontendMap() map[string]interface{} {
	result := map[string]interface{}{
		"ID":                           c.ID,
		"AttendanceId":                 c.AttendanceId,
		"AttendanceGroupParticipantId": c.AttendanceGroupParticipantId,
		"ConditionId":                  c.ConditionId,
		"counter":                      c.Counter,
	}

	// Add embedded condition if it exists
	if c.Condition != nil {
		// Create a ConditionWithCounter for the embedded condition
		condWithCounter := NewConditionWithCounter(*c.Condition, c.Counter)

		// Get the frontend-friendly map for the condition
		condMap := condWithCounter.ToFrontendMap()

		// Add the condition map to the result
		result["condition"] = condMap
	}

	return result
}

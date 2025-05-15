package checkins

// ConditionWithCounter extends the Condition struct with a Counter field for display purposes
type ConditionWithCounter struct {
	Condition
	Counter int `json:"counter"`
}

// NewConditionWithCounter creates a new ConditionWithCounter from a Condition
func NewConditionWithCounter(condition Condition, counter int) ConditionWithCounter {
	return ConditionWithCounter{
		Condition: condition,
		Counter:   counter,
	}
}

// ToFrontendMap converts the condition with counter to a map that matches the frontend's expected format
func (c ConditionWithCounter) ToFrontendMap() map[string]interface{} {
	result := map[string]interface{}{
		"ID":              c.ID,
		"AttendanceId":    c.AttendanceId,
		"GroupId":         c.GroupId,
		"AreaId":          c.AreaId,
		"StartAt":         c.StartAt,
		"EndAt":           c.EndAt,
		"IsPass":          c.IsPass,
		"counter":         c.Counter,
		"ShowLuckyNumber": c.ShowLuckyNumber,
		"Message":         c.Message,
	}

	// Add structure expected by frontend conditionString function
	if c.Group != nil {
		result["group"] = map[string]interface{}{
			"name": c.Group.Name,
		}
	}

	if c.Area != nil && c.Area.Area != nil {
		result["area"] = map[string]interface{}{
			"area": map[string]interface{}{
				"name": c.Area.Area.Name,
			},
		}
	}

	result["startAt"] = c.StartAt
	result["endAt"] = c.EndAt

	return result
}

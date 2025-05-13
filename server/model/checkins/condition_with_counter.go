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

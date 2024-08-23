package utils

import (
	"math"
	"strings"
)

func GetArrayValue(array []string, index int) string {
	if index >= len(array) {
		return ""
	}

	return strings.TrimSpace(array[index])
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

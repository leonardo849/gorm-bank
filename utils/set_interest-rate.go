package utils

import "time"

const (
	shortTermInterestRate = 3
	mediumTermInterestRate = 5
	longTermInterestRate = 7
)


func SetInterestRate(timeDifference time.Duration) float64 {
	day := 24 * time.Hour
	if timeDifference >= 30*day && timeDifference<90*day {
		return shortTermInterestRate
	} else if timeDifference >= 90*day && timeDifference < 365*day {
		return mediumTermInterestRate
	} else {
		return longTermInterestRate
	}
}
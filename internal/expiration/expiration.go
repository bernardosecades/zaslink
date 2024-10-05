package expiration

import (
	"fmt"
	"time"
)

type Expiration string

const (
	SevenDays     Expiration = "7d"
	ThreeDays     Expiration = "3d"
	OneDay        Expiration = "1d"
	TwelveHours   Expiration = "12h"
	FourHours     Expiration = "4h"
	OneHour       Expiration = "1h"
	ThirtyMinutes Expiration = "30m"
	FiveMinutes   Expiration = "5m"
)

func (e Expiration) Duration() time.Duration {
	switch e {
	case SevenDays:
		return time.Hour * 24 * 7
	case ThreeDays:
		return time.Hour * 24 * 3
	case OneDay:
		return time.Hour * 24
	case TwelveHours:
		return time.Hour * 12
	case FourHours:
		return time.Hour * 4
	case OneHour:
		return time.Hour * 1
	case ThirtyMinutes:
		return time.Minute * 30
	case FiveMinutes:
		return time.Minute * 5
	default:
		return 0
	}
}

// ValidateExpiration function to check if the provided string is a valid expiration
func ValidateExpiration(expStr string) (Expiration, error) {
	switch Expiration(expStr) {
	case SevenDays, ThreeDays, OneDay, TwelveHours, FourHours, OneHour, ThirtyMinutes, FiveMinutes:
		return Expiration(expStr), nil
	default:
		return "", fmt.Errorf("invalid expiration")
	}
}

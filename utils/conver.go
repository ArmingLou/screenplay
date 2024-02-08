package utils

import "fmt"

func SecondsToClockFormat(sec int) string {
	if sec < 60 {
		return fmt.Sprintf("00:%02d", sec)
	}
	if sec < 3600 {
		second := sec % 60
		minute := sec / 60
		return fmt.Sprintf("%02d:%02d", minute, second)
	}

	second := sec % 60
	min := sec / 60
	hour := min / 60
	minute := min % 60
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)

}

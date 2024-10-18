package validators

import "time"

func ValidateDate(date time.Time) bool {
	return date.After(time.Now())
}

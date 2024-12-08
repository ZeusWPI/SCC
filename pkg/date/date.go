// Package date makes working with dates without timezones easier
package date

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Date represents a date without a timezone
type Date time.Time

const dateLayout = "2006-01-02"

// UnmarshalJSON converts bytes to a Date
func (d *Date) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	if str == "" {
		*d = Date(time.Time{})
		return nil
	}

	// Parse the date
	parsedTime, err := time.Parse(dateLayout, str)
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}
	*d = Date(parsedTime)

	return nil
}

// ToTime converts a date to a time object
func (d Date) ToTime() time.Time {
	return time.Time(d)
}

// ValidateDate adds validation support for go-playground/validator for the Date type
func ValidateDate(f1 validator.FieldLevel) bool {
	date, ok := f1.Field().Interface().(Date)
	return ok && !date.ToTime().IsZero()
}

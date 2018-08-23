// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package date

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// A Date represents a date (year, month, day).
//
// This type does not include location information, and therefore does not
// describe a unique 24-hour timespan.
type Date struct {
	Year  int        // Year (e.g., 2014).
	Month time.Month // Month of the year (January = 1, ...).
	Day   int        // Day of the month, starting at 1.
}

// Parse parses a string in RFC3339 full-date format and returns the date value it represents.
func Parse(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return FromTime(t), nil
}

// FromTime returns the Date in which a time occurs in that time's location.
func FromTime(t time.Time) Date {
	var d Date
	d.Year, d.Month, d.Day = t.Date()
	return d
}

// String returns the date in RFC3339 full-date format.
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// IsValid reports whether the date is valid.
func (d Date) IsValid() bool {
	return FromTime(d.In(time.UTC)) == d
}

// In returns the time corresponding to time 00:00:00 of the date in the location.
//
// In is always consistent with time.Date, even when time.Date returns a time
// on a different day. For example, if loc is America/Indiana/Vincennes, then both
//     time.Date(1955, time.May, 1, 0, 0, 0, 0, loc)
// and
//     date.Date{Year: 1955, Month: time.May, Day: 1}.In(loc)
// return 23:00:00 on April 30, 1955.
//
// In panics if loc is nil.
func (d Date) In(loc *time.Location) time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, loc)
}

// AddDays returns the date that is n days in the future.
// n can also be negative to go into the past.
func (d Date) AddDays(n int) Date {
	return FromTime(d.In(time.UTC).AddDate(0, 0, n))
}

// DaysSince returns the signed number of days between the date and s, not including the end day.
// This is the inverse operation to AddDays.
func (d Date) DaysSince(s Date) (days int) {
	// We convert to Unix time so we do not have to worry about leap seconds:
	// Unix time increases by exactly 86400 seconds per day.
	deltaUnix := d.In(time.UTC).Unix() - s.In(time.UTC).Unix()
	return int(deltaUnix / 86400)
}

// Before reports whether d occurs before p.
func (d Date) Before(p Date) bool {
	if d.Year != p.Year {
		return d.Year < p.Year
	}
	if d.Month != p.Month {
		return d.Month < p.Month
	}
	return d.Day < p.Day
}

// After reports whether d1 occurs after p.
func (d Date) After(p Date) bool {
	return p.Before(d)
}

// Equal returns true if d equal d2.
func (d Date) Equal(d2 Date) bool {
	return d.Year == d2.Year && d.Month == d2.Month && d.Day == d2.Day
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is the result of d.String().
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected to be a string in a format accepted by ParseDate.
func (d *Date) UnmarshalText(data []byte) error {
	var err error
	*d, err = Parse(string(data))
	return err
}

// Value implements database/sql Valuer.
func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements the database/sql Scanner interface.
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date{}

		return nil
	}

	if tv, ok := value.(time.Time); ok {
		*d = FromTime(tv)

		return nil
	}

	if sv, ok := value.(string); ok {
		parsed, err := Parse(sv)
		if err != nil {
			return errors.New("failed to scan Date: " + err.Error())
		}

		*d = parsed

		return nil
	}

	return errors.New("failed to scan Date")
}

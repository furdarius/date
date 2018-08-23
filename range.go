package date

import (
	"sort"
)

// Range represents range between two dates.
type Range struct {
	Start Date `json:"start"`
	End   Date `json:"end"`
}

// IsValid reports if r is valid range.
func (r Range) IsValid() bool {
	return r.Start.IsValid() && r.End.IsValid()
}

// Empty returns true if Start equal End.
func (r Range) Empty() bool {
	return r.Start.Equal(r.End)
}

// Contains reports whether d is within r.
func (r Range) Contains(d Date) bool {
	return !d.Before(r.Start) && !d.After(r.End)
}

// Encloses returns true if the bounds of the inner range do not extend outside the bounds of the outer range.
func (r Range) Encloses(dr Range) bool {
	return r.Contains(dr.Start) && r.Contains(dr.End)
}

// RangeSet used to works with date ranges list.
type RangeSet []Range

// Sub returns RangeSet after dr subtraction.
func (a RangeSet) Sub(dr ...Range) RangeSet {
	busy := a.buildDictExcludeRanges(dr...)
	dates := linearizeBusy(busy)

	sort.Sort(ByAsc(dates))

	list := buildRanges(dates)

	return RangeSet(list)
}

// SubSet returns RangeSet list after set subtraction.
func (a RangeSet) SubSet(set RangeSet) RangeSet {
	return a.Sub(set.List()...)
}

// Impose returns RangeSet after dr imposition.
// Range intersections will be merged:
// [10, 15] impose with [13, 22] got [10, 22]
// [10, 15] impose with [20, 25] got [10, 15], [20, 25]
func (a RangeSet) Impose(dr ...Range) RangeSet {
	busy := a.buildDictIncludeRanges(dr...)
	dates := linearizeBusy(busy)

	sort.Sort(ByAsc(dates))

	list := buildRanges(dates)

	return RangeSet(list)
}

// ImposeSet returns ImposeSet list after set imposition.
func (a RangeSet) ImposeSet(set RangeSet) RangeSet {
	return a.Impose(set.List()...)
}

// TrimEnd shifts end date a day earlier.
func (a RangeSet) TrimEnd() RangeSet {
	return a.ShiftEnd(-1)
}

// ExtendEnd shifts end date a day later.
func (a RangeSet) ExtendEnd() RangeSet {
	return a.ShiftEnd(1)
}

// ShiftEnd shifts end date to n days.
// n can also be negative to go into the past.
func (a RangeSet) ShiftEnd(n int) RangeSet {
	for i := 0; i < len(a); i++ {
		a[i].End = a[i].End.AddDays(n)
	}

	return a
}

// RangeTest used to filter RangeSet.
type RangeTest func(Range) bool

// RangeNotEmpty tests that Range is not empty.
func RangeNotEmpty(dr Range) bool {
	return !dr.Empty()
}

// Filter returns RangeSet with all elements that pass the test.
func (a RangeSet) Filter(test RangeTest) RangeSet {
	res := RangeSet{}
	for i := 0; i < len(a); i++ {
		if test(a[i]) {
			res = append(res, a[i])
		}
	}

	return res
}

// FilterEmpty returns RangeSet with not empty Ranges from a.
func (a RangeSet) FilterEmpty() RangeSet {
	return a.Filter(RangeNotEmpty)
}

// List returns list of Ranges in set.
func (a RangeSet) List() []Range {
	return []Range(a)
}

func (a RangeSet) buildBusyDict() map[Date]bool {
	busy := map[Date]bool{}

	for _, dr := range a {
		days := dr.End.DaysSince(dr.Start)
		for i := 0; i <= days; i++ {
			busy[dr.Start.AddDays(i)] = true
		}
	}

	return busy
}

func (a RangeSet) buildDictExcludeRanges(drs ...Range) map[Date]bool {
	busy := a.buildBusyDict()

	for _, dr := range drs {
		days := dr.End.DaysSince(dr.Start)
		for i := 0; i <= days; i++ {
			busy[dr.Start.AddDays(i)] = false
		}
	}

	return busy
}

func (a RangeSet) buildDictIncludeRanges(drs ...Range) map[Date]bool {
	busy := a.buildBusyDict()

	for _, dr := range drs {
		days := dr.End.DaysSince(dr.Start)
		for i := 0; i <= days; i++ {
			busy[dr.Start.AddDays(i)] = true
		}
	}

	return busy
}

func linearizeBusy(busy map[Date]bool) []Date {
	res := []Date{}
	for d, busy := range busy {
		if busy {
			res = append(res, d)
		}
	}

	return res
}

func buildRanges(dates []Date) []Range {
	N := len(dates)

	if N == 0 {
		return []Range{}
	}

	ranges := []Range{}
	start := dates[0]

	for i := 1; i < N; i++ {
		if dates[i].DaysSince(dates[i-1]) == 1 {
			continue
		}

		ranges = append(ranges, Range{start, dates[i-1]})

		start = dates[i]
	}

	ranges = append(ranges, Range{start, dates[N-1]})

	return ranges
}

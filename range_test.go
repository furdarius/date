package date

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntervalJSONMarshaling(t *testing.T) {
	var dr Range
	for _, test := range []struct {
		data string
		ptr  interface{}
		want interface{}
	}{
		{
			`{"start": "2018-07-15","end": "2018-07-20"}`, &dr,
			&Range{
				Start: Date{Year: 2018, Month: 07, Day: 15},
				End:   Date{Year: 2018, Month: 07, Day: 20},
			},
		},
		{
			`{"end": "2040-11-01","start": "2025-07-15"}`, &dr,
			&Range{
				Start: Date{Year: 2025, Month: 07, Day: 15},
				End:   Date{Year: 2040, Month: 11, Day: 01},
			},
		},
	} {
		err := json.Unmarshal([]byte(test.data), test.ptr)
		assert.NoError(t, err)
		assert.Equal(t, test.ptr, test.want)
	}
}

func TestBuildRanges(t *testing.T) {
	for _, test := range []struct {
		dates    []Date
		expected []Range
	}{
		{
			dates: []Date{{2018, 10, 15}, {2018, 10, 16}, {2018, 10, 17}, {2018, 10, 20}, {2018, 10, 21}},
			expected: []Range{
				{Date{2018, 10, 15}, Date{2018, 10, 17}},
				{Date{2018, 10, 20}, Date{2018, 10, 21}},
			},
		},
		{
			dates: []Date{{2018, 10, 7}, {2018, 10, 11}, {2018, 10, 12}, {2018, 10, 13}, {2018, 10, 14}, {2018, 10, 15},
				{2018, 10, 19}, {2018, 10, 20}, {2018, 10, 21}, {2018, 10, 24}, {2018, 10, 25}},
			expected: []Range{
				{Date{2018, 10, 7}, Date{2018, 10, 7}},
				{Date{2018, 10, 11}, Date{2018, 10, 15}},
				{Date{2018, 10, 19}, Date{2018, 10, 21}},
				{Date{2018, 10, 24}, Date{2018, 10, 25}},
			},
		},
		{
			dates: []Date{{2018, 10, 7}},
			expected: []Range{
				{Date{2018, 10, 7}, Date{2018, 10, 7}},
			},
		},
		{
			dates: []Date{{2018, 10, 7}, {2018, 10, 8}, {2018, 10, 20}, {2018, 10, 21}, {2018, 10, 25}},
			expected: []Range{
				{Date{2018, 10, 7}, Date{2018, 10, 8}},
				{Date{2018, 10, 20}, Date{2018, 10, 21}},
				{Date{2018, 10, 25}, Date{2018, 10, 25}},
			},
		},
		{
			dates: []Date{{2018, 10, 7}, {2018, 10, 8}, {2018, 10, 9}, {2018, 10, 10}, {2018, 10, 11}},
			expected: []Range{
				{Date{2018, 10, 7}, Date{2018, 10, 11}},
			},
		},
		{
			dates: []Date{{2018, 10, 5}, {2018, 10, 8}, {2018, 10, 10}, {2018, 10, 22}, {2018, 10, 28}},
			expected: []Range{
				{Date{2018, 10, 5}, Date{2018, 10, 5}},
				{Date{2018, 10, 8}, Date{2018, 10, 8}},
				{Date{2018, 10, 10}, Date{2018, 10, 10}},
				{Date{2018, 10, 22}, Date{2018, 10, 22}},
				{Date{2018, 10, 28}, Date{2018, 10, 28}},
			},
		},
	} {
		actual := buildRanges(test.dates)
		assert.Equal(t, test.expected, actual)
	}
}

func TestRangesIntersections(t *testing.T) {
	for _, test := range []struct {
		r1, r2   Range
		expected bool
	}{
		{
			r1: Range{
				Start: Date{Year: 2018, Month: 10, Day: 15},
				End:   Date{Year: 2018, Month: 10, Day: 25},
			},
			r2: Range{
				Start: Date{Year: 2018, Month: 10, Day: 15},
				End:   Date{Year: 2018, Month: 10, Day: 25},
			},
			expected: true,
		},
		{
			r1: Range{
				Start: Date{Year: 2018, Month: 10, Day: 15},
				End:   Date{Year: 2018, Month: 10, Day: 25},
			},
			r2: Range{
				Start: Date{Year: 2018, Month: 10, Day: 18},
				End:   Date{Year: 2018, Month: 10, Day: 28},
			},
			expected: true,
		},
		{
			r1: Range{
				Start: Date{Year: 2018, Month: 10, Day: 15},
				End:   Date{Year: 2018, Month: 10, Day: 25},
			},
			r2: Range{
				Start: Date{Year: 2018, Month: 10, Day: 18},
				End:   Date{Year: 2018, Month: 10, Day: 28},
			},
			expected: true,
		},
		{
			r1: Range{
				Start: Date{Year: 2018, Month: 10, Day: 15},
				End:   Date{Year: 2018, Month: 10, Day: 25},
			},
			r2: Range{
				Start: Date{Year: 2018, Month: 10, Day: 25},
				End:   Date{Year: 2018, Month: 10, Day: 28},
			},
			expected: true,
		},
		{
			r1: Range{
				Start: Date{Year: 2018, Month: 10, Day: 15},
				End:   Date{Year: 2018, Month: 10, Day: 24},
			},
			r2: Range{
				Start: Date{Year: 2018, Month: 10, Day: 25},
				End:   Date{Year: 2018, Month: 10, Day: 28},
			},
			expected: false,
		},
		{
			r1: Range{
				Start: Date{Year: 2018, Month: 10, Day: 15},
				End:   Date{Year: 2018, Month: 10, Day: 24},
			},
			r2: Range{
				Start: Date{Year: 2018, Month: 10, Day: 16},
				End:   Date{Year: 2018, Month: 10, Day: 18},
			},
			expected: true,
		},
		{
			r1: Range{
				Start: Date{Year: 2018, Month: 10, Day: 15},
				End:   Date{Year: 2018, Month: 10, Day: 24},
			},
			r2: Range{
				Start: Date{Year: 2018, Month: 10, Day: 10},
				End:   Date{Year: 2018, Month: 10, Day: 29},
			},
			expected: true,
		},
	} {
		actual := test.r1.Intersects(test.r2)
		assert.Equalf(t, test.expected, actual, "%v intersection with %v is %t, but %t expected",
			test.r1, test.r2, actual, test.expected)
	}
}

func TestRangeSet_Sub(t *testing.T) {
	for _, test := range []struct {
		basement []Range
		sub      []Range
		expected []Range
	}{
		{
			basement: []Range{
				{Date{2018, 7, 11}, Date{2018, 7, 14}},
				{Date{2018, 7, 7}, Date{2018, 7, 8}},
				{Date{2018, 7, 27}, Date{2018, 7, 28}},
				{Date{2018, 7, 23}, Date{2018, 7, 23}},
			},
			sub: []Range{
				{Date{2018, 7, 15}, Date{2018, 7, 24}},
			},
			expected: []Range{
				{Date{2018, 7, 7}, Date{2018, 7, 8}},
				{Date{2018, 7, 11}, Date{2018, 7, 14}},
				{Date{2018, 7, 27}, Date{2018, 7, 28}},
			},
		},
		{
			basement: []Range{
				{Date{2018, 2, 11}, Date{2018, 2, 15}},
				{Date{2018, 2, 19}, Date{2018, 2, 21}},
				{Date{2018, 2, 7}, Date{2018, 2, 7}},
				{Date{2018, 2, 24}, Date{2018, 2, 25}},
			},
			sub: []Range{
				{Date{2018, 2, 14}, Date{2018, 2, 22}},
			},
			expected: []Range{
				{Date{2018, 2, 7}, Date{2018, 2, 7}},
				{Date{2018, 2, 11}, Date{2018, 2, 13}},
				{Date{2018, 2, 24}, Date{2018, 2, 25}},
			},
		},
		{
			basement: []Range{
				{Date{2018, 2, 11}, Date{2018, 2, 15}},
				{Date{2018, 2, 19}, Date{2018, 2, 21}},
				{Date{2018, 2, 7}, Date{2018, 2, 7}},
				{Date{2018, 2, 24}, Date{2018, 2, 25}},
			},
			sub: []Range{
				{Date{2018, 2, 14}, Date{2018, 2, 14}},
			},
			expected: []Range{
				{Date{2018, 2, 7}, Date{2018, 2, 7}},
				{Date{2018, 2, 11}, Date{2018, 2, 13}},
				{Date{2018, 2, 15}, Date{2018, 2, 15}},
				{Date{2018, 2, 19}, Date{2018, 2, 21}},
				{Date{2018, 2, 24}, Date{2018, 2, 25}},
			},
		},
		{
			basement: []Range{
				{Date{2018, 2, 11}, Date{2018, 2, 15}},
				{Date{2018, 2, 19}, Date{2018, 2, 21}},
				{Date{2018, 2, 7}, Date{2018, 2, 7}},
				{Date{2018, 2, 24}, Date{2018, 2, 25}},
			},
			sub: []Range{
				{Date{2018, 1, 14}, Date{2018, 3, 14}},
			},
			expected: []Range{},
		},
		{
			basement: []Range{
				{Date{2018, 2, 11}, Date{2018, 2, 15}},
				{Date{2018, 2, 19}, Date{2018, 2, 21}},
				{Date{2018, 2, 7}, Date{2018, 2, 7}},
				{Date{2018, 2, 24}, Date{2018, 2, 25}},
			},
			sub: []Range{{}},
			expected: []Range{
				{Date{2018, 2, 7}, Date{2018, 2, 7}},
				{Date{2018, 2, 11}, Date{2018, 2, 15}},
				{Date{2018, 2, 19}, Date{2018, 2, 21}},
				{Date{2018, 2, 24}, Date{2018, 2, 25}},
			},
		},
		{
			basement: []Range{
				{Date{2018, 2, 11}, Date{2018, 2, 15}},
			},
			sub: []Range{
				{Date{2018, 2, 2}, Date{2018, 2, 10}},
				{Date{2018, 2, 11}, Date{2018, 2, 13}},
				{Date{2018, 2, 18}, Date{2018, 2, 22}},
			},
			expected: []Range{
				{Date{2018, 2, 14}, Date{2018, 2, 15}},
			},
		},
		{
			basement: []Range{
				{Date{2018, 2, 11}, Date{2018, 2, 15}},
			},
			sub: []Range{
				{Date{2018, 2, 2}, Date{2018, 2, 11}},
				{Date{2018, 2, 13}, Date{2018, 2, 13}},
				{Date{2018, 2, 15}, Date{2018, 2, 18}},
			},
			expected: []Range{
				{Date{2018, 2, 12}, Date{2018, 2, 12}},
				{Date{2018, 2, 14}, Date{2018, 2, 14}},
			},
		},
		{
			basement: []Range{
				{Date{2018, 7, 11}, Date{2018, 2, 13}},
			},
			sub: []Range{
				{Date{2018, 7, 10}, Date{2018, 7, 15}},
			},
			expected: []Range{},
		},
	} {
		actual := RangeSet(test.basement).Sub(test.sub...).List()
		assert.Equal(t, test.expected, actual)
	}
}

func TestRangeSet_Filter(t *testing.T) {
	for _, test := range []struct {
		in       RangeSet
		test     RangeTest
		expected RangeSet
	}{
		{
			in: RangeSet([]Range{
				{},
				{Date{2018, 2, 12}, Date{2018, 2, 12}},
				{Date{2018, 2, 14}, Date{2018, 2, 14}},
				{},
			}),
			test:     RangeNotEmpty,
			expected: RangeSet([]Range{}),
		},
		{
			in: RangeSet([]Range{
				{},
				{Date{2018, 2, 12}, Date{2018, 2, 12}},
				{Date{2018, 2, 14}, Date{2018, 2, 22}},
				{},
			}),
			test: RangeNotEmpty,
			expected: RangeSet([]Range{
				{Date{2018, 2, 14}, Date{2018, 2, 22}},
			}),
		},
		{
			in:       RangeSet([]Range{}),
			test:     RangeNotEmpty,
			expected: RangeSet([]Range{}),
		},
		{
			in: RangeSet([]Range{
				{},
				{Date{2018, 2, 22}, Date{2018, 2, 25}},
				{Date{2018, 2, 28}, Date{2018, 3, 22}},
				{},
			}),
			test: func(dr Range) bool { return true },
			expected: RangeSet([]Range{
				{},
				{Date{2018, 2, 22}, Date{2018, 2, 25}},
				{Date{2018, 2, 28}, Date{2018, 3, 22}},
				{},
			}),
		},
		{
			in: RangeSet([]Range{
				{},
				{Date{2018, 2, 22}, Date{2018, 2, 25}},
				{Date{2018, 2, 28}, Date{2018, 3, 22}},
				{},
			}),
			test:     func(dr Range) bool { return false },
			expected: RangeSet([]Range{}),
		},
	} {
		actual := test.in.Filter(test.test)
		assert.Equal(t, test.expected, actual)
	}
}

func TestRangeSet_Impose(t *testing.T) {
	for _, test := range []struct {
		basement RangeSet
		impose   []Range
		expected RangeSet
	}{
		{
			basement: []Range{
				{Date{2018, 2, 15}, Date{2018, 2, 16}},
				{Date{2018, 2, 20}, Date{2018, 2, 21}},
			},
			impose: []Range{
				{Date{2018, 2, 9}, Date{2018, 2, 19}},
			},
			expected: []Range{
				{Date{2018, 2, 9}, Date{2018, 2, 21}},
			},
		},
		{
			basement: []Range{
				{Date{2018, 2, 2}, Date{2018, 2, 5}},
				{Date{2018, 2, 15}, Date{2018, 2, 21}},
				{Date{2018, 3, 10}, Date{2018, 4, 17}},
			},
			impose: []Range{
				{Date{2018, 2, 9}, Date{2018, 2, 19}},
			},
			expected: []Range{
				{Date{2018, 2, 2}, Date{2018, 2, 5}},
				{Date{2018, 2, 9}, Date{2018, 2, 21}},
				{Date{2018, 3, 10}, Date{2018, 4, 17}},
			},
		},
	} {
		actual := test.basement.Impose(test.impose...)
		assert.Equal(t, test.expected, actual)
	}
}

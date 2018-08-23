package date

import (
	"testing"
	"sort"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestByStartSorting(t *testing.T) {
	dates := []Date{{2018, 07, 12}, {2018, 05, 15}, {2018, 07, 22}, {2018, 10, 13}, {2018, 05, 05}}

	sort.Sort(ByAsc(dates))

	for i := 0; i < len(dates)-1; i++ {
		inOrder := dates[i].In(time.UTC).Before(dates[i+1].In(time.UTC))
		assert.True(t, inOrder)
	}
}

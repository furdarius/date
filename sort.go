package date

// ByAsc sorting type - used to sort Date slice from the earliest to the latest
type ByAsc []Date

// Len returns length of underlying array.
func (a ByAsc) Len() int { return len(a) }

// Swap swaps i and j elements.
func (a ByAsc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less compare i and j elements and returns true if element i less than element j.
func (a ByAsc) Less(i, j int) bool { return a[i].Before(a[j]) }

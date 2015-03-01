package size

type AllocFunc func(int, []T) []int

func AllocMax(value int, sizes []T) []int {
	subvals := make([]int, len(sizes))
	for i, size := range sizes {
		switch t := size.(type) {
		case ConstT:
			subvals[i] = lower(int(t), value)
		case RangeT:
			subvals[i] = lower(t.max, value)
		default:
			subvals[i] = value
		}
	}
	return subvals
}

// range.max allocations:
// 1. fair
// 2. prioritize rightmost
// 3. prioritize leftmost
// 4. 1 then 2
// 5. 1 then 3

// end allocations:
// 1. fair
// 2. leftmost
// 3. rightmost

// range.max and end allocations:
// r4, e1

func AllocFair(value int, sizes []T) []int {
	subvals := make([]int, len(sizes))

	deduct := func(i, x int) {
		y := lower(value, x)
		subvals[i] = y
		value = zero(value - y)
	}

	indices := []int{}
	x := 0

	// TODO: reduce number of iterations

	// allocate for const and min of range
	for i, size := range sizes {
		switch t := size.(type) {
		case ConstT:
			deduct(i, int(t))
		case RangeT:
			deduct(i, t.min)
			x += t.max - t.min
			indices = append(indices, i)
		case FreeT:
			indices = append(indices, i)
		}
	}

	remc := higher(1, len(indices)) // avoid div by zero
	rem := value % remc
	n := value / remc

	for _, i := range indices {
		switch t := sizes[i].(type) {
		case FreeT:
			subvals[i] += n
		case RangeT:
			m := subvals[i] + n
			if m > t.max {
				subvals[i] = t.max
				rem += m - t.max
			} else {
				subvals[i] += n
			}
		}
	}
	value = rem

	for i := len(indices) - 1; i >= 0 && value > 0; i-- {
		switch t := sizes[i].(type) {
		case FreeT:
			subvals[i] += value
			value = 0
		case RangeT:
			m := subvals[i] + value
			if m > t.max {
				subvals[i] = t.max
				value = m - t.max
			} else {
				subvals[i] += value
				value = 0
			}
		}
	}

	return subvals
}

package size

func lower(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func higher(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func zero(x int) int {
	if x >= 0 {
		return x
	}
	return 0
}

package sum

// Returns sum of a list of integers
func Ints(s ...int) int {
	return Sum(s)
}

func Sum(s []int) int {
	if len(s) == 0 {
		return 0
	}

	return Sum(s[1:]) + s[0]
}

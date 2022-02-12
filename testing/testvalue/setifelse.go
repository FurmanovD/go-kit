package testvalue

func SetIfElseString(cond bool, valTrue, valFalse string) string {
	if cond {
		return valTrue
	}
	return valFalse
}

func SetIfElseInt(cond bool, valTrue, valFalse int) int {
	if cond {
		return valTrue
	}
	return valFalse
}

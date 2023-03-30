package testvalue

// IfElse returns the value of the valTrue if the condition is true and returns valFalse otherwise.
func IfElse[T comparable](cond bool, valTrue, valFalse T) T {
	if cond {
		return valTrue
	}
	return valFalse
}

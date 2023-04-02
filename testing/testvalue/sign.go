package testvalue

type NumberPosNeg interface {
	int | int32 | int64 | float32 | float64
}

// Pos returns abs(n) of the number !=0 or 1 otherwise.
func Pos[T NumberPosNeg](n T) T {
	if n == 0 {
		n = 1
	} else if n < 0 {
		n *= -1
	}
	return n
}

// Neg returns -abs(n) of the number !=0 or -1 otherwise.
func Neg[T NumberPosNeg](n T) T {
	// nolint:typecheck
	return -Pos(n)
}

// NonPos returns -abs(n) or zero.
func NonPos[T NumberPosNeg](n T) T {
	if n > 0 {
		n *= -1
	}
	return n
}

// NonNeg returns abs(n) or zero.
func NonNeg[T NumberPosNeg](n T) T {
	return -NonPos(n)
}

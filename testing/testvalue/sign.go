package testvalue

type NumberPosNeg interface {
	int | int32 | int64 | float32 | float64
}

func Pos[T NumberPosNeg](n T) T {
	if n == 0 {
		n = 1
	} else if n < 0 {
		n *= -1
	}
	return n
}

func Neg[T NumberPosNeg](n T) T {
	return -Pos(n)
}

func NonPos[T NumberPosNeg](n T) T {
	if n > 0 {
		n *= -1
	}
	return n
}

func NonNeg[T NumberPosNeg](n T) T {
	return -NonPos(n)
}

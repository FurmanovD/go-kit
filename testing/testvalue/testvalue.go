package testvalue

import (
	"math/rand"
	"strings"

	"github.com/FurmanovD/go-kit/randomstring"
	"github.com/ericlagergren/decimal"
)

func NonEmptyNaturalIntSlice(nums []int) []int {
	if len(nums) == 0 {
		return []int{3}
	}
	for i, n := range nums {
		// nolint:typecheck
		nums[i] = Pos(n)
	}
	return nums
}

func RandItemStr(items ...string) string {
	if len(items) == 0 {
		return ""
	}

	return items[rand.Intn(len(items))]
}

func RandItemInt(items ...int) int {
	if len(items) == 0 {
		return 0
	}

	return items[rand.Intn(len(items))]
}

func RandInt8Flag() int8 {
	var val int8
	if rand.Intn(100) >= 50 {
		val = 1
	}
	return val
}

func RandDecimalBig(allowNegative bool, beforeCommaDigits int, afterCommaDigits int) decimal.Big {
	var before int
	if beforeCommaDigits > 0 {
		before = beforeCommaDigits
	}

	var after int
	if afterCommaDigits > 0 {
		after = afterCommaDigits
	}

	var numString strings.Builder
	if allowNegative && rand.Intn(100) >= 50 {
		numString.WriteRune('-')
	}

	if before > 0 {
		// non-zero char first
		numString.WriteString(randomstring.FromSet(randomstring.Decimal[1:], 1, nil))
		if beforeCommaDigits > 1 {
			numString.WriteString(randomstring.FromSet(randomstring.Decimal, beforeCommaDigits-1, nil))
		}
	} else {
		numString.WriteRune('0')
	}

	if after > 0 {
		numString.WriteRune('.')
		numString.WriteString(randomstring.FromSet(randomstring.Decimal, afterCommaDigits, nil))
	}

	var val decimal.Big
	val.SetString(numString.String())

	return val
}

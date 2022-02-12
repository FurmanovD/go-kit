package testvalue

import (
	"math/rand"
	"time"

	"github.com/FurmanovD/go-kit/randomstring"
	"github.com/ericlagergren/decimal"
)

func NotNilRnd(r *rand.Rand) *rand.Rand {
	if r == nil {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return r
}

func NaturalInt(n int) int {
	if n <= 0 {
		n = (n * -1) + 1 // make a natural number
	}
	return n
}

func NonEmptyNaturalIntSlice(nums []int) []int {
	if len(nums) == 0 {
		return []int{3}
	}
	for i, n := range nums {
		nums[i] = NaturalInt(n)
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

	before := beforeCommaDigits
	if before < 0 {
		before = 0
	}
	after := afterCommaDigits
	if after < 0 {
		after = 0
	}

	numString := ""
	if allowNegative && rand.Intn(100) >= 50 {
		numString = "-"
	}

	if beforeCommaDigits > 0 {
		numString = randomstring.FromSet(randomstring.Decimal[1:], 1, nil)
		if beforeCommaDigits > 1 {
			numString += randomstring.FromSet(randomstring.Decimal, beforeCommaDigits-1, nil)
		}
	} else {
		numString = "0"
	}

	if afterCommaDigits > 0 {
		numString += "." + randomstring.FromSet(randomstring.Decimal, afterCommaDigits, nil)
	}

	var val decimal.Big
	val.SetString(numString)

	return val
}

package testvalue

import (
	"math/rand"
	"time"

	"github.com/ericlagergren/decimal"
)

func NilOrTime(nilPercent int, isPast bool) *time.Time {

	percent := nilPercent
	if percent <= 0 || percent >= 100 {
		percent = 50
	}

	if percent < 1+rand.Intn(100) {
		return nil
	}

	now := time.Now().UTC()
	timeDiff := time.Duration(rand.Intn(365*24)) * time.Hour
	if isPast {
		now.Add(-timeDiff)
	} else {
		now.Add(timeDiff)
	}
	return &now
}

func NilOrBool(nilPercent int) *bool {

	percent := nilPercent
	if percent <= 0 || percent >= 100 {
		percent = 50
	}

	if percent < 1+rand.Intn(100) {
		return nil
	}

	val := rand.Intn(10) <= 4
	return &val
}

func NilOrInt8Flag(nilPercent int) *int8 {

	percent := nilPercent
	if percent <= 0 || percent >= 100 {
		percent = 50
	}

	if percent < 1+rand.Intn(100) {
		return nil
	}

	val := RandInt8Flag()
	return &val
}

func NilOrFloat32(nilPercent int, maxValBeforeComma uint) *float32 {

	percent := nilPercent
	if percent <= 0 || percent >= 100 {
		percent = 50
	}

	if percent < 1+rand.Intn(100) {
		return nil
	}

	val := rand.Float32() + float32(rand.Intn(int(maxValBeforeComma)))
	return &val
}

func NilOrIntn(nilPercent int, maxValue int) *int {

	percent := nilPercent
	if percent <= 0 || percent >= 100 {
		percent = 50
	}

	var val int
	if maxValue <= 0 {
		return &val
	}

	if percent < 1+rand.Intn(100) {
		return nil
	}

	val = rand.Intn(maxValue)
	return &val
}

func NilOrInt8(allowNegative bool, nilPercent int, maxValue int) *int8 {

	percent := nilPercent
	if percent <= 0 || percent >= 100 {
		percent = 50
	}

	var val int8
	if maxValue <= 0 {
		return &val
	}

	if percent < 1+rand.Intn(100) {
		return nil
	}

	val = int8(rand.Intn(maxValue))
	return &val
}

func NilOrUInt8(allowNegative bool, nilPercent int, maxValue int) *uint8 {

	int8Value := NilOrInt8(allowNegative, nilPercent, maxValue)
	if int8Value == nil {
		return nil
	}
	val := uint8(*int8Value)
	return &val
}

func NilOrDecimalBig(nilPercent int, allowNegative bool, beforeCommaDigits int, afterCommaDigits int) *decimal.Big {

	percent := nilPercent
	if percent <= 0 || percent >= 100 {
		percent = 50
	}

	if percent < 1+rand.Intn(100) {
		return nil
	}

	val := RandDecimalBig(allowNegative, beforeCommaDigits, afterCommaDigits)
	return &val
}

// NilOrStr returns a nil or pointer to a given string value.
// Used to construct nilPercent of invalid null.String values using FromStringPtr() method.
// e.g.:
// null.StringFromPtr(testvalue.NilOrStr(10, randomstring.UTF8Printable(20, rnd)))
// will return 10% of NULL DB values and 90% of random UTF strings
// null.StringFromPtr(testvalue.NilOrStr(10, testvalue.RandItemStr("ABC", "DEF")))
// will fill out a NULLABLE DB field of enum: "ABC", "DEF"
func NilOrStr(nilPercent int, str string) *string {

	percent := nilPercent
	if percent <= 0 || percent >= 100 {
		percent = 50
	}

	if percent < 1+rand.Intn(100) {
		return nil
	}
	val := str
	return &val
}

func NilOrRandStrItem(nilPercent int, items ...string) *string {

	percent := nilPercent
	if percent <= 0 || percent >= 100 {
		percent = 50
	}

	if percent < 1+rand.Intn(100) {
		return nil
	}
	val := items[rand.Intn(len(items))]
	return &val
}

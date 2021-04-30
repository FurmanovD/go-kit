package randomstring

import (
	"errors"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

const (
	LatinUpper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LatinLower   = "abcdefghijklmnopqrstuvwxyz"
	Latin        = LatinUpper + LatinLower
	Decimal      = "0123456789"
	Hexadecimal  = Decimal + "ABCDEF"
	Alphanumeric = Latin + Decimal
	Punctuation  = ",.:;!?-"
	Bracket      = "{}[]()"
	SpecCharsKbd = "`~!@#$%^&*()|"
)

var (
	totalUTF8Printable int = totalInRange(unicode.GraphicRanges)
)

// FromSet returns a random string built of runes from a set provided
func FromSet(set string, length int, rnd *rand.Rand) string {
	if len(set) == 0 {
		return ""
	}

	r := notNilRand(rnd)

	str := make([]rune, length)
	for i := 0; i < length; i++ {
		str[i] = rune(set[r.Intn(len(set)-1)])
	}
	return string(str)
}

// URL returns a random URL string with given schema, 1-4 random subdomains and 0-5 subdirectories
func URL(schema string, rnd *rand.Rand) string {
	r := notNilRand(rnd)

	str := schema + "://"

	// subdomains
	// 1-4 subdomain
	// of names [1 latin lower]+[2-12 Latin or Decimal or -_].
	for i := 1 + r.Intn(3); i > 0; i-- {
		str += FromSet(LatinLower, 1, r) +
			FromSet(
				Latin+Decimal+"-_",
				2+r.Intn(10),
				r,
			) + "."
	}
	str += "net"

	// directories
	for i := r.Intn(5); i > 0; i-- {
		str += FromSet(
			Latin+Decimal+"-_",
			2+r.Intn(10),
			r,
		) + "/"
	}
	return str
}

// NonEmptyUTF8Printable returns a random string of the length between 1 and maxlength
// comprised of printable UTF8 characters
func NonEmptyUTF8Printable(maxlength int, rnd *rand.Rand) string {

	r := notNilRand(rnd)
	if maxlength <= 1 {
		return UTF8Printable(1, r)
	}

	len := uint64(1 + r.Int63n(int64(maxlength-1)))
	return UTF8Printable(len, r)
}

// UTF8Printable generates a random UTF8 string of a given character length
// consist of the printable UTF8 characters
func UTF8Printable(length uint64, rnd *rand.Rand) string {
	if length == 0 {
		return ""
	}

	r := notNilRand(rnd)

	s := make([]rune, length)
	var i uint64
	for i = 0; i < length; i++ {
		s[i] = RandomRuneUTF8Printable(r)
	}

	return string(s)
}

// RandomRuneUTF8Printable returns a random rune of a printable range.
func RandomRuneUTF8Printable(rnd *rand.Rand) rune {
	x, _ := getItemFromRangeTable(
		notNilRand(rnd).Intn(totalUTF8Printable-1),
		unicode.GraphicRanges,
	)
	return rune(x)
}

// UTF8FromRanges generates a random UTF8 string of a given character length that exists on the
// given RangeTables. For a list of valid RangeTables, see
// http://golang.org/pkg/unicode/#pkg-variables
func UTF8FromRanges(length int, rnd *rand.Rand, tables ...*unicode.RangeTable) string {
	r := notNilRand(rnd)

	s := make([]rune, length)
	for i := 0; i < length; i++ {
		s[i] = RandomRuneUTF8FromRanges(tables, r)
	}

	return string(s)
}

// RandomRuneUTF8FromRanges returns a random rune in the given RangeTable.
func RandomRuneUTF8FromRanges(tables []*unicode.RangeTable, rnd *rand.Rand) rune {
	x, _ := getItemFromRangeTable(
		notNilRand(rnd).Intn(totalInRange(tables)),
		tables,
	)
	return rune(x)
}

// ToRandomCase returns randomly cased(uppercase, lowercase or RaNDomcASe) input string or the input string without any changes
func ToRandomCase(input string, rnd *rand.Rand) string {
	if input == "" {
		return ""
	}

	r := notNilRand(rnd)
	switch r.Intn(3) {
	case 0:
		return strings.ToLower(input)
	case 1:
		return strings.ToUpper(input)
	case 2:
		output := make([]rune, len([]rune(input)))
		for i, char := range input {
			// apply or not casing to the next rune
			switch r.Intn(2) {
			case 0:
				output[i] = unicode.ToLower(char)
			case 1:
				output[i] = unicode.ToUpper(char)
			default:
				output[i] = char
			}
		}
		return string(output)

	default:
		return input
	}
}

/// Returns the nth item contained in the array of ranges.
func getItemFromRangeTable(n int, tables []*unicode.RangeTable) (int, error) {
	nPointer := n
	var picked int
	found := false

	for _, table := range tables {
		if !found {
			for _, r16 := range table.R16 {
				countInRange := int((r16.Hi-r16.Lo)/r16.Stride) + 1
				if nPointer <= countInRange-1 {
					picked = int(r16.Lo) + (int(r16.Stride) * nPointer)
					found = true
					break
				} else {
					nPointer -= countInRange
				}
			}

			if !found {
				for _, r32 := range table.R32 {
					countInRange := int((r32.Hi-r32.Lo)/r32.Stride) + 1
					if nPointer <= countInRange-1 {
						picked = int(r32.Lo) + (int(r32.Stride) * nPointer)
						found = true
						break
					} else {
						nPointer -= countInRange
					}
				}
			}
		}
	}

	if found {
		return picked, nil
	} else {
		return -1, errors.New("Value not found in range")
	}
}

// Counts the total number of items contained in the array of ranges.
func totalInRange(tables []*unicode.RangeTable) int {
	total := 0
	for _, table := range tables {
		for _, r16 := range table.R16 {
			total += int((r16.Hi-r16.Lo)/r16.Stride) + 1
		}
		for _, r32 := range table.R32 {
			total += int((r32.Hi-r32.Lo)/r32.Stride) + 1
		}
	}
	return total
}

func notNilRand(rnd *rand.Rand) *rand.Rand {
	if rnd != nil {
		return rnd
	}
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

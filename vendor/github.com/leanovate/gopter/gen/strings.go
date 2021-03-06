package gen

import (
	"unicode"
	"unicode/utf8"

	"github.com/leanovate/gopter"
)

// RuneRange generates runes within a given range
func RuneRange(min, max rune) gopter.Gen {
	return genRune(Int64Range(int64(min), int64(max)))
}

// Rune generates an arbitrary character rune
func Rune() gopter.Gen {
	return genRune(Frequency(map[int]gopter.Gen{
		0xD800:                Int64Range(0, 0xD800),
		utf8.MaxRune - 0xDFFF: Int64Range(0xDFFF, int64(utf8.MaxRune)),
	}))
}

// RuneNoControl generates an arbitrary character rune that is not a control character
func RuneNoControl() gopter.Gen {
	return genRune(Frequency(map[int]gopter.Gen{
		0xD800:                Int64Range(32, 0xD800),
		utf8.MaxRune - 0xDFFF: Int64Range(0xDFFF, int64(utf8.MaxRune)),
	}))
}

func genRune(int64Gen gopter.Gen) gopter.Gen {
	return int64Gen.Map(func(value interface{}) interface{} {
		return rune(value.(int64))
	}).SuchThat(func(v interface{}) bool {
		return utf8.ValidRune(v.(rune))
	})
}

// NumChar generates arbitrary numberic character runes
func NumChar() gopter.Gen {
	return RuneRange('0', '9')
}

// AlphaUpperChar generates arbitrary uppercase alpha character runes
func AlphaUpperChar() gopter.Gen {
	return RuneRange('A', 'Z')
}

// AlphaLowerChar generates arbitrary lowercase alpha character runes
func AlphaLowerChar() gopter.Gen {
	return RuneRange('a', 'z')
}

// AlphaChar generates arbitrary character runes (upper- and lowercase)
func AlphaChar() gopter.Gen {
	return Frequency(map[int]gopter.Gen{
		0: AlphaUpperChar(),
		9: AlphaLowerChar(),
	})
}

// AlphaNumChar generate arbitrary alpha-numeric character runes
func AlphaNumChar() gopter.Gen {
	return Frequency(map[int]gopter.Gen{
		0: NumChar(),
		9: AlphaChar(),
	})
}

// AnyString generates an arbitrary string
func AnyString() gopter.Gen {
	return genString(Rune(), utf8.ValidRune)
}

// AlphaString generates an arbitrary string with letters
func AlphaString() gopter.Gen {
	return genString(AlphaChar(), unicode.IsLetter)
}

// NumString generates an arbitrary string with digits
func NumString() gopter.Gen {
	return genString(NumChar(), unicode.IsDigit)
}

// Identifier generates an arbitrary identifier string
// Identitiers are supporsed to start with a lowercase letter and contain only
// letters and digits
func Identifier() gopter.Gen {
	return gopter.CombineGens(
		AlphaLowerChar(),
		SliceOf(AlphaNumChar()),
	).Map(func(v interface{}) interface{} {
		values := v.([]interface{})
		first := values[0].(rune)
		tail := values[1].([]rune)
		result := make([]rune, 0, len(tail)+1)
		return string(append(append(result, first), tail...))
	}).SuchThat(func(v interface{}) bool {
		str := v.(string)
		if len(str) < 1 || !unicode.IsLower(([]rune(str))[0]) {
			return false
		}
		for _, ch := range str {
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
				return false
			}
		}
		return true
	}).WithShrinker(StringShrinker)
}

func genString(runeGen gopter.Gen, runeSieve func(ch rune) bool) gopter.Gen {
	return SliceOf(runeGen).Map(runesToString).SuchThat(func(v interface{}) bool {
		for _, ch := range v.(string) {
			if !runeSieve(ch) {
				return false
			}
		}
		return true
	}).WithShrinker(StringShrinker)
}

func runesToString(v interface{}) interface{} {
	return string(v.([]rune))
}

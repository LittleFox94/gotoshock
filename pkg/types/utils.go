package types

import "strings"

// bitstring takes a bool array and returns a string with a 1 or 0 rune for
// each entry in the given array.
func bitstring(v []bool) string {
	ret := strings.Builder{}
	for _, e := range v {
		if e {
			ret.WriteRune('1')
		} else {
			ret.WriteRune('0')
		}
	}

	return ret.String()
}

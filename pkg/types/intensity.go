package types

import (
	"fmt"
	"strconv"
)

// Intensity defines how strong a given operation (shock or vibrate) should be
// delivered.
type Intensity uint8

// String returns a string representation of the Intensity.
func (i Intensity) String() string {
	return fmt.Sprintf("%d", i)
}

// Set parses the given string into the Intensity this method was called on.
// Returns an error if not parsable.
func (i *Intensity) Set(v string) error {
	parsed, err := strconv.ParseUint(v, 10, 8)
	if err != nil {
		return fmt.Errorf("error parsing number: %w", err)
	}

	if parsed > 100 {
		return fmt.Errorf("intensity out of range: %w", ErrUnparsable)
	}

	*i = Intensity(parsed)
	return nil
}

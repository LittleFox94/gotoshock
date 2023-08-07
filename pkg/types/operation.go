package types

import "fmt"

// Operation defines what the shocker is supposed to do when receiving the
// message. See the Operation* constants in this package.
type Operation uint8

const (
	// OperationShock tells the shocker to deliver an electroshock with a given intensity (see SetIntensity).
	OperationShock Operation = 1

	// OperationVibrate tells the shocker to vibrate with a given intensity (see SetIntensity).
	OperationVibrate Operation = 2

	// OperationBeep tells the shocker to beep (intensity not used).
	OperationBeep Operation = 4
)

// String returns a string representation of the Operation.
func (op Operation) String() string {
	switch op {
	case OperationShock:
		return "shock"
	case OperationVibrate:
		return "vibrate"
	case OperationBeep:
		return "beep"
	default:
		return fmt.Sprintf("unknown operation (%v)", int(op))
	}
}

// Set takes the given string and parses it into the Operating this method is
// called on, returning an error if it cannot be parsed.
func (op *Operation) Set(s string) error {
	switch s {
	case "shock", "zap":
		*op = OperationShock
	case "vibrate":
		*op = OperationVibrate
	case "beep":
		*op = OperationBeep
	default:
		return ErrUnknownOperation
	}

	return nil
}

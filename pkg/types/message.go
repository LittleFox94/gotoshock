package types

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// header of all messages
	messageHeader = "01"

	// not-yet-known data in the message, probably remote ID - see
	// documentation on the Message type for more information.
	messageUnknown = "00101110001010110"

	// footer of all messages
	messageFooter = "00"
)

// Message is the digital protocol message understood by the Petrainer shockers.
//
// To ensure having useful defaults in the Message, please use the NewMessage
// function to create an instance. If all values will be set anyway, you can
// just create an instance in any other way, but you then have to ensure
// calling ALL Set* methods on it, as some of those also set a checksum for the
// value. A totally empty Message does NOT have valid checksums and calling
// .Build() does NOT calculate the checksums.
//
// It seems to be made of:
//   * 2 bits header (constant 01)
//   * 4 bits channel indicator (0000 for channel 1, 1110 for channel two)
//   * 3 bits operation (see the Operation constants)
//   * 17 bits of unknown use, probably remote ID
//     - 00101110001010110 is used in this library, while
//       10111010010101110 is used in the example on buttplug.io
//     - both those values have 110 at the end, maybe that's fixed then?
//	   - 14 bit remote ID and 3 bits for an even more unknown purpose sounds nice
//   * 7 bits for intensity (but capped at 100 by the shockers)
//   * 3 bits operation but different (swapped bit-order and negated bits)
//   * 4 bits channel but different (swapped bit-order and negated bits)
//   * 2 bits footer (constant 00)
type Message [42]bool

// NewMessage constructs a new instance of the Message type and sets safe
// defaults for all the values (Channel 1, Operation Beep, Intensity 0). Since
// most usages will set other values, it does not call .Build(), meaning the
// returned Message is not ready to be transmitted.
func NewMessage() *Message {
	return new(Message).
		SetChannel(Channel1).
		SetOperation(OperationBeep).
		SetIntensity(0)
}

// String returns a string representation of the message.
func (m Message) String() string {
	return bitstring(m[0:])
}

// DebugString returns a human-readable string representation of the message.
func (m Message) DebugString() string {
	op, opVerify, err := m.GetOperation()
	opVerified := "(verified)"

	if err != nil && errors.Is(err, ErrVerificationFailed) {
		opVerified = fmt.Sprintf(
			"(NOT verified: %s)",
			opVerify.String(),
		)
	}

	opString := fmt.Sprintf("operation: %s %s", op, opVerified)

	ch, chVerify, err := m.GetChannel()
	chVerified := "(verified)"

	if err != nil && errors.Is(err, ErrVerificationFailed) {
		chVerified = fmt.Sprintf(
			"(NOT verified: %s)",
			chVerify,
		)
	}

	chString := fmt.Sprintf("channel: %s %s", ch, chVerified)

	return strings.Join([]string{
		fmt.Sprintf("header: %s", m.GetHeader()),
		chString,
		opString,
		fmt.Sprintf("unknown part: %s", m.GetUnknown()),
		fmt.Sprintf("intensity: %v", m.GetIntensity()),
		fmt.Sprintf("footer: %s", m.GetFooter()),
		fmt.Sprintf("full message: %s", m.String()),
	}, "; ")
}

// SetChannel sets the channel this message is sent on, returning the Message
// so you can use this as a Builder-pattern method.
func (m *Message) SetChannel(ch Channel) *Message {
	for i := 0; i < 4; i++ {
		m[2+i] = (ch >> (3 - i) & 1) == 1
		m[36+i] = (ch >> i & 1) == 0
	}

	return m
}

// SetOperation sets the given operation on the Message, returning the Message
// so you can use this as a Builder-pattern method.
func (m *Message) SetOperation(op Operation) *Message {
	for i := 0; i < 3; i++ {
		m[6+i] = (op >> (2 - i) & 1) == 1
		m[33+i] = (op >> (i) & 1) == 0
	}

	return m
}

// SetIntensity sets the intensity value in the Message, returning the Message
// so you can use this as a Builder-pattern method.
//
// The intensity is capped at 100 by the shockers.
func (m *Message) SetIntensity(intensity uint8) *Message {
	for i := 0; i < 7; i++ {
		m[25+i] = (intensity >> (7 - i) & 1) == 1
	}

	return m
}

// Build prepares the Message to be sent, setting the constant bits as needed.
func (m *Message) Build() *Message {
	for i := 0; i < len(messageHeader); i++ {
		m[i] = messageHeader[i] == '1'
	}

	for i := 0; i < len(messageUnknown); i++ {
		m[9+i] = messageUnknown[i] == '1'
	}

	for i := 0; i < len(messageFooter); i++ {
		m[len(m)-len(messageFooter)+i] = messageFooter[i] == '1'
	}

	return m
}

// GetHeader extracts the header from the message.
func (m Message) GetHeader() string {
	return bitstring(m[0:2])
}

// GetFooter extracts the footer from the message.
func (m Message) GetFooter() string {
	return bitstring(m[40:42])
}

// GetUnknown extracts the unknown part from the message.
func (m Message) GetUnknown() string {
	return bitstring(m[9:25])
}

// GetChannel returns the channel value and channel verify value encoded
// in the message, returning an error if the verify value of the channel is
// not matching.
func (m Message) GetChannel() (Channel, Channel, error) {
	ch := Channel(0)
	chVerify := Channel(0)

	for i := 0; i < 4; i++ {
		if m[2+i] {
			ch |= 1 << (3 - i)
		}

		if !m[36+i] {
			chVerify |= 1 << i
		}
	}

	if ch != chVerify {
		return ch, chVerify, ErrChannelVerificationFailed
	}

	return ch, chVerify, nil
}

// GetOperation returns the operating value and operation verify value encoded
// in the message, returning an error if the verify value of the operation is
// not matching.
func (m Message) GetOperation() (Operation, Operation, error) {
	op := Operation(0)
	opVerify := Operation(0)

	for i := 0; i < 3; i++ {
		if m[6+i] {
			op |= 1 << (2 - i)
		}

		if !m[33+i] {
			opVerify |= 1 << i
		}
	}

	if op != opVerify {
		return op, opVerify, ErrOperationVerificationFailed
	}

	return op, opVerify, nil
}

// GetIntensity extracts the intensity from the message.
func (m Message) GetIntensity() uint8 {
	intensity := 0
	for i := 0; i < 7; i++ {
		if m[25+i] {
			intensity |= 1 << (7 - i)
		}
	}

	return uint8(intensity)
}

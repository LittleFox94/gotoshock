package types

import "fmt"

// Channel defines the channel this message is sent on. One remote can send
// messages on two different channels. The protocol might allow more channels,
// but this is not verified yet - so best use the Channel* constants in this
// package.
type Channel uint8

const (
	// Channel1 defines shockers listening on channel 1 should act on the Message.
	Channel1 Channel = 0

	// Channel2 defines shockers listening on channel 2 should act on the Message.
	Channel2 Channel = 14
)

// String returns a string representation of the Channel.
func (ch Channel) String() string {
	switch ch {
	case Channel1:
		return "1"
	case Channel2:
		return "2"
	default:
		return fmt.Sprintf("unknown channel (%v)", int(ch))
	}
}

// Set takes the given string and parses it into the Channel this method is
// called on, returning an error if it cannot be parsed.
func (ch *Channel) Set(s string) error {
	switch s {
	case "1":
		*ch = Channel1
	case "2":
		*ch = Channel2
	default:
		return ErrUnknownChannel
	}

	return nil
}

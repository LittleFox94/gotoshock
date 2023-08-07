package types

import (
	"errors"
	"fmt"
)

var (
	// ErrVerificationFailed is returned when a given value has a checksum in
	// the Message and the checksum and value do not match.
	ErrVerificationFailed = errors.New("value in message does not match verification value")

	// ErrOperationVerificationFailed is returned when the operation and
	// operation checksum in a Message don't match.
	ErrOperationVerificationFailed = fmt.Errorf("operation verification failed: %w", ErrVerificationFailed)

	// ErrOperationVerificationFailed is returned when the channel and channel
	// checksum in a Message don't match.
	ErrChannelVerificationFailed = fmt.Errorf("channel verification failed: %w", ErrVerificationFailed)

	// ErrUnparsable is returned when a given value cannot be parsed into a
	// specific type.
	ErrUnparsable = errors.New("parsing failed")

	// ErrUnknownOperation is returned when Operation.Set() is called with a
	// string that does not match any known Operation.
	ErrUnknownOperation = fmt.Errorf("%w: unknown operation", ErrUnparsable)

	// ErrUnknownChannel is returned when Channel.Set() is called with a
	// string that does not match any known Channel.
	ErrUnknownChannel = fmt.Errorf("%w: unknown channel", ErrUnparsable)
)

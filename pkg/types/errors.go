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
)

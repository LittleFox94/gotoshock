package driver

import (
	"time"

	"praios.lf-net.org/littlefox/gotoshock/pkg/types"
)

// BitstreamDriver is the interface of lowlevel I/O drivers, like for GPIOs
type BitstreamDriver interface {
	Output(stream []bool, between time.Duration) error
}

// MessageDriver is the interface of highlevel drivers, directly sending
// digital protocol data.
type MessageDriver interface {
	Output(message *types.Message) error
}

// BindableMessageDriver is a MessageDriver that can be bound to a given
// BitstreamDriver instead of being independent from any other driver.
type BindableMessageDriver interface {
	MessageDriver
	Bind(io BitstreamDriver) error
}

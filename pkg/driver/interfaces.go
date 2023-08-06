package driver

import (
	"time"

	"praios.lf-net.org/littlefox/gotoshock/pkg/types"
)

// IODriver is the interface of lowlevel I/O drivers, like for GPIOs
type IODriver interface {
	Output(stream []bool, between time.Duration) error
}

// PWMDriver is the interface of highlevel drivers, directly sending digitial
// protocol data.
type PWMDriver interface {
	Output(message *types.Message) error
}

// BindablePWMDriver is a PWMDriver that can be bound to a given IODriver
// instead of being independent from any other driver.
type BindablePWMDriver interface {
	PWMDriver
	Bind(io IODriver) error
}

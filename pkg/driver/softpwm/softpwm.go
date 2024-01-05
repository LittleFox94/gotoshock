package softpwm

import (
	"time"

	"praios.lf-net.org/littlefox/gotoshock/pkg/driver"
	"praios.lf-net.org/littlefox/gotoshock/pkg/types"
)

type softpwm struct {
	io driver.BitstreamDriver
}

func (s softpwm) Output(m *types.Message) error {
	if s.io == nil {
		return driver.ErrIODriverNotBound
	}

	preamble := "00000000000000011111"
	trailer := ""

	bitstring := make([]bool, 0, len(preamble)+len(trailer)+len(m)*4)

	for _, v := range preamble {
		bitstring = append(bitstring, v == '1')
	}

	for _, v := range m {
		bitstring = append(bitstring, true)

		bitstring = append(bitstring, v)
		bitstring = append(bitstring, v)
		bitstring = append(bitstring, false)
	}

	for _, v := range trailer {
		bitstring = append(bitstring, v == '1')
	}

	return s.io.Output(bitstring, 250*time.Microsecond)
}

func (s *softpwm) Bind(io driver.BitstreamDriver) error {
	s.io = io
	return nil
}

func init() {
	driver.RegisterMessage("softpwm", func(args []string) (driver.MessageDriver, error) {
		return &softpwm{}, nil
	})
}

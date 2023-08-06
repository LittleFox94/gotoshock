package gpio

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
	"praios.lf-net.org/littlefox/gotoshock/pkg/driver"
)

type raspi_gpio struct {
	pin rpio.Pin
}

func (r raspi_gpio) Output(stream []bool, d time.Duration) error {
	t := time.Now()
	for i := 0; i < len(stream); {
		if t.Add(d).After(time.Now()) {
			continue
		}

		t = time.Now()

		s := rpio.Low
		if stream[i] {
			s = rpio.High
		}

		r.pin.Write(s)
		i++
	}

	time.Sleep(time.Duration(len(stream)+1) * d)

	return nil
}

func init() {
	driver.RegisterIO("raspi_gpio", func(args []string) (driver.IODriver, error) {
		if len(args) != 1 {
			return nil, errors.New("invalid arguments, needs exactly one argument: pin number to use (BCM2835 pin numbering)")
		}

		if err := rpio.Open(); err != nil {
			return nil, fmt.Errorf("error opening GPIO: %w", err)
		}

		pinNumber, err := strconv.ParseUint(args[0], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("error parsing pin to use: %w", err)
		}

		ret := raspi_gpio{
			pin: rpio.Pin(pinNumber),
		}
		log.Printf("raspi_gpio: pin number %v", pinNumber)

		ret.pin.Output()

		return ret, nil
	})
}

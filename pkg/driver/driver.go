package driver

import (
	"errors"
	"fmt"
	"strings"
	"text/scanner"
)

var messageDriverRegistry map[string]MessageDriverFactory

var bitstreamDriverRegistry map[string]BitstreamDriverFactory

func RegisterMessage(name string, fac MessageDriverFactory) {
	if messageDriverRegistry == nil {
		messageDriverRegistry = make(map[string]MessageDriverFactory)
	}

	messageDriverRegistry[name] = fac
}

func RegisterBitstream(name string, fac BitstreamDriverFactory) {
	if bitstreamDriverRegistry == nil {
		bitstreamDriverRegistry = make(map[string]BitstreamDriverFactory)
	}

	bitstreamDriverRegistry[name] = fac
}

func Setup(conn string) (MessageDriver, error) {
	type driverWithArgs struct {
		driver string
		args   []string
	}

	drivers := make([]driverWithArgs, 0)

	s := scanner.Scanner{
		Mode: scanner.ScanIdents | scanner.ScanStrings,
	}

	s.Init(strings.NewReader(conn))
	for token := s.Scan(); token != scanner.EOF; token = s.Scan() {
		switch token {
		case scanner.Ident:
			drivers = append(drivers, driverWithArgs{
				driver: s.TokenText(),
				args:   make([]string, 0),
			})
		case scanner.Int, scanner.Float, scanner.Char, scanner.String, scanner.RawString:
			if len(drivers) == 0 {
				return nil, errors.New("syntax error in driver string")
			}

			drivers[len(drivers)-1].args = append(drivers[len(drivers)-1].args, s.TokenText())
		default:
			return nil, errors.New("syntax error in driver string")
		}
	}

	if len(drivers) == 0 || len(drivers) > 2 {
		return nil, errors.New("invalid driver number")
	}

	messageDriverFactory, ok := messageDriverRegistry[drivers[0].driver]
	if !ok {
		return nil, fmt.Errorf("PWM driver %q not found", drivers[0].driver)
	}

	messageDriver, err := messageDriverFactory(drivers[0].args)
	if err != nil {
		return nil, fmt.Errorf("error initializing PWM driver: %w", err)
	}

	if len(drivers) == 1 {
		return messageDriver, nil
	}

	bindableMessageDriver, ok := messageDriver.(BindableMessageDriver)
	if !ok {
		return nil, fmt.Errorf("PWM driver %q cannot be bound to another driver", drivers[0].driver)
	}

	ioDriverFactory, ok := bitstreamDriverRegistry[drivers[1].driver]
	if !ok {
		return nil, fmt.Errorf("I/O driver %q not found", drivers[1].driver)
	}

	ioDriver, err := ioDriverFactory(drivers[1].args)
	if err != nil {
		return nil, fmt.Errorf("error initializing I/O driver: %w", err)
	}

	bindableMessageDriver.Bind(ioDriver)
	return bindableMessageDriver, nil
}

package driver

type PWMDriverFactory func(arguments []string) (PWMDriver, error)

type IODriverFactory func(arguments []string) (IODriver, error)

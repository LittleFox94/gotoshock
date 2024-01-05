package driver

type BitstreamDriverFactory func(arguments []string) (BitstreamDriver, error)

type MessageDriverFactory func(arguments []string) (MessageDriver, error)

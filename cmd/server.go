package main

import (
	"log"
	"os"

	"praios.lf-net.org/littlefox/gotoshock/pkg/driver"
	"praios.lf-net.org/littlefox/gotoshock/pkg/types"

	_ "praios.lf-net.org/littlefox/gotoshock/pkg/driver/raspi/gpio"
	_ "praios.lf-net.org/littlefox/gotoshock/pkg/driver/softpwm"
)

func main() {
	pwmDriver, err := driver.Setup(os.Args[1])
	if err != nil {
		log.Fatalf("error initializing driver: %v", err)
	}

	msg := types.NewMessage().
		SetIntensity(100).
		SetOperation(types.OperationVibrate).
		Build()

	for i := 0; i < 8; i++ {
		if err := pwmDriver.Output(msg); err != nil {
			log.Fatalf("error sending message: %v", err)
		}
	}
}

package main

import (
	"log"
	"net/http"
	"os"

	"praios.lf-net.org/littlefox/gotoshock/pkg/driver"
	"praios.lf-net.org/littlefox/gotoshock/pkg/server/api/v1alpha1"

	_ "praios.lf-net.org/littlefox/gotoshock/pkg/driver/raspi/gpio"
	_ "praios.lf-net.org/littlefox/gotoshock/pkg/driver/softpwm"
)

func main() {
	pwmDriver, err := driver.Setup(os.Args[1])
	if err != nil {
		log.Fatalf("error initializing driver: %v", err)
	}

	routes, err := v1alpha1.Routes(pwmDriver)
	if err != nil {
		log.Fatalf("error initializing router: %v", err)
	}
	http.ListenAndServe(":8080", routes)

	/*
		msg := types.NewMessage().
			SetIntensity(20).
			SetOperation(types.OperationVibrate).
			Build()

		for i := 0; i < 4; i++ {
			if err := pwmDriver.Output(msg); err != nil {
				log.Fatalf("error sending message: %v", err)
			}
		}*/
}

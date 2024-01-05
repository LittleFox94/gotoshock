package v1alpha1

import (
	"fmt"
	"net/http"

	"praios.lf-net.org/littlefox/gotoshock/pkg/driver"
	"praios.lf-net.org/littlefox/gotoshock/pkg/server/typesafe_router"
	"praios.lf-net.org/littlefox/gotoshock/pkg/types"
)

type apikey string

func (a *apikey) Set(v string) error {
	*a = apikey(v)
	return nil
}

func (a apikey) String() string {
	return string(a)
}

type routes struct {
	typesafe_router.TypeSafeRouter

	driver driver.MessageDriver
}

func (routes routes) postMessageHandler(res http.ResponseWriter, req *http.Request, key apikey, channel types.Channel, operation types.Operation, intensity types.Intensity) {
	if key == "hellorld!" {
		res.Write([]byte(fmt.Sprintf("Hello %v!\nsending %v with intensity %v on channel %v\n", key, operation, intensity, channel)))

		msg := types.NewMessage().
			SetChannel(channel).
			SetOperation(operation).
			SetIntensity(intensity).
			Build()

		for i := 0; i < 4; i++ {
			if err := routes.driver.Output(msg); err != nil {
				res.WriteHeader(500)
				res.Write([]byte(err.Error()))
			}
		}
	} else {
		res.WriteHeader(401)
	}
}

func Routes(driver driver.MessageDriver) (http.Handler, error) {
	ret := routes{driver: driver}

	type route struct {
		method  string
		path    string
		handler any
	}

	routes := map[string]route{
		"postMessage": {"POST", "/v1alpha1/message/:/:/:/:", ret.postMessageHandler},
	}

	for name, route := range routes {
		if err := ret.AddRoute(route.method, route.path, route.handler); err != nil {
			return nil, fmt.Errorf("error adding route %q: %w", name, err)
		}
	}

	return &ret, nil
}

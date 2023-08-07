package v1alpha1

import (
	"fmt"
	"net/http"

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

func postMessageHandler(res http.ResponseWriter, req *http.Request, key apikey, channel types.Channel, operation types.Operation, intensity types.Intensity) {
	res.Write([]byte(fmt.Sprintf("Hello %v!\nsending %v with intensity %v on channel %v\n", key, operation, intensity, channel)))
}

func Routes() (http.Handler, error) {
	router := typesafe_router.TypeSafeRouter{}

	type route struct {
		method  string
		path    string
		handler any
	}

	routes := map[string]route{
		"postMessage": {"POST", "/v1alpha1/message/:/:/:/:", postMessageHandler},
	}

	for name, route := range routes {
		if err := router.AddRoute(route.method, route.path, route.handler); err != nil {
			return nil, fmt.Errorf("error adding route %q: %w", name, err)
		}
	}

	return &router, nil
}

package typesafe_router

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var (
	errInvalidHandler = errors.New("invalid handler function")
)

var placeHolderRegex = regexp.MustCompile("/:")

type TypeSafeRouter struct {
	handlers map[string]map[*regexp.Regexp]http.HandlerFunc
}

// AddRoute takes a path with placeholders ('/:') to mark path elements going
// into arguments to the given handler function and adds that to the
// typesafeRouter instance for serving.
func (tr *TypeSafeRouter) AddRoute(method, path string, handler any) error {
	placeHolderPositions := placeHolderRegex.FindAllStringIndex(path, -1)

	handlerType := reflect.TypeOf(handler)
	if handlerType.NumIn() != len(placeHolderPositions)+2 {
		return fmt.Errorf("%w: number of path placeholders and handler arguments don't match", errInvalidHandler)
	}

	if !reflect.TypeOf((*http.ResponseWriter)(nil)).Elem().AssignableTo(handlerType.In(0)) ||
		!reflect.TypeOf(&http.Request{}).AssignableTo(handlerType.In(1)) {
		return fmt.Errorf("%w: first two handler arguments have to be a http.ResponseWriter and *http.Request", errInvalidHandler)
	}

	valueType := reflect.TypeOf((*flag.Value)(nil)).Elem()
	for i := 2; i < handlerType.NumIn(); i++ {
		if !reflect.PtrTo(handlerType.In(i)).Implements(valueType) {
			return fmt.Errorf("%w: all handlers parameters after http.ResponseWriter and http.Request must implement flag.Value", errInvalidHandler)
		}
	}

	pathRegexString := strings.Builder{}

	lastPosition := []int{0, 0}
	for _, placeHolderPosition := range placeHolderPositions {
		pathRegexString.WriteString(path[lastPosition[1]:placeHolderPosition[0]])
		pathRegexString.WriteString(`/([^/]+)`)
		lastPosition = placeHolderPosition
	}

	pathRegex, err := regexp.Compile(pathRegexString.String())
	if err != nil {
		return fmt.Errorf("error compiling path regex: %w", err)
	}

	if tr.handlers == nil {
		tr.handlers = make(map[string]map[*regexp.Regexp]http.HandlerFunc)
	}

	if _, ok := tr.handlers[method]; !ok {
		tr.handlers[method] = make(map[*regexp.Regexp]http.HandlerFunc)
	}

	tr.handlers[method][pathRegex] = func(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		matches := pathRegex.FindStringSubmatch(path)
		if len(matches) != len(placeHolderPositions)+1 {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		args := make([]reflect.Value, 2+len(matches)-1)
		args[0] = reflect.ValueOf(res)
		args[1] = reflect.ValueOf(req)

		for i, match := range matches[1:] {
			v := reflect.New(handlerType.In(2 + i))
			err := v.Interface().(flag.Value).Set(match)
			if err != nil {
				res.WriteHeader(400)
				res.Write([]byte(err.Error()))
				return
			}

			args[2+i] = v.Elem()
		}

		reflect.ValueOf(handler).Call(args)
	}

	return nil
}

func (tr TypeSafeRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	paths, ok := tr.handlers[req.Method]
	if ok {
		for pathRegex, handler := range paths {
			if pathRegex.MatchString(req.URL.Path) {
				handler(res, req)
				return
			}
		}
	}

	http.NotFoundHandler().ServeHTTP(res, req)
}

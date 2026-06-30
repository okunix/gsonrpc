package main

import (
	"errors"
	"fmt"
	"reflect"
)

type ExampleApp struct {
}

func (e *ExampleApp) Add(a, b float64) float64 {
	return a + b
}

func (e *ExampleApp) Echo(s string) string {
	return s
}

type RPC struct {
	app any
}

func NewRPC(app any) RPC {
	return RPC{app: app}
}

func (r *RPC) Call(name string, args ...any) ([]any, error) {
	defer func() {
		if issue := recover(); issue != nil {
			fmt.Printf("unknown issue: %v\n", issue)
		}
	}()
	appType := reflect.ValueOf(r.app)
	if appType.IsZero() {
		return nil, errors.New("app is nil")
	}
	method := appType.MethodByName(name)
	if method.IsZero() {
		return nil, errors.New("method not found")
	}
	if len(args) != method.Type().NumIn() {
		return nil, errors.New("arguments mismatch")
	}
	vals := make([]reflect.Value, 0, len(args))
	for i, arg := range args {
		want := method.Type().In(i)
		val := reflect.ValueOf(arg)
		if val.Type() == want {
			vals = append(vals, val)
		} else if val.Type().ConvertibleTo(want) {
			vals = append(vals, val.Convert(want))
		} else {
			return nil, fmt.Errorf("argument %d can't convert type %s to %s", i, val.Type(), want)
		}
	}
	results := method.Call(vals)
	out := make([]any, 0, len(results))
	for _, result := range results {
		out = append(out, result.Interface())
	}
	return out, nil
}

func main() {
	exampleApp := &ExampleApp{}
	rpc := NewRPC(exampleApp)
	results, err := rpc.Call("Echo", "2")
	if err != nil {
		panic(err)
	}
	fmt.Printf("results: %v\n", results)
}

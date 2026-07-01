package main

import (
	"errors"
	"fmt"
	"reflect"
)

type RPC struct {
	methods map[string]reflect.Value
}

func NewRPC() RPC {
	return RPC{methods: make(map[string]reflect.Value)}
}

func (r *RPC) Register(apps ...any) error {
	return r.RegisterWithPrefix("", apps...)
}

func (r *RPC) RegisterWithPrefix(prefix string, apps ...any) error {
	for _, app := range apps {
		structValue := reflect.ValueOf(app)
		if structValue.IsZero() {
			return errors.New("app is nil")
		}

		for k, v := range structValue.Methods() {
			r.methods[prefix+k.Name] = v
		}
	}
	return nil
}

func (r *RPC) Call(name string, params ...any) (out []any, err error) {
	defer func() {
		if issue := recover(); issue != nil {
			rpcErr := ErrInternalError
			rpcErr.Data = fmt.Sprintf("%v", issue)
			err = rpcErr
		}
	}()

	method, ok := r.methods[name]
	if !ok {
		return nil, ErrMethodNotFound
	}

	if (len(params) != method.Type().NumIn() && !method.Type().IsVariadic()) ||
		(method.Type().IsVariadic() && len(params) < method.Type().NumIn()-1) {
		return nil, ErrInvalidParams
	}

	// type conversion stuff
	vals := make([]reflect.Value, 0, len(params))
	for i, arg := range params {
		val := reflect.ValueOf(arg)

		var want reflect.Type
		if method.Type().NumIn()-1 <= i && method.Type().IsVariadic() {
			want = method.Type().In(method.Type().NumIn() - 1).Elem()
		} else {
			want = method.Type().In(i)
		}

		if val.Type() == want {
			vals = append(vals, val)
		} else if val.Type().ConvertibleTo(want) {
			vals = append(vals, val.Convert(want))
		} else {
			return nil, ErrInvalidParams
		}
	}

	// calling actual method
	results := method.Call(vals)
	out = make([]any, 0, len(results))
	for _, result := range results {
		out = append(out, result.Interface())
	}

	return out, nil
}

func (r *RPC) CallNamed(name string, params map[string]any) (out []any, err error) {
	defer func() {
		if issue := recover(); issue != nil {
			rpcErr := ErrInternalError
			rpcErr.Data = fmt.Sprintf("%v", issue)
			err = rpcErr
		}
	}()
	panic("call by named parameters is unimplemented")

	method, ok := r.methods[name]
	if !ok {
		return nil, ErrMethodNotFound
	}
	if len(params) != method.Type().NumIn() {
		return nil, ErrInvalidParams
	}

	//values := make([]reflect.Value, 0, len(params))
	inputValues := method.Type().Ins()
	for i := range inputValues {
		fmt.Printf("i: %+v\n", i.String())
	}

	return nil, nil
}

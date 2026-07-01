package main

import (
	"fmt"
	"testing"
)

type ExampleApp struct {
	accumulator uint
}

func (e *ExampleApp) Add(a, b float64) float64 {
	e.accumulator++
	return a + b
}

func (e *ExampleApp) Sum(a ...float64) float64 {
	e.accumulator++
	sum := float64(0)
	for _, v := range a {
		sum += v
	}
	return sum
}

func (e *ExampleApp) SumMul(mul float64, a ...float64) float64 {
	e.accumulator++
	sum := float64(0)
	for _, v := range a {
		sum += v
	}
	return sum * mul
}

func (e *ExampleApp) Echo(s string) string {
	e.accumulator++
	return s
}

func (e *ExampleApp) PrintAccumulator() {
	fmt.Println(e.accumulator)
}

var (
	app = &ExampleApp{}
)

func TestJsonRPC(t *testing.T) {
	rpc := NewRPC()
	rpc.RegisterWithPrefix("math.", app)
	jsonRPC := NewJsonRPC(&rpc)

	req := []byte(`{"jsonrpc": "2.0", "method": "math.Sum", "params": [3, 4], "id": "111"}`)
	resp := jsonRPC.ProcessRequest(req)
	t.Log(string(resp))
}

func TestEcho(t *testing.T) {
	rpc := NewRPC()
	rpc.RegisterWithPrefix("math.", app)
	jsonRPC := NewJsonRPC(&rpc)

	req := []byte(
		`{"jsonrpc": "2.0", "method": "math.SumMul", "params": [2, 2, 2, 2], "id": "111"}`,
	)
	resp := jsonRPC.ProcessRequest(req)
	t.Log(string(resp))
}

func TestSumMul(t *testing.T) {
	rpc := NewRPC()
	rpc.RegisterWithPrefix("math.", app)
	jsonRPC := NewJsonRPC(&rpc)

	req := []byte(`{"jsonrpc": "2.0", "method": "math.Echo", "params": ["hello"], "id": "111"}`)
	resp := jsonRPC.ProcessRequest(req)
	t.Log(string(resp))
}

func TestNotify(t *testing.T) {
	rpc := NewRPC()
	rpc.RegisterWithPrefix("math.", app)
	jsonRPC := NewJsonRPC(&rpc)

	req := []byte(`{"jsonrpc": "2.0", "method": "math.PrintAccumulator", "id": null}`)
	resp := jsonRPC.ProcessRequest(req)
	t.Log(resp)
}

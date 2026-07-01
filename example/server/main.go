package main

import (
	"log"

	"github.com/okunix/gsonrpc"
)

type ExampleApp struct {
}

func (e *ExampleApp) Example(name string) string {
	log.Println("called method example")
	return "hello " + name
}

func main() {
	app := &ExampleApp{}
	rpc := gsonrpc.NewRPC()
	rpc.RegisterWithPrefix("example.", app)
	jsonrpc := gsonrpc.NewJsonRPC(&rpc)
	jsonrpc.ListenAndServe(":8080")
}

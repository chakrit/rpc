package main

//go:generate rpc -gen go -out ./api todo.rpc
//go:generate go fmt ./api/...

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chakrit/rpc/todo/api"
	"github.com/chakrit/rpc/todo/api/client"
	"github.com/chakrit/rpc/todo/api/server"
)

func main() {
	var addr string
	flag.StringVar(&addr, "a", "0.0.0.0:9999", "address to bind")
	flag.Parse()

	go runServer(addr)
	runClient(addr)
	os.Exit(0) // also terminates server routine
}

func runServer(addr string) {
	opts := server.Options{Addr: addr}
	srv := server.New(&opts)
	srv.Handler = &handler{}
	if err := srv.Listen(); err != nil {
		log.Fatal(err)
	}
}

func runClient(addr string) {
	opts := client.Options{Addr: addr}
	cl := client.New(&opts)

	if items, err := cl.List(); err != nil {
		log.Fatal(err)
	} else {
		logOutput("List", items...)
	}

	alpha := &api.TodoItem{Description: "alpha", Done: false}
	if item, err := cl.Update("alpha", alpha); err != nil {
		log.Fatal(err)
	} else {
		logOutput("Update", item)
	}

	beta := &api.TodoItem{Description: "beta", Done: true}
	if item, err := cl.Update("beta", beta); err != nil {
		log.Fatal(err)
	} else {
		logOutput("Update", item)
	}

	if items, err := cl.List(); err != nil {
		log.Fatal(err)
	} else {
		logOutput("List", items...)
	}

	if item, err := cl.Destroy("alpha"); err != nil {
		log.Fatal(err)
	} else {
		logOutput("Destroy", item)
	}

	if item, err := cl.Destroy("beta"); err != nil {
		log.Fatal(err)
	} else {
		logOutput("Destroy", item)
	}

	if items, err := cl.List(); err != nil {
		log.Fatal(err)
	} else {
		logOutput("List", items...)
	}
}

func logOutput(name string, items ...*api.TodoItem) {
	fmt.Printf("%s\n", name)
	for _, item := range items {
		if item.Done {
			fmt.Printf("[%s] %s !DONE!\n", item.ID, item.Description)
		} else {
			fmt.Printf("[%s] %s\n", item.ID, item.Description)
		}
	}
}

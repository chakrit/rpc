package main

//go:generate rpc -gen go -out ./api todo.rpc
//go:generate go fmt ./api/...

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chakrit/rpc-todo/api"
	"github.com/chakrit/rpc-todo/api/client"
	"github.com/chakrit/rpc-todo/api/server"
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
	var (
		opts = client.Options{Addr: addr}
		cl   = client.New(&opts)
		ctx  = context.Background()
	)

	if items, err := cl.List(ctx); err != nil {
		log.Fatal(err)
	} else {
		logOutput("List", items...)
	}

	alpha := &api.TodoItem{Description: "alpha", Done: false}
	if item, err := cl.Update(ctx, "alpha", alpha); err != nil {
		log.Fatal(err)
	} else {
		logOutput("Update", item)
	}

	beta := &api.TodoItem{Description: "beta", Done: true}
	if item, err := cl.Update(ctx, "beta", beta); err != nil {
		log.Fatal(err)
	} else {
		logOutput("Update", item)
	}

	if items, err := cl.List(ctx); err != nil {
		log.Fatal(err)
	} else {
		logOutput("List", items...)
	}

	if item, err := cl.Destroy(ctx, "alpha"); err != nil {
		log.Fatal(err)
	} else {
		logOutput("Destroy", item)
	}

	if item, err := cl.Destroy(ctx, "beta"); err != nil {
		log.Fatal(err)
	} else {
		logOutput("Destroy", item)
	}

	if items, err := cl.List(ctx); err != nil {
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

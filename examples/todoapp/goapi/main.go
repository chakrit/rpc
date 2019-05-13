package main

//go:generate rpc -gen go -out ./api todo.rpc
//go:generate go fmt ./api/...

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	"github.com/gorilla/handlers"

	"github.com/chakrit/rpc/todo/api"
	"github.com/chakrit/rpc/todo/api/client"
	"github.com/chakrit/rpc/todo/api/server"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Example application using the RPC",
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server part of the RPC",
	Run:   runServerCmd,
}

var ClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Starts the client part of the RPC",
	Run:   runClientCmd,
}

var flags = struct {
	Addr string
}{}

func init() {
	RootCmd.AddCommand(ServerCmd, ClientCmd)
	RootCmd.PersistentFlags().StringVarP(&flags.Addr,
		"address", "a",
		"", "Server bind or client connect address.")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runServerCmd(cmd *cobra.Command, args []string) {
	opts := server.Options(flags)
	srv := server.New(&opts)
	srv.Handler = &handler{}

	// compose some middlewares cors/logging
	httpHandler := srv.HTTPHandler()
	corsHandler := cors.AllowAll().Handler(httpHandler)
	logHandler := handlers.LoggingHandler(os.Stdout, corsHandler)

	if err := http.ListenAndServe(opts.Addr, logHandler); err != nil {
		log.Fatal(err)
	}
}

func runClientCmd(cmd *cobra.Command, args []string) {
	opts := client.Options(flags)
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

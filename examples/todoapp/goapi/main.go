package main

//go:generate rpc -gen go -out ./api todo.rpc
//go:generate go fmt ./api/...

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chakrit/rpc/todo/api"
	"github.com/chakrit/rpc/todo/api/client"
	"github.com/chakrit/rpc/todo/api/server"
	"github.com/gorilla/handlers"
	"github.com/rs/cors"
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
	opts := server.Options{Addr: flags.Addr}
	srv := server.New(&opts)
	srv.Provider = provider{}

	// compose some middlewares cors/logging
	httpHandler := srv.HTTPHandler()
	corsHandler := cors.AllowAll().Handler(httpHandler)
	logHandler := handlers.LoggingHandler(os.Stdout, corsHandler)

	if err := http.ListenAndServe(opts.Addr, logHandler); err != nil {
		log.Fatal(err)
	}
}

func runClientCmd(cmd *cobra.Command, args []string) {
	var (
		opts = client.Options(flags)
		c    = client.New(&opts)
		ctx  = context.Background()

		check = func(err error) {
			if err != nil {
				log.Fatal(err)
			}
		}
	)

	list, err := c.List(ctx)
	check(err)
	logOutput("List", list...)

	alpha, err := c.Create(ctx, "alpha")
	check(err)
	logOutput("Create", alpha)

	list, err = c.List(ctx)
	check(err)
	logOutput("List", list...)

	beta, err := c.Create(ctx, "beta")
	check(err)
	logOutput("Create", beta)

	list, err = c.List(ctx)
	check(err)
	logOutput("List", list...)

	alpha, err = c.Destroy(ctx, alpha.ID)
	check(err)
	logOutput("Destroy", list...)

	list, err = c.List(ctx)
	check(err)
	logOutput("List", list...)
}

func logOutput(name string, items ...*api.TodoItem) {
	fmt.Printf("%s\n", name)
	for _, item := range items {
		fmt.Printf("[%d] %s\n", item.ID, item.Description)
	}
}

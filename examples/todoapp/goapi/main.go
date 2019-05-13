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
	var (
		opts = client.Options(flags)
		c    = client.New(&opts)

		check = func(err error) {
			if err != nil {
				log.Fatal(err)
			}
		}
	)

	list, err := c.List()
	check(err)
	logOutput("List", list...)

	alpha, err := c.Create("alpha")
	check(err)
	logOutput("Create", alpha)

	list, err = c.List()
	check(err)
	logOutput("List", list...)

	beta, err := c.Create("beta")
	check(err)
	logOutput("Create", beta)

	list, err = c.List()
	check(err)
	logOutput("List", list...)

	alpha, err = c.Destroy(alpha.ID)
	check(err)
	logOutput("Destroy", list...)

	list, err = c.List()
	check(err)
	logOutput("List", list...)
}

func logOutput(name string, items ...*api.TodoItem) {
	fmt.Printf("%s\n", name)
	for _, item := range items {
		fmt.Printf("[%d] %s\n", item.ID, item.Description)
	}
}

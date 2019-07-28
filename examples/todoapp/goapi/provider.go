package main

import (
	"github.com/chakrit/rpc/todo/api"
	"github.com/chakrit/rpc/todo/api/server"
)

type provider struct{}

var _ server.Provider_rpc_root = provider{}

func (p provider) Provide_rpc_root() api.Interface {
	return &handler{}
}

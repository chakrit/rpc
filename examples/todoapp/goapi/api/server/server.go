// <auto-generated />
//
// expected import: github.com/chakrit/rpc/todo/api/server
package server

import (
	"encoding/json"
	"net/http"

	rpc_root "github.com/chakrit/rpc/todo/api"
)

type Handler_rpc_root struct {
	Handler rpc_root.Interface
}

type Result struct {
	Error   error         `json:"error"`
	Returns []interface{} `json:"returns"`
}

type Server struct {
	Options
	Handler_rpc_root
}

type Options struct {
	Addr string
}

func New(opts *Options) *Server {
	return &Server{Options: *opts}
}

func (s *Server) Listen() error {
	return http.ListenAndServe(s.Options.Addr, s.HTTPHandler())
}

func (s *Server) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	s.register_rpc_root(mux, s.Handler_rpc_root)
	return mux
}

func (s *Server) register_rpc_root(
	mux *http.ServeMux,
	handler Handler_rpc_root,
) *http.ServeMux {

	mux.HandleFunc("/api/Create", func(resp http.ResponseWriter, req *http.Request) {
		var err error
		resp.Header().Set("Content-Type", "application/json")

		var arg0 string
		args := [1]interface{}{
			&arg0,
		}

		if req.Body != nil {
			if err := json.NewDecoder(req.Body).Decode(&args); err != nil {
				resp.WriteHeader(400)
				renderError(resp, err)
				return
			}
		}

		var out0 *rpc_root.TodoItem
		out0, err = handler.Handler.Create(arg0)
		result := &Result{
			Error: err,
			Returns: []interface{}{
				out0,
			},
		}

		bytes, err := json.Marshal(result)
		if err != nil {
			resp.WriteHeader(500)
			renderError(resp, err)
		}

		resp.WriteHeader(200)
		_, _ = resp.Write(bytes)
	})

	mux.HandleFunc("/api/Destroy", func(resp http.ResponseWriter, req *http.Request) {
		var err error
		resp.Header().Set("Content-Type", "application/json")

		var arg0 int64
		args := [1]interface{}{
			&arg0,
		}

		if req.Body != nil {
			if err := json.NewDecoder(req.Body).Decode(&args); err != nil {
				resp.WriteHeader(400)
				renderError(resp, err)
				return
			}
		}

		var out0 *rpc_root.TodoItem
		out0, err = handler.Handler.Destroy(arg0)
		result := &Result{
			Error: err,
			Returns: []interface{}{
				out0,
			},
		}

		bytes, err := json.Marshal(result)
		if err != nil {
			resp.WriteHeader(500)
			renderError(resp, err)
		}

		resp.WriteHeader(200)
		_, _ = resp.Write(bytes)
	})

	mux.HandleFunc("/api/List", func(resp http.ResponseWriter, req *http.Request) {
		var err error
		resp.Header().Set("Content-Type", "application/json")

		var out0 []*rpc_root.TodoItem
		out0, err = handler.Handler.List()
		result := &Result{
			Error: err,
			Returns: []interface{}{
				out0,
			},
		}

		bytes, err := json.Marshal(result)
		if err != nil {
			resp.WriteHeader(500)
			renderError(resp, err)
		}

		resp.WriteHeader(200)
		_, _ = resp.Write(bytes)
	})

	return mux
}

func renderError(resp http.ResponseWriter, e error) {
	result := &Result{
		Error:   e,
		Returns: nil,
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		_, _ = resp.Write([]byte(`{"error":"json processing error"}`))
	} else {
		_, _ = resp.Write(bytes)
	}
}

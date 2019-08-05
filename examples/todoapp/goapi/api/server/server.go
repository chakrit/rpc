// <auto-generated />
//
// expected import: github.com/chakrit/rpc/todo/api/server
package server

import (
	"context"
	"encoding/json"
	"net/http"

	rpc_root "github.com/chakrit/rpc/todo/api"
)

type Provider_rpc_root interface {
	Provide_rpc_root() rpc_root.Interface
}

type Result struct {
	Error   error         `json:"error"`
	Returns []interface{} `json:"returns"`
}

type Server struct {
	options  Options
	Provider Provider_rpc_root
}

type Options struct {
	Addr      string
	CtxFilter func(req *http.Request, method string) context.Context
	ErrFilter func(req *http.Request, method string, err error) error
	ErrLog    func(req *http.Request, method string, err error)
	FormatErr func(err error) string
}

func New(opts *Options) *Server {
	srv := &Server{options: *opts}
	if srv.options.CtxFilter == nil {
		srv.options.CtxFilter = func(req *http.Request, _ string) context.Context {
			return req.Context()
		}
	}
	if srv.options.ErrFilter == nil {
		srv.options.ErrFilter = func(_ *http.Request, _ string, err error) error {
			return err
		}
	}
	return srv
}

func (s *Server) Listen() error {
	return http.ListenAndServe(s.options.Addr, s.HTTPHandler())
}

func (s *Server) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	s.register_rpc_root(mux, s.Provider)
	return mux
}

func (s *Server) register_rpc_root(
	mux *http.ServeMux,
	provider Provider_rpc_root,
) *http.ServeMux {
	handler := provider.Provide_rpc_root()

	mux.HandleFunc("/api/Create", func(resp http.ResponseWriter, req *http.Request) {
		var (
			err error
			ctx context.Context
		)

		ctx = s.options.CtxFilter(req, "api/Create")
		req = req.WithContext(ctx)

		var arg0 string
		args := [1]interface{}{
			&arg0,
		}

		if req.Body != nil {
			if err := json.NewDecoder(req.Body).Decode(&args); err != nil {
				renderResult(s.options, resp, 400, &Result{
					Error:   err,
					Returns: nil,
				})
				return
			}
		}

		var (
			out0 *rpc_root.TodoItem
		)

		out0, err = handler.Create(
			ctx, arg0)

		result := &Result{}
		if err != nil {
			err = s.options.ErrFilter(req, "api/Create", err)
			s.options.ErrLog(req, "api/Create", err)
			result.Error = err
		} else {
			result.Returns = []interface{}{
				out0,
			}
		}

		renderResult(s.options, resp, 200, result)
	})

	mux.HandleFunc("/api/Destroy", func(resp http.ResponseWriter, req *http.Request) {
		var (
			err error
			ctx context.Context
		)

		ctx = s.options.CtxFilter(req, "api/Destroy")
		req = req.WithContext(ctx)

		var arg0 int64
		args := [1]interface{}{
			&arg0,
		}

		if req.Body != nil {
			if err := json.NewDecoder(req.Body).Decode(&args); err != nil {
				renderResult(s.options, resp, 400, &Result{
					Error:   err,
					Returns: nil,
				})
				return
			}
		}

		var (
			out0 *rpc_root.TodoItem
		)

		out0, err = handler.Destroy(
			ctx, arg0)

		result := &Result{}
		if err != nil {
			err = s.options.ErrFilter(req, "api/Destroy", err)
			s.options.ErrLog(req, "api/Destroy", err)
			result.Error = err
		} else {
			result.Returns = []interface{}{
				out0,
			}
		}

		renderResult(s.options, resp, 200, result)
	})

	mux.HandleFunc("/api/List", func(resp http.ResponseWriter, req *http.Request) {
		var (
			err error
			ctx context.Context
		)

		ctx = s.options.CtxFilter(req, "api/List")
		req = req.WithContext(ctx)

		var (
			out0 []*rpc_root.TodoItem
		)

		out0, err = handler.List(
			ctx)

		result := &Result{}
		if err != nil {
			err = s.options.ErrFilter(req, "api/List", err)
			s.options.ErrLog(req, "api/List", err)
			result.Error = err
		} else {
			result.Returns = []interface{}{
				out0,
			}
		}

		renderResult(s.options, resp, 200, result)
	})

	return mux
}

func renderResult(options Options, resp http.ResponseWriter, status int, result *Result) {
	resp.Header().Set("Content-Type", "application/json")

	shim := struct {
		Error   *string       `json:"error"`
		Returns []interface{} `json:"returns"`
	}{}

	if result.Error != nil {
		var errstr string
		if options.FormatErr != nil {
			errstr = options.FormatErr(result.Error)
		} else {
			errstr = result.Error.Error()
		}

		shim.Returns, shim.Error = nil, &errstr

	} else {
		shim.Returns, shim.Error = result.Returns, nil
	}

	buf, err := json.Marshal(shim)
	if err != nil {
		resp.WriteHeader(500)
		_, _ = resp.Write([]byte(`{"error":"json processing error","result":null}`))
	} else {
		resp.WriteHeader(status)
		_, _ = resp.Write(buf)
	}
}

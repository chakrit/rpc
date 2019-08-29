package main

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/chakrit/rpc/todo/api"
)

type handler struct {
	counter int64
	items   []*api.TodoItem
}

var errNotFound = errors.New("item not found")

var _ api.Interface = &handler{}

func (h *handler) List(ctx context.Context) ([]*api.TodoItem, error) {
	return h.items, nil
}

func (h *handler) Create(ctx context.Context, desc string) (*api.TodoItem, error) {
	id := atomic.AddInt64(&h.counter, 1)
	item := &api.TodoItem{
		ID:          id,
		Description: desc,
		Ctime:       time.Now(),
		Metadata:    []byte(desc),
	}

	h.items = append(h.items, item)
	return item, nil
}

func (h *handler) UpdateState(ctx context.Context, id int64, state api.State) (*api.TodoItem, error) {
	for idx, item := range h.items {
		if item.ID == id {
			item.State = state
			h.items[idx] = item
			return item, nil
		}
	}

	return nil, errNotFound
}

func (h *handler) Destroy(ctx context.Context, id int64) (*api.TodoItem, error) {
	for idx, item := range h.items {
		if item.ID == id {
			h.items = append(h.items[0:idx], h.items[idx+1:]...)
			return item, nil
		}
	}

	return nil, errNotFound
}

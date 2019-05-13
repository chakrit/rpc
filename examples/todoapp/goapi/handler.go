package main

import (
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

func (h *handler) Destroy(id int64) (*api.TodoItem, error) {
	for idx, item := range h.items {
		if item.ID == id {
			h.items = append(h.items[0:idx], h.items[idx+1:]...)
			return item, nil
		}
	}

	return nil, errNotFound
}

func (h *handler) List() ([]*api.TodoItem, error) {
	return h.items, nil
}

func (h *handler) Create(desc string) (*api.TodoItem, error) {
	id := atomic.AddInt64(&h.counter, 1)
	item := &api.TodoItem{
		ID:          id,
		Description: desc,
		Ctime:       time.Now(),
	}

	h.items = append(h.items, item)
	return item, nil
}

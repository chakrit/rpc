package main

import (
	"context"
	"errors"

	"github.com/chakrit/rpc/todo/api"
)

type handler struct {
	items []*api.TodoItem
}

var errNotFound = errors.New("item not found")

var _ api.Interface = &handler{}

func (h *handler) Destroy(ctx context.Context, id string) (*api.TodoItem, error) {
	for idx, item := range h.items {
		if item.ID == id {
			h.items = append(h.items[0:idx], h.items[idx+1:]...)
			return item, nil
		}
	}

	return nil, errNotFound
}

func (h *handler) List(ctx context.Context) ([]*api.TodoItem, error) {
	return h.items, nil
}

func (h *handler) Retrieve(ctx context.Context, id string) (*api.TodoItem, error) {
	for _, item := range h.items {
		if item.ID == id {
			return item, nil
		}
	}

	return nil, errNotFound
}

func (h *handler) Update(ctx context.Context, id string, item *api.TodoItem) (*api.TodoItem, error) {
	item.ID = id
	for idx, find := range h.items {
		if find.ID == id {
			h.items[idx] = item
			return item, nil
		}
	}

	h.items = append(h.items, item)
	return item, nil
}

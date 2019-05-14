// <auto-generated />
//
// expected import: github.com/chakrit/rpc/todo/api
package api

import (
	time "time"
)

type TodoItem struct {
	Ctime       time.Time `json:"ctime" yaml:"ctime" db:"ctime"`
	Description string    `json:"description" yaml:"description" db:"description"`
	ID          int64     `json:"id" yaml:"id" db:"id"`
}

type Interface interface {
	Create(string) (*TodoItem, error,
	)
	Destroy(int64) (*TodoItem, error,
	)
	List() ([]*TodoItem, error,
	)
}
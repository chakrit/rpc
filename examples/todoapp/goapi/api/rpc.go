// <auto-generated />
// @generated by github.com/chakrit/rpc
//
// expected import: github.com/chakrit/rpc/todo/api
package api

import (
	"context"
	"encoding/json"
	"math"

	time "time"
)

var (
	_ context.Context = nil
	_                 = json.Marshal
	_                 = math.Pi
)

type TodoItem struct {
	Ctime       time.Time `json:"ctime" yaml:"ctime" db:"ctime"`
	Description string    `json:"description" yaml:"description" db:"description"`
	ID          int64     `json:"id" yaml:"id" db:"id"`
	Metadata    []byte    `json:"metadata" yaml:"metadata" db:"metadata"`
	State       State     `json:"state" yaml:"state" db:"state"`
}

func (obj *TodoItem) MarshalJSON() ([]byte, error) {
	outobj := struct {
		Ctime       float64 `json:"ctime"`
		Description string  `json:"description"`
		ID          int64   `json:"id"`
		Metadata    []byte  `json:"metadata"`
		State       string  `json:"state"`
	}{
		Ctime: (func(t time.Time) float64 {
			sec, nsec := t.Unix(), t.Nanosecond()
			return float64(sec) + (float64(nsec) / float64(time.Second))
		})(obj.Ctime),
		Description: (obj.Description),
		ID:          (obj.ID),
		Metadata:    (obj.Metadata),
		State:       (func(v State) string { return string(v) })(obj.State),
	}
	return json.Marshal(outobj)
}

func (obj *TodoItem) UnmarshalJSON(buf []byte) error {
	inobj := struct {
		Ctime       float64 `json:"ctime"`
		Description string  `json:"description"`
		ID          int64   `json:"id"`
		Metadata    []byte  `json:"metadata"`
		State       string  `json:"state"`
	}{}

	if err := json.Unmarshal(buf, &inobj); err != nil {
		return err
	}

	obj.Ctime = (func(t float64) time.Time {
		fsec, fnsec := math.Modf(t)
		sec, nsec := int64(fsec), int64(math.Round(fnsec*float64(time.Second)))
		return time.Unix(sec, nsec)
	})(inobj.Ctime)
	obj.Description = (inobj.Description)
	obj.ID = (inobj.ID)
	obj.Metadata = (inobj.Metadata)
	obj.State = (func(v string) State { return State(v) })(inobj.State)
	return nil
}

type State string

const (
	StateNew        = State("new")
	StateInProgress = State("in-progress")
	StateOverdue    = State("overdue")
	StateCompleted  = State("completed")
)

type Interface interface {
	Create(context.Context, string) (*TodoItem, error,
	)
	Destroy(context.Context, int64) (*TodoItem, error,
	)
	List(context.Context) ([]*TodoItem, error,
	)
	UpdateState(context.Context, int64, State) (*TodoItem, error,
	)
}

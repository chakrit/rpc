// <auto-generated />
//
// expected import: github.com/chakrit/rpc/examples/todo
package Todo

import (
	"encoding/json"
	"math"

	time "time"
)

var _ = math.Pi
var _ = json.Marshal

type User struct {
	Ctime    time.Time `json:"ctime" db:"ctime"`
	Username string    `json:"username" db:"username"`
}

func (obj *User) MarshalJSON() ([]byte, error) {
	outobj := struct {
		Ctime    float64 `json:"ctime"`
		Username string  `json:"username"`
	}{
		Ctime: (func(t time.Time) float64 {
			sec, nsec := t.Unix(), t.Nanosecond()
			return float64(sec) + (float64(nsec) / float64(time.Second))
		})(obj.Ctime),
		Username: (obj.Username),
	}
	return json.Marshal(outobj)
}

func (obj *User) UnmarshalJSON(buf []byte) error {
	inobj := struct {
		Ctime    float64 `json:"ctime"`
		Username string  `json:"username"`
	}{}

	if err := json.Unmarshal(buf, &inobj); err != nil {
		return err
	}

	obj.Ctime = (func(t float64) time.Time {
		fsec, fnsec := math.Modf(t)
		sec, nsec := int64(fsec), int64(math.Round(fnsec*float64(time.Second)))
		return time.Unix(sec, nsec)
	})(inobj.Ctime)
	obj.Username = (inobj.Username)
	return nil
}

type Interface interface {
}

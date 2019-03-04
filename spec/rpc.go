package spec

type RPC struct {
	Name        string     `json:"name"`
	InputTypes  []*TypeRef `json:"input"`
	OutputTypes []*TypeRef `json:"output"`
}

var _ Node = &RPC{}

func (r *RPC) name() string { return r.Name }
func (r *RPC) node()        {}

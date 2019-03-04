package spec

type TypeRef struct {
	Name      string     `json:"name"`
	Arguments []*TypeRef `json:"arguments"`
}

var _ Node = &TypeRef{}

func (t *TypeRef) name() string { return t.Name }
func (t *TypeRef) node()        {}

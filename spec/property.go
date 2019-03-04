package spec

type Property struct {
	Name string   `json:"name"`
	Type *TypeRef `json:"type"`
}

var _ Node = &Property{}

func (p *Property) name() string { return p.Name }
func (p *Property) node()        {}

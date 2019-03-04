package spec

type Type struct {
	Name       string   `json:"name"`
	Properties Mappings `json:"properties"`
}

var _ Node = &Type{}
var _ merger = &Type{}

func (t *Type) name() string { return t.Name }
func (t *Type) node()        {}

func (t *Type) Merge(node Node) Node {
	another, ok := node.(*Type)
	if !ok { // TODO: Warn about this
		return another
	}

	for name, prop := range another.Properties {
		t.Properties[name] = prop
	}
	return t
}

package spec

type Namespace struct {
	Name     string                 `json:"name"`
	Children Mappings               `json:"children"`
	Options  map[string]interface{} `json:"options"`

	Types Mappings `json:"types"`
	Enums Mappings `json:"enums"`
	RPCs  Mappings `json:"rpcs"`
}

var _ Node = &Namespace{}
var _ merger = &Namespace{}

func (ns *Namespace) name() string { return ns.Name }
func (ns *Namespace) node()        {}

func (ns *Namespace) Merge(node Node) Node {
	another, ok := node.(*Namespace)
	if !ok { // TODO: Warn about this.
		return ns
	}

	if ns.Name == "" {
		ns.Name = another.Name
	}
	if ns.Options == nil && len(another.Options) > 0 {
		ns.Options = map[string]interface{}{}
	}

	for _, child := range another.Children {
		ns.Children.Add(child)
	}
	for _, typ := range another.Types {
		ns.Types.Add(typ)
	}
	for _, enum := range another.Enums {
		ns.Enums.Add(enum)
	}
	for _, rpc := range another.RPCs {
		ns.RPCs.AddIfNew(rpc)
	}
	for key, value := range another.Options {
		ns.Options[key] = value
	}

	return ns
}

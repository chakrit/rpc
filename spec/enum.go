package spec

type Enum struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

var _ Node = &Enum{}
var _ merger = &Enum{}

func (e Enum) name() string { return e.Name }
func (e Enum) node()        {}

func (e Enum) Merge(node Node) Node {
	another, ok := node.(*Enum)
	if !ok { // TODO: Warn
		return another
	}

	existing := map[string]struct{}{}
	for _, member := range e.Members {
		existing[member] = struct{}{}
	}
	for _, member := range another.Members {
		if _, exists := existing[member]; !exists {
			e.Members = append(e.Members, member)
		}
	}
	return e
}

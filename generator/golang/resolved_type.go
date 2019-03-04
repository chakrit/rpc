package golang

import "github.com/chakrit/rpc/spec"

const (
	TypeSimple = 1 << iota
	TypeTime
	TypeMap
	TypeList
	TypeUserDefined
	TypeIsBuiltin
	TypeIsLocal
	TypeIsImported
)

type ResolvedType struct {
	Name  string
	Pkg   *Pkg
	Type  *spec.Type
	Args  []*ResolvedType
	Flags int
}

func (rt *ResolvedType) IsSimple() bool      { return (rt.Flags & TypeSimple) != 0 }
func (rt *ResolvedType) IsTime() bool        { return (rt.Flags & TypeTime) != 0 }
func (rt *ResolvedType) IsMap() bool         { return (rt.Flags & TypeMap) != 0 }
func (rt *ResolvedType) IsList() bool        { return (rt.Flags & TypeList) != 0 }
func (rt *ResolvedType) IsUserDefined() bool { return (rt.Flags & TypeUserDefined) != 0 }

func (rt *ResolvedType) IsBuiltin() bool  { return (rt.Flags & TypeIsBuiltin) != 0 }
func (rt *ResolvedType) IsLocal() bool    { return (rt.Flags & TypeIsLocal) != 0 }
func (rt *ResolvedType) IsImported() bool { return (rt.Flags & TypeIsImported) != 0 }

func (rt *ResolvedType) WithoutFlags(flags int) *ResolvedType {
	clone := *rt
	clone.Flags = clone.Flags & ^flags
	clone.Args = make([]*ResolvedType, len(rt.Args))
	for idx := range clone.Args {
		clone.Args[idx] = rt.Args[idx].WithoutFlags(flags)
	}

	return &clone
}

func (rt *ResolvedType) WithFlags(flags int) *ResolvedType {
	clone := *rt
	clone.Flags = clone.Flags | flags
	clone.Args = make([]*ResolvedType, len(rt.Args))
	for idx := range clone.Args {
		clone.Args[idx] = rt.Args[idx].WithFlags(flags)
	}

	return &clone
}

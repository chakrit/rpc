package golang

import (
	"github.com/chakrit/rpc/spec"
)

const (
	TypeSimple = 1 << iota
	TypeTime
	TypeBytes
	TypeMap
	TypeList
	TypeUserDefined
	OriginBuiltin
	OriginPkgLocal
	OriginImported
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
func (rt *ResolvedType) IsBytes() bool       { return (rt.Flags & TypeBytes) != 0 }
func (rt *ResolvedType) IsMap() bool         { return (rt.Flags & TypeMap) != 0 }
func (rt *ResolvedType) IsList() bool        { return (rt.Flags & TypeList) != 0 }
func (rt *ResolvedType) IsUserDefined() bool { return (rt.Flags & TypeUserDefined) != 0 }

func (rt *ResolvedType) IsBuiltin() bool  { return (rt.Flags & OriginBuiltin) != 0 }
func (rt *ResolvedType) IsPkgLocal() bool { return (rt.Flags & OriginPkgLocal) != 0 }
func (rt *ResolvedType) IsImported() bool { return (rt.Flags & OriginImported) != 0 }

func (rt *ResolvedType) Arg(n int) *ResolvedType {
	if len(rt.Args) > n {
		return rt.Args[n]
	} else {
		return nil
	}
}

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

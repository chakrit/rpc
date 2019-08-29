package elm

import (
	"text/template"
)

func funcMap() template.FuncMap {
	f := template.FuncMap{}
	f["ifFirst"] = func(idx int, firstStr, restStr string) string {
		if idx == 0 {
			return firstStr
		} else {
			return restStr
		}
	}

	f["resolve"] = func(ref *TypeRef) *TypeResolution {
		return ref.Module.Registry.Resolve(ref)
	}

	return f
}

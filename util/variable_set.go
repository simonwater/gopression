package util

import (
	"sort"
	"strings"
)

type VariableSet struct {
	assigns map[string]bool
	depends map[string]bool
}

func NewVariableSet() *VariableSet {
	return &VariableSet{
		assigns: make(map[string]bool),
		depends: make(map[string]bool),
	}
}

func (vs *VariableSet) GetAssigns() map[string]bool {
	return vs.assigns
}

func (vs *VariableSet) SetAssigns(assigns map[string]bool) {
	vs.assigns = make(map[string]bool)
	for k := range assigns {
		vs.assigns[k] = true
	}
}

func (vs *VariableSet) GetDepends() map[string]bool {
	return vs.depends
}

func (vs *VariableSet) SetDepends(depends map[string]bool) {
	vs.depends = make(map[string]bool)
	for k := range depends {
		vs.depends[k] = true
	}
}

func (vs *VariableSet) AddAssign(name string) {
	vs.assigns[name] = true
}

func (vs *VariableSet) AddDepend(name string) {
	vs.depends[name] = true
}

func (vs *VariableSet) Combine(other *VariableSet) {
	if other == nil {
		return
	}
	for k := range other.assigns {
		vs.assigns[k] = true
	}
	for k := range other.depends {
		vs.depends[k] = true
	}
}

func FromDepends(names ...string) *VariableSet {
	vs := NewVariableSet()
	for _, name := range names {
		vs.depends[name] = true
	}
	return vs
}

func FromAssigns(names ...string) *VariableSet {
	vs := NewVariableSet()
	for _, name := range names {
		vs.assigns[name] = true
	}
	return vs
}

func (vs *VariableSet) String() string {
	assigns := make([]string, 0, len(vs.assigns))
	for k := range vs.assigns {
		assigns = append(assigns, k)
	}
	sort.Strings(assigns)
	assignStr := strings.Join(assigns, ",")

	depends := make([]string, 0, len(vs.depends))
	for k := range vs.depends {
		depends = append(depends, k)
	}
	sort.Strings(depends)
	dependStr := strings.Join(depends, ",")

	if assignStr != "" {
		return assignStr + " = " + dependStr
	}
	return dependStr
}

package util

import (
	"strings"
)

type Field struct {
	Name  string
	Owner *Field
	src   string
}

func NewField(name string) *Field {
	return &Field{
		Name:  name,
		Owner: nil,
	}
}

func NewFieldWithOwner(name string, owner *Field) *Field {
	return &Field{
		Name:  name,
		Owner: owner,
	}
}

func NewFieldFromPath(src string) *Field {
	names := strings.Split(src, ".")
	var cur *Field
	for _, name := range names {
		cur = NewFieldWithOwner(name, cur)
	}
	return cur
}

func (f *Field) GetName() string {
	return f.Name
}

func (f *Field) GetOwner() *Field {
	return f.Owner
}

func (f *Field) search(field *Field, path *[]string) {
	if field == nil {
		return
	}
	f.search(field.Owner, path)
	*path = append(*path, field.Name)
}

func (f *Field) String() string {
	if IsEmpty(f.src) {
		path := []string{}
		f.search(f, &path)
		f.src = strings.Join(path, ".")
	}
	return f.src
}

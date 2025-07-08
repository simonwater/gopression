package functions

type Function struct {
	Name  string
	Title string
	Group string
}

func NewFunction(name, title, group string) *Function {
	return &Function{
		Name:  name,
		Title: title,
		Group: group,
	}
}

func (f *Function) GetName() string {
	return f.Name
}

func (f *Function) GetTitle() string {
	return f.Title
}

func (f *Function) GetGroup() string {
	return f.Group
}

package datastore

// Filter ...
type Filter struct {
	field, value, operator string
}

// NewFilter returns a new member instance
func NewFilter(f, o, v string) *Filter {
	return &Filter{field: f, operator: o, value: v}
}

// SQL ...
func (f *Filter) SQL() string {
	return "(" + f.field + " " + f.operator + " '" + f.value + "')"
}

func (f *Filter) String() string {
	return "<" + f.SQL() + ">"
}

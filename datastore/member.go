package datastore

import "strconv"

// Member - encapsulation of a member
type Member struct {
	ID   int64
	Name string
	Age  int16
}

func (m *Member) String() string {
	return "<" + strconv.FormatInt(m.ID, 10) + ":" + m.Name + ">"
}

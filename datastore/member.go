package datastore

import (
	"strconv"
	"time"
)

// Member - encapsulation of a member
type Member struct {
	ID                  int64
	FirstName, LastName string
	Birthday            time.Time
	IsMale              bool
}

// NewMemberWithoutID returns a new member instance
func NewMemberWithoutID(fname, lname string, bd time.Time, ism bool) *Member {
	return &Member{FirstName: fname, LastName: lname, Birthday: bd, IsMale: ism}
}

// NewMember returns a new member instance
func NewMember(id int64, fname, lname string, bd time.Time, ism bool) *Member {
	m := NewMemberWithoutID(fname, lname, bd, ism)
	m.ID = id
	return m
}

func (m *Member) String() string {
	return "<" + strconv.FormatInt(m.ID, 10) + ":" + m.FirstName + ">"
}

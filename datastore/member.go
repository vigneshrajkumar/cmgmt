package datastore

import (
	"strconv"
	"time"
)

// Member - encapsulation of a member
type Member struct {
	ID        int64     `json:"ID"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Birthday  time.Time `json:"birthday"`
	Gender    string    `json:"gender"`
	FamilyID  int64     `json:"fID"`
}

// NewMemberWithoutID returns a new member instance
func NewMemberWithoutID(fname, lname string, bd time.Time, gender string, fID int64) *Member {
	return &Member{FirstName: fname, LastName: lname, Birthday: bd, Gender: gender, FamilyID: fID}
}

// NewMember returns a new member instance
func NewMember(id int64, fname, lname string, bd time.Time, gender string, fID int64) *Member {
	m := NewMemberWithoutID(fname, lname, bd, gender, fID)
	m.ID = id
	return m
}

func (m *Member) String() string {
	return "<" + strconv.FormatInt(m.ID, 10) + ":" + m.FirstName + ">"
}

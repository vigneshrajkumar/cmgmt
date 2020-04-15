package datastore

import (
	"strconv"
)

// Family - encapsulation of a Family
type Family struct {
	ID     int64
	HeadID int64
}

// NewFamilyWithoutID returns a new Family instance
func NewFamilyWithoutID(hID int64) *Family {
	return &Family{HeadID: hID}
}

// NewFamily returns a new Family instance
func NewFamily(id, hID int64) *Family {
	f := NewFamilyWithoutID(hID)
	f.ID = id
	return f
}

func (m *Family) String() string {
	return "<" + strconv.FormatInt(m.ID, 10) + ":H" + strconv.FormatInt(m.HeadID, 10) + ">"
}

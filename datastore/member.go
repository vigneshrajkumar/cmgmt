package datastore

import "fmt"

// Member - encapsulation of a member
type Member struct {
	ID float64 `json:"ID"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Phone string `json:"phone"`
	Home  string `json:"home"`

	Email string `json:"email"`

	DateOfBirth string `json:"dateOfBirth"`
	Gender      string `json:"gender"`

	FamilyID float64 `json:"fID"`

	Address string `json:"address"`
	Pincode string `json:"pincode"`

	BloodGroup   string  `json:"bloodGroup"`
	ProfessionID float64 `json:"professionID"`

	Photo   string `json:"photo"`
	Remarks string `json:"remarks"`
}

// // NewMemberWithoutID returns a new member instance
// func NewMemberWithoutID(fname, lname string, bd time.Time, gender string, fID int64) *Member {
// 	return &Member{FirstName: fname, LastName: lname, Birthday: bd, Gender: gender, FamilyID: fID}
// }

// // NewMember returns a new member instance
// func NewMember(id int64, fname, lname string, bd time.Time, gender string, fID int64) *Member {
// 	m := NewMemberWithoutID(fname, lname, bd, gender, fID)
// 	m.ID = id
// 	return m
// }

func (m *Member) String() string {
	return "<" + fmt.Sprintf("%f", m.ID) + ":" + m.FirstName + " " + m.LastName + ">"
}

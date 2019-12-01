package datastore

// MembersData - Data carriers
type MembersData struct {
	Username string
	Members  []*Member
}

// ErrorData - Data carriers
type ErrorData struct {
	ErrorMessage string
}

// MemberData - Data carriers
type MemberData struct {
	ID       int64
	Username string
	Name     string
	Age      int16
}

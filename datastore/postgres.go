package datastore

import (
	"fmt"
	"strings"
)

func getSQLToFetchProfessions() string {
	return strings.Join([]string{"SELECT ", ID, ",", Description, " FROM cm.profession"}, "")
}

func getSQLToInsertMember(m *Member) string {
	colNames := []string{ID, FirstName, LastName, Phone, Home, Email, DateOfBirth, Gender, FamilyID, Address, Pincode, BloodGroup, ProfessionID, Photo, Remarks}
	colValues := []string{fmt.Sprintf("%.0f", m.ID), m.FirstName, m.LastName, m.Phone, m.Home, m.Email, m.DateOfBirth, m.Gender, fmt.Sprintf("%.0f", m.FamilyID), m.Address, m.Pincode, m.BloodGroup, fmt.Sprintf("%.0f", m.ProfessionID), m.Photo, m.Remarks}

	for ix, cv := range colValues {
		colValues[ix] = strings.Join([]string{"'", cv, "'"}, "")
	}
	return strings.Join([]string{"INSERT INTO cm.member (", strings.Join(colNames, ", "), ") VALUES (", strings.Join(colValues, ", "), ")"}, "")
}

func getSQLToFetchMembersOverview(fitlers ...*Filter) string {
	// explicitly casting bigint _id as float because _id in Member is float64
	colNames := []string{ID + "::float8", FirstName, LastName, Phone, Email}
	query := strings.Join([]string{"SELECT ", strings.Join(colNames, ", "), " FROM cm.member"}, "")
	if len(fitlers) > 0 {
		whereClause := make([]string, len(fitlers))
		for ix, f := range fitlers {
			whereClause[ix] = f.SQL()
		}
		query = strings.Join([]string{query, " WHERE ", strings.Join(whereClause, " AND ")}, "")
	}
	return query
}

func getSQLToFetchMembersByID(id float64) string {
	selCols := []string{ID + "::float8", FirstName, LastName, Phone, Home, Email, DateOfBirth + "::varchar", Gender, Address, Pincode, BloodGroup, ProfessionID + "::float8", Remarks}
	return strings.Join([]string{"SELECT ", strings.Join(selCols, ", "), " FROM cm.member WHERE ", ID, " = ", fmt.Sprintf("%.0f", id)}, "")
}

func getSQLToResolveProfessionID(id float64) string {
	return strings.Join([]string{"SELECT ", Description, " FROM cm.profession WHERE ", ID, " = ", fmt.Sprintf("%.0f", id)}, "")
}

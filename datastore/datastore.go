package datastore

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Store is the interface which lets you interact with you rdbms
type Store struct {
	cxn *pgxpool.Pool
	gen int64
}

// NextID - inhouse id gen
func (d *Store) NextID() int64 {
	d.gen++
	return d.gen
}

// NewStore - creates a new Store instance
func NewStore() (*Store, error) {
	connectionString := "user=" + os.Getenv("USER") + " password=" + os.Getenv("PASSWORD") + " host=" + os.Getenv("HOST") + " port=" + os.Getenv("PORT") + " dbname=" + os.Getenv("DATABASE")
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.ConnectConfig(context.TODO(), config)
	if err != nil {
		return nil, err
	}
	return &Store{cxn: pool}, nil
}

// Initialize - preliminary setup for DB access
func (d *Store) Initialize() error {
	queries := []string{
		"CREATE SCHEMA cm",
		"CREATE TYPE cm.membership_status AS ENUM ('active', 'prep', 'inactive')",
		"CREATE TYPE cm.blood_group AS ENUM ('A+', 'A-', 'B+', 'B-', 'O+', 'O-', 'AB+', 'AB-')",
		"CREATE TYPE cm.family_role AS ENUM ('spouse', 'children', 'sibling', 'parent')",
		"CREATE TYPE cm.txn_type AS ENUM ('tithe', 'donation')",
		"CREATE TABLE IF NOT EXISTS cm.user (_id BIGSERIAL, username VARCHAR, password VARCHAR, session_token VARCHAR)",

		//To Add: Blood Group, date of baptism, date of confirmation, membership status, profession, remarks, photo
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS cm.member (%s BIGINT, %s VARCHAR, %s VARCHAR, %s date, %s VARCHAR, %s VARCHAR, %s VARCHAR, %s VARCHAR, %s VARCHAR, %s VARCHAR, %s BIGINT, %s VARCHAR, %s BIGINT, %s cm.blood_group, %s BYTEA)", ID, FirstName, LastName, DateOfBirth, Gender, Phone, Home, Email, Address, Pincode, FamilyID, Remarks, ProfessionID, BloodGroup, Photo),

		"CREATE TABLE IF NOT EXISTS cm.family (_id BIGINT, family_head BIGINT)",
		"CREATE TABLE IF NOT EXISTS cm.transactions (_id BIGINT, amount numeric, date date, type cm.txn_type, remarks VARCHAR)",
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS cm.profession (%s BIGINT, %s VARCHAR)", ID, Description),

		// Default Values
		fmt.Sprintf("INSERT INTO cm.profession(%s, %s) VALUES (%s, 'Doctor'), (%s, 'Engineer'), (%s, 'Architect'), (%s, 'Other')", ID, Description, strconv.FormatInt(d.NextID(), 10), strconv.FormatInt(d.NextID(), 10), strconv.FormatInt(d.NextID(), 10), strconv.FormatInt(d.NextID(), 10)),
	}
	for _, q := range queries {
		_, err := d.cxn.Exec(context.TODO(), q)
		if err != nil {
			fmt.Println("err: ", q)
			return err
		}
	}
	return nil
}

// Reset - drops everything
func (d *Store) Reset() error {
	_, err := d.cxn.Exec(context.TODO(), "DROP SCHEMA cm CASCADE")
	if err != nil {
		return err
	}
	return nil
}

// EstablishAdminAccess if not found, will add access for admin account
func (d *Store) EstablishAdminAccess() error {
	exists, err := d.CheckUserExistance("nate")
	if err != nil {
		return err
	}
	if !exists {
		err = d.AddUser("nate", "robinson")
		if err != nil {
			fmt.Println("error adding user")
			return err
		}
		fmt.Println("Established Admin Access")
	}
	return nil
}

// AddUser adds access to new user
func (d *Store) AddUser(username, password string) error {
	_, err := d.cxn.Exec(context.TODO(), "INSERT INTO cm.user (username, password) VALUES ('"+username+"', '"+password+"')")
	if err != nil {
		return err
	}
	return nil
}

// AddMember adds access to new user
func (d *Store) AddMember(m *Member) error {
	log.Println("AddMember()", m)
	// .Format("2006-01-02 15:04:05")

	stmt := getSQLToInsertMember(m)

	log.Println("Q::", stmt)
	_, err := d.cxn.Exec(context.TODO(), stmt)
	if err != nil {
		return err
	}
	return nil
}

// AddFamily adds a new family
func (d *Store) AddFamily(f *Family) error {
	log.Println("AddFamily()", f)
	query := fmt.Sprintf("INSERT INTO cm.family (_id, family_head) VALUES ('%s','%s')", strconv.FormatInt(f.ID, 10), strconv.FormatInt(f.HeadID, 10))
	log.Println("Q::", query)
	_, err := d.cxn.Exec(context.TODO(), query)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMember deletes a member with given ID
func (d *Store) DeleteMember(id int64) error {
	log.Println("DeleteMember()", id)
	query := fmt.Sprintf("DELETE FROM cm.member WHERE _id IN  ( %s ) ", strconv.FormatInt(id, 10))
	log.Println("Q::", query)
	_, err := d.cxn.Exec(context.TODO(), query)
	if err != nil {
		return err
	}
	return nil
}

// GetMembers retreives all members from DB
func (d *Store) GetMembers(fitlers ...*Filter) (mems []*Member, err error) {
	log.Println("GetMembers()", fitlers)
	query := getSQLToFetchMembersOverview(fitlers...)
	log.Println("Q::", query)
	rows, err := d.cxn.Query(context.TODO(), query)
	if err != nil {
		log.Println("ERROR ::", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.ID, &m.FirstName, &m.LastName, &m.Phone, &m.Email); err != nil {
			return nil, err
		}
		mems = append(mems, &m)
	}
	return
}

// GetFamilies retreives all families from DB
func (d *Store) GetFamilies() (fams []*Family, err error) {
	log.Println("GetFamilies()")
	query := "SELECT _id, family_head FROM cm.family"
	log.Println("Q::", query)
	rows, err := d.cxn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		fams = append(fams, NewFamily(vals[0].(int64), vals[1].(int64)))
	}
	return
}

// GetMemberByID retreives a member by ID
func (d *Store) GetMemberByID(id float64) (*Member, error) {
	query := getSQLToFetchMembersByID(id)
	log.Println("Q: ", query)
	rows, err := d.cxn.Query(context.TODO(), query)
	if err != nil {
		log.Println("ERR: ", err)
		return nil, err
	}
	defer rows.Close()
	var mem Member

	for rows.Next() {
		if err := rows.Scan(&mem.ID, &mem.FirstName, &mem.LastName, &mem.Phone, &mem.Home, &mem.Email, &mem.DateOfBirth, &mem.Gender, &mem.Address, &mem.Pincode, &mem.BloodGroup, &mem.ProfessionID, &mem.Remarks); err != nil {
			return nil, err
		}
	}
	return &mem, nil
}

// GetFamilyByID retreives a member by ID
func (d *Store) GetFamilyByID(id int64) (fam *Family, err error) {
	row, err := d.cxn.Query(context.TODO(), "SELECT _id, family_head FROM cm.family where _id = "+strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		vals, err := row.Values()
		if err != nil {
			return nil, err
		}
		fam = NewFamily(vals[0].(int64), vals[1].(int64))
	}
	return
}

// ValidateUser check whether a given user exists
func (d *Store) ValidateUser(username, password string) (bool, error) {
	var exists bool
	row := d.cxn.QueryRow(context.TODO(), "SELECT EXISTS (SELECT _id FROM cm.user WHERE username = '"+username+"' AND password = '"+password+"')")
	err := row.Scan(&exists)
	if err != nil {
		return exists, err
	}
	return exists, nil
}

// CheckUserExistance check whether a given user exists
func (d *Store) CheckUserExistance(username string) (bool, error) {
	var exists bool
	row := d.cxn.QueryRow(context.TODO(), "SELECT EXISTS (SELECT _id FROM cm.user WHERE username = '"+username+"')")
	err := row.Scan(&exists)
	if err != nil {
		return exists, err
	}
	return exists, nil
}

// UpdateSessionToken updates session token for the given user
func (d *Store) UpdateSessionToken(username, token string) error {
	_, err := d.cxn.Exec(context.TODO(), "UPDATE cm.user SET session_token = '"+token+"' WHERE username = '"+username+"'")
	if err != nil {
		fmt.Println("err UpdateSessionToken")
		return err
	}
	return nil
}

// GetUser returns username from session token
func (d *Store) GetUser(token string) (username string, err error) {
	rows, err := d.cxn.Query(context.TODO(), "SELECT username FROM cm.user WHERE session_token = '"+token+"'")
	if err != nil {
		fmt.Println("err GetUser")
		return "nil", err
	}
	defer rows.Close()

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			fmt.Println("err GetUser")
			return "nil", err
		}
		username = vals[0].(string)
	}
	return
}

// GetProfessionsWithID a map of professions and respective ID
func (d *Store) GetProfessionsWithID() (map[int64]string, error) {
	log.Println("GetProfessionsWithID()")
	mapping := make(map[int64]string)
	query := getSQLToFetchProfessions()
	log.Println("Q::", query)
	rows, err := d.cxn.Query(context.TODO(), query)
	if err != nil {
		log.Println("Error while executing the query: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		log.Println("row vals:", vals)
		mapping[vals[0].(int64)] = vals[1].(string)
	}
	log.Println("retruning", mapping)
	return mapping, nil
}

// GetFamilyNamesWithID a map of famaily name and respective ID
func (d *Store) GetFamilyNamesWithID() (map[int64]string, error) {
	log.Println("GetFamilyNamesWithID()")
	mapping := make(map[int64]string)
	query := "select A._id as id, B.first_name || ' ' || B.last_name as name from cm.family A inner join cm.member B on A._id = B.family_id;"
	log.Println("Q::", query)
	rows, err := d.cxn.Query(context.TODO(), query)
	if err != nil {
		log.Println("ERROR::", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		log.Println(vals)
		mapping[vals[0].(int64)] = vals[1].(string)
	}
	log.Println("returning ", mapping)
	return mapping, nil
}

// GetMemberFamilyInfo returns the given member's family info
func (d *Store) GetMemberFamilyInfo(mID string) (map[int64][]interface{}, error) {
	mapping := make(map[int64][]interface{})
	query := strings.Join([]string{"select C._id as id, C.first_name || ' ' || C.last_name as name, CASE WHEN C._id = A.family_head THEN true ELSE false END as is_head from cm.family A join cm.member B on A._id = B.family_id inner join cm.member C on A._id = C.family_id where B._id = ", mID}, " ")
	log.Println("Q::", query)
	rows, err := d.cxn.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		fmt.Println(vals)
		mapping[vals[0].(int64)] = []interface{}{vals[1].(string), vals[2].(bool)}
	}
	log.Println("returning::", mapping)
	return mapping, nil
}

// ResolveProfession retreives the profession description
func (d *Store) ResolveProfession(id float64) (desc string, err error) {
	query := getSQLToResolveProfessionID(id)
	log.Println("Q: ", query)
	rows, err := d.cxn.Query(context.TODO(), query)
	if err != nil {
		log.Println("ERR: ", err)
		return "resolve error", err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&desc); err != nil {
			return "resolve error", err
		}
	}
	return
}

package datastore

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
)

// Store is the interface which lets you interact with you rdbms
type Store struct {
	cxn *pgx.Conn
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
	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}
	conn, err := pgx.ConnectConfig(context.TODO(), config)
	if err != nil {
		return nil, err
	}
	return &Store{cxn: conn}, nil
}

// Initialize - preliminary setup for DB access
func (d *Store) Initialize() error {
	queries := []string{
		"CREATE SCHEMA cm",
		"CREATE TYPE cm.membership_status AS ENUM ('active', 'prep', 'inactive')",
		"CREATE TYPE cm.family_role AS ENUM ('spouse', 'children', 'sibling', 'parent')",
		"CREATE TYPE cm.txn_type AS ENUM ('tithe', 'donation')",
		"CREATE TABLE IF NOT EXISTS cm.user (_id BIGSERIAL, username VARCHAR, password VARCHAR, session_token VARCHAR)",
		"CREATE TABLE IF NOT EXISTS cm.member (_id BIGINT, first_name VARCHAR, last_name VARCHAR,  date_of_birth date, gender VARCHAR, family_id BIGINT)",
		// "CREATE TABLE IF NOT EXISTS cm.member (_id BIGINT, first_name VARCHAR, middle_name VARCHAR, last_name VARCHAR, phone VARCHAR, home VARCHAR, email VARCHAR, address VARCHAR, pincode VARCHAR, date_of_birth date, date_of_baptism date, date_of_confirmation date, blood_group VARCHAR, profession SMALLINT, remarks VARCHAR, photo BYTEA, status cm.membership_status, family_id bigint)",
		"CREATE TABLE IF NOT EXISTS cm.family (_id BIGINT, family_head BIGINT)",
		"CREATE TABLE IF NOT EXISTS cm.transactions (_id BIGINT, amount numeric, date date, type cm.txn_type, remarks VARCHAR)",
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
		fmt.Println("error checking user existance")
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
	query := fmt.Sprintf("INSERT INTO cm.member (_id, first_name, last_name,  date_of_birth, gender, family_id) VALUES ('%s','%s', '%s', '%s', '%s', '%s')", strconv.FormatInt(m.ID, 10), m.FirstName, m.LastName, m.Birthday.Format("2006-01-02 15:04:05"), m.Gender, strconv.FormatInt(m.FamilyID, 10))
	fmt.Println(query)
	_, err := d.cxn.Exec(context.TODO(), query)
	if err != nil {
		return err
	}
	return nil
}

// AddFamily adds a new family
func (d *Store) AddFamily(f *Family) error {
	query := fmt.Sprintf("INSERT INTO cm.family (_id, family_head) VALUES ('%s','%s')", strconv.FormatInt(f.ID, 10), strconv.FormatInt(f.HeadID, 10))
	fmt.Println(query)
	_, err := d.cxn.Exec(context.TODO(), query)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMember deletes a member with given ID
func (d *Store) DeleteMember(id int64) error {
	query := fmt.Sprintf("DELETE FROM cm.member WHERE _id IN  ( %s ) ", strconv.FormatInt(id, 10))
	fmt.Println(query)
	_, err := d.cxn.Exec(context.TODO(), query)
	if err != nil {
		return err
	}
	return nil
}

// GetMembers retreives all members from DB
func (d *Store) GetMembers() (mems []*Member, err error) {
	rows, err := d.cxn.Query(context.TODO(), "SELECT _id, first_name, last_name,  date_of_birth, gender, family_id FROM cm.member")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		mems = append(mems, NewMember(vals[0].(int64), vals[1].(string), vals[2].(string), vals[3].(time.Time), vals[4].(string), vals[5].(int64)))
	}
	return
}

// GetFamilies retreives all families from DB
func (d *Store) GetFamilies() (fams []*Family, err error) {
	rows, err := d.cxn.Query(context.TODO(), "SELECT _id, family_head FROM cm.family")
	if err != nil {
		return nil, err
	}
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
func (d *Store) GetMemberByID(id int64) (mem *Member, err error) {
	row, err := d.cxn.Query(context.TODO(), "SELECT _id, first_name, last_name,  date_of_birth, gender, family_id FROM cm.member where _id = "+strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		vals, err := row.Values()
		if err != nil {
			return nil, err
		}
		mem = NewMember(vals[0].(int64), vals[1].(string), vals[2].(string), vals[3].(time.Time), vals[4].(string), vals[5].(int64))
	}
	return
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

// GetFamilyNamesWithID a map of famaily name and respective ID
func (d *Store) GetFamilyNamesWithID() (mapping map[int64]string, err error) {
	mapping = make(map[int64]string)
	rows, err := d.cxn.Query(context.TODO(), "select A._id as id, B.first_name || ' ' || B.last_name as name from cm.family A inner join cm.member B on A._id = B.family_id;")
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
		mapping[vals[0].(int64)] = vals[1].(string)
	}
	return
}

// GetMemberFamilyInfo returns the given member's family info
func (d *Store) GetMemberFamilyInfo(mID string) (mapping map[int64][]interface{}, err error) {
	mapping = make(map[int64][]interface{})
	rows, err := d.cxn.Query(context.TODO(), "select C._id as id, C.first_name || ' ' || C.last_name as name, CASE WHEN C._id = A.family_head THEN true ELSE false END as is_head from cm.family A join cm.member B on A._id = B.family_id inner join cm.member C on A._id = C.family_id where B._id = "+mID+";")
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
	return
}

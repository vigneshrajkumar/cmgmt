package datastore

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

// Store is the interface which lets you interact with you rdbms
type Store struct {
	cxn *pgx.Conn
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
		"CREATE TABLE IF NOT EXISTS cm.member (_id BIGSERIAL, first_name VARCHAR, last_name VARCHAR,  date_of_birth date, is_male BOOLEAN)",
		// "CREATE TABLE IF NOT EXISTS cm.member (_id BIGSERIAL, first_name VARCHAR, middle_name VARCHAR, last_name VARCHAR, phone VARCHAR, home VARCHAR, email VARCHAR, address VARCHAR, pincode VARCHAR, date_of_birth date, date_of_baptism date, date_of_confirmation date, blood_group VARCHAR, profession SMALLINT, remarks VARCHAR, photo BYTEA, status cm.membership_status, family_id bigint)",
		"CREATE TABLE IF NOT EXISTS cm.family (_id BIGSERIAL, date_of_anniversary date, role cm.family_role)",
		"CREATE TABLE IF NOT EXISTS cm.transactions (_id BIGSERIAL, amount numeric, date date, type cm.txn_type, remarks VARCHAR)",
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
	isMale := "False"
	if m.IsMale {
		isMale = "True"
	}
	query := fmt.Sprintf("INSERT INTO cm.member (first_name, last_name,  date_of_birth, is_male) VALUES ('%s', '%s', '%s', '%s')", m.FirstName, m.LastName, m.Birthday.Format("2006-01-02 15:04:05"), isMale)
	fmt.Println(query)
	_, err := d.cxn.Exec(context.TODO(), query)
	if err != nil {
		return err
	}
	return nil
}

// GetMembers retreives all members from DB
func (d *Store) GetMembers() (mems []*Member, err error) {
	rows, err := d.cxn.Query(context.TODO(), "SELECT _id, first_name, last_name,  date_of_birth, is_male FROM cm.member")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		mems = append(mems, NewMember(vals[0].(int64), vals[1].(string), vals[2].(string), vals[3].(time.Time), vals[4].(bool)))
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

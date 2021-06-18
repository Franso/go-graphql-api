package postgres

import (
	"database/sql"
	"fmt"
)

// Db Structure
// Db is our database struct used for interacting with the db
// this is an awesome way since it prevents a global decalaration
// and we interact with the db through this struct
type Db struct {
	*sql.DB
}

// New makes a new database using the connection string and returns it,
// otherwise, returns an error
func New(connString string) (*Db, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Check that our connection is good
	// we do this through a Ping to the db
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Db{db}, nil
}

// ConnString returns a connection string based on the parameters it is given
// this would normally contain a password
// Please use a password in production @todo
func ConnString(host string, port int, user string, dbName string) string {
	return fmt.Sprintf(
		// enable ssl for secure connections
		"host=%s port=%d user=%s dbName=%s sslmode=disable",
		host, port, user, dbName,
	)
}

// User structure
type User struct {
	ID         int
	Name       string
	Age        int
	Profession string
	Friendly   bool
}

// GetUsersByName is called within our user query for graphql
func (d *Db) GetUsersByName(name string) []User {
	// Prepare query, takes a name argument, protects from sql injection
	stmt, err := d.Prepare("SELECT * FROM users WHERE name=$1")
	if err != nil {
		fmt.Println("GetUserByName Preparation error: ", err)
	}

	// Make query with our stmt, passing it name argument
	rows, err := stmt.Query(name)
	if err != nil {
		fmt.Println("GetUsrByName Query error: ", err)
	}

	// Create User struct for holding each rows data
	var r User
	// Create slice of Users for our response
	users := []User{}
	// copy the columns from row into the values pointed at by r(User)
	for rows.Next() {
		err = rows.Scan(
			&r.ID,
			&r.Name,
			&r.Age,
			&r.Profession,
			&r.Friendly,
		)
		if err != nil {
			fmt.Println("Error scanning rows: ", err)
		}
		users = append(users, r)
	}
	return users
}

package psql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDBClient struct {
	host         string
	port         int
	dbname       string
	user         string
	password     string
	DbConnection *sql.DB
}

// NewPostgresDBClient acts like a constructor to generating a new object of PostgresDBClient with a DbConnection as well
func NewPostgresDBClient(host string, port int, dbname, user, password string) (*PostgresDBClient, error) {
	pc := &PostgresDBClient{
		host:     host,
		port:     port,
		dbname:   dbname,
		user:     user,
		password: password,
	}
	if _, err := pc.setupAndReturnDbConnection(); err != nil {
		return nil, err
	}
	return pc, nil
}

// Receives a PostgresDBClient object and setups a new DbConnection connection for it (and returns it) if need be, else returns the existing one functional one
// Not setting up a new connection blindly everything this method is called so as to maintain idempotency
func (pc *PostgresDBClient) setupAndReturnDbConnection() (*sql.DB, error) {
	// if the .DbConnection attribute is nil, then we clearly need to setup a new connection as none exists presently
	// else if it exists but has some other issues (basically, some error occurred while "Ping"ing the DB with current connection), then setup a new connection
	// PS: expiry of connection won't be a concern because .Ping() not only checks the connection but resets the connections too automatically with new connections, if need be.
	// Ref: https://cs.opensource.google/go/go/+/refs/tags/go1.17.1:src/database/sql/sql.go;l=868
	setupNewConnection := false
	if pc.DbConnection == nil {
		setupNewConnection = true
	} else if err := pc.DbConnection.Ping(); err != nil {
		setupNewConnection = true
	}

	if setupNewConnection {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			pc.host, pc.port, pc.user, pc.password, pc.dbname)
		newDbConnection, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			return nil, err
		}
		pc.DbConnection = newDbConnection
	}
	return pc.DbConnection, nil
}

// Insert inserts a row into the DB to which the receiver PostgresDBClient poins
func (pc *PostgresDBClient) Insert(id, table, name string, age int, country string) error {
	dbConnection, err := pc.setupAndReturnDbConnection()
	if err != nil {
		return err
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (id, name, age, country) VALUES ('%s', '%s', %d, '%s') ON CONFLICT (id) DO NOTHING;", table, id, name, age, country)
	if _, err := dbConnection.Exec(insertQuery); err != nil {
		return err
	}
	return nil
}

// Delete deletes row from the DB to which the receiver PostgresDBClient points
func (pc *PostgresDBClient) Delete(id, table string) error {
	dbConnection, err := pc.setupAndReturnDbConnection()
	if err != nil {
		return err
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id='%s';", table, id)
	if _, err := dbConnection.Exec(deleteQuery); err != nil {
		return err
	}
	return nil
}

// Close is used to gracefully close and wrap up the DbConnection associated with the PostgresDBClient object so as to avoid any associated memory leaks
func (pc *PostgresDBClient) Close() {
	pc.DbConnection.Close()
}

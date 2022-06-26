package tools

import (
	"database/sql"
	"fmt"
    "os"
    "strconv"

	_ "github.com/lib/pq"
)

type Confer interface {
    Conf() (err error)
}
type SqlOpener interface {
    SqlOpen() (db *sql.DB, err error)
}
type SqlConferOpener interface {
    Confer
    SqlOpener
}
func Open(loc SqlConferOpener) (db *sql.DB, err error){
    err = loc.Conf()
    db, err = loc.SqlOpen()
    return db, err
}

type PostgresDb struct {
    Host string
    Port int
    User string
    Pass string
    Dnam string
    Sslm string

}
func (s *PostgresDb) Conf() (err error) {
    // env variables must be available
    s.Host = os.Getenv("PGHOST")
    port, err := strconv.Atoi(os.Getenv("PGPORT"))
    s.Port = port
    s.User = os.Getenv("PGUSER")
    s.Pass = os.Getenv("PGPASS")
    s.Dnam = os.Getenv("PGDNAM")
    s.Sslm = os.Getenv("PGSSLM")
    if err != nil {
        return fmt.Errorf("Error configuring postgres db. %s", err)
    }
    return nil
}
func (s *PostgresDb) SqlOpen() (db *sql.DB, err error) {
	// Connect to the postgres development database
    conn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.Host,
		s.Port,
		s.User,
		s.Pass,
		s.Dnam,
		s.Sslm,
	)
	db, err = sql.Open("postgres", conn)
	if err != nil {
		return db, fmt.Errorf("Unable to connect to the postgres db. %s", err)
	}
	return db, nil
}

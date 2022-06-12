package tools

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/robertsmoto/skustor/configs"
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

type PostgresDev struct {
    Host string
    Port int
    User string
    Pass string
    Dnam string
    Sslm string

}
func (s *PostgresDev) Conf() (err error) {
	config := configs.Config{}
    err = configs.Load(&config)
    s.Host = config.DbDevelopment.Host
    s.Port = config.DbDevelopment.Port
    s.User = config.DbDevelopment.User
    s.Pass = config.DbDevelopment.Pass
    s.Dnam = config.DbDevelopment.Dnam
    s.Sslm = config.DbDevelopment.Sslm
    if err != nil {
        log.Print("Error configuring postgres development db.", err)
        fmt.Println("##s -->", s)
    }
    return err
}
func (s *PostgresDev) SqlOpen() (db *sql.DB, err error) {
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
		log.Print("Unable to connect to the postgres db", err)
	}
	return db, err
}

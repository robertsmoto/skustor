package postgres

import (
    "example.com/internal/conf"
    "database/sql"
    "fmt"

    _ "github.com/lib/pq"
)


var Conf = conf.Conf

func Connect() {

    fmt.Println(Conf)
    psqlInfo := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        Conf.DbSwimExpress.Host, 
        Conf.DbSwimExpress.Port,
        Conf.DbSwimExpress.User,
        Conf.DbSwimExpress.Pass,
        Conf.DbSwimExpress.Dnam,
        Conf.DbSwimExpress.Sslm,
    )

    db, err := sql.Open("postgres", psqlInfo)

    sqlStatement := `
        SELECT id, source_headline_url
        FROM headlines_headlinepost`
    rows, err := db.Query(sqlStatement)

    if err != nil {
        panic(err)
    }

    for rows.Next(){

		var id int
        var source_headline_url string

		if err := rows.Scan(&id, &source_headline_url); err != nil {
			panic(err)
		}

        fmt.Printf("%d: %s", id, source_headline_url)
        fmt.Println()
    }
}

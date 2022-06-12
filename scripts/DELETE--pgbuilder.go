package scripts

import (
    "fmt"
    "github.com/robertsmoto/skustor/tools"
)

func BuildDb() (err error) {
    qstr := ``
    fmt.Println(qstr)

    // open each db
    devPostgres := tools.PostgresDev{}
    devConn, err := tools.Open(&devPostgres)
    fmt.Println(devConn)
    devConn.Close()
    
    // run the query on each db

    // close the db
    return err
}

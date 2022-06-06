package main

import (
    "example.com/internal/conf"
    "example.com/repos/postgres"
    "fmt"
    "strings"
)

var Conf = conf.Conf

func main() {
    fmt.Println("Starting the application ...")
    //// put config func here
    fmt.Println(Conf.TempFileDir)
    s := "Hello World"
    fmt.Println("Hello World")
    fmt.Println()
    fmt.Println(s[:4])
    fmt.Println(s[5:])
    fmt.Println(strings.Replace(s, "l", "-", -1))
    var k string
    k = strings.Replace(s, "l", "_", -1)
    fmt.Println(k)

    for i, c := range(s) {
        fmt.Printf("%v: %c", i, c)
        fmt.Println()
    }

    // connect to a database, query it and then loop thru the results
    postgres.Connect()

}

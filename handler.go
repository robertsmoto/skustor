package main

import (
    "fmt"
    "log"
    //"os"
    //"encoding/json"
    //"io/ioutil"
    "net/http"

    //"github.com/gorilla/mux"
    "github.com/robertsmoto/skustor/internal/configs"
    //"github.com/robertsmoto/skustor/internal/models"
	//"github.com/robertsmoto/skustor/internal/postgres"

)


func upsertDataHandler(w http.ResponseWriter, r *http.Request) {

    fmt.Println("Endpoint hit: upsertData")
    fmt.Println("request: r")

    //vars := mux.Vars(r)

    //name := vars["name"]

	var err error

	// loading env variables (will eventually be loaded by main)
	conf := configs.Config{}
	configs.Load(&conf)
	if err != nil {
        log.Print("error")
	}

	// read file (will eventually come from the request)
	//if err != nil {
        //log.Print("error")
	//}

	//// open the db connections
	//pgDb, err := postgres.Open(&postgres.PostgresDb{})

	//// instantiate the structs
	////collection := Collection{}
	//collections := models.CollectionNodes{}

	//// Little Johnnie user
	//userId := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"

	//structs := []models.LoaderProcesserUpserter{&collections}
	//for _, s := range structs {
		//err = models.JsonLoaderUpserterHandler(s, userId, &testFile, pgDb)
		//if err != nil {
            //log.Print("error")
		//}
	//}

}


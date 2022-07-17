/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
    "errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
    "time"

	"github.com/robertsmoto/skustor/cmd"
	"github.com/robertsmoto/skustor/internal/configs"
	"github.com/robertsmoto/skustor/internal/models"
	"github.com/robertsmoto/skustor/internal/postgres"
    "github.com/graphql-go/graphql"
    "github.com/graphql-go/graphql-go-handler"
)

func main() {
	// for the cli
	cmd.Execute()

	configs.Load(&configs.Config{})

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/upsert", UpsertDataHandler).Methods("POST")
	r.HandleFunc("/delete", DeleteDataHandler).Methods("DELETE")
    r.HandleFunc("/gql", GqlHandler).Methods("GET") 
	fmt.Println("Listening on port 10000 ...")
	// change to only listen on https in production
	log.Fatal(http.ListenAndServe(":10000", r))

}

func AuthenticateRequest(w http.ResponseWriter, r *http.Request) (aid string, err error) {

	// the request validation can be written in a separate function
	// limit number requests per min
	reqOrig := r.RemoteAddr
	fmt.Println("Request originated from: ", reqOrig)

	// limit size of request
	r.Body = http.MaxBytesReader(w, r.Body, 10000000) // 10 Mb

	header := r.Header
	// authorize
	auth := header.Get("auth")
	ckAuth := "88355a6d-bffa-448f-8c73-1b420031dc95"
	prefix := header.Get("prefix")
	ckPrefix := "c083bf1d-8878-4ae3-97bd-e521f70c4717"
	aid = header.Get("aid")
	ckAid := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"
	sub := "CURRENT"
	ckSub := "CURRENT"

	// this will come from sv-user table / redis eventually
	if auth != ckAuth || prefix != ckPrefix || aid != ckAid || sub != ckSub {
		return aid, errors.New("User not authenticated.")
	}

	contentType := header.Get("content-type")
	if contentType != "application/json" {
        return aid, errors.New(
            "Content-Type application/json is required in the request header.")
	}
    return aid, nil
}

func UpsertDataHandler(w http.ResponseWriter, r *http.Request) {

    start := time.Now()

    aid, err := AuthenticateRequest(w, r)
    if err != nil {
        log.Printf("AuthenticationRequest %s %s", err, aid)
        w.WriteHeader(500)
        w.Write([]byte("Internal error."))
        return
    }

	// read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
        log.Printf("Reading request body. %s", err)
        w.WriteHeader(500)
        w.Write([]byte("Internal error."))
        return
	}

	// open the db connections
	devPostgres := postgres.PostgresDb{}
	pgDb, err := postgres.Open(&devPostgres)

	// instantiate the structs
	accountNodes := models.AccountNodes{}
	collectionNodes := models.CollectionNodes{}
    contentNodes := models.ContentNodes{}
    itemNodes := models.ItemNodes{}
    imageNodes := models.ImageNodes{}
    personNodes := models.PersonNodes{}
    placeNodes := models.PlaceNodes{}

	// loader validator
    lvNodes := []models.LoaderValidator{
        &accountNodes,
        &collectionNodes,
        &contentNodes,
        &itemNodes,
        &imageNodes,
        &personNodes,
        &placeNodes,
    }

	for _, node := range lvNodes {
        err = models.LoadValidateHandler(node, &body)
        if err != nil {
            log.Println("main 01", err)
            w.WriteHeader(500)
            w.Write([]byte("Internal error."))
        }
	}

    // preProcessing
    ppNodes := []models.PreProcessor{
        &imageNodes,
    }
	for _, node := range ppNodes {
        err = models.PreProcessHandler(node, aid, pgDb)
        if err != nil {
            log.Println("main 02", err)
            w.WriteHeader(500)
            w.Write([]byte("Internal error."))
        }
	}

	// upsert
    upsertNodes := []models.Upserter{
        &accountNodes,
        &collectionNodes,
        &contentNodes,
        &itemNodes,
        &imageNodes,
        &personNodes,
        &placeNodes,
    }
	for _, node := range upsertNodes {
        err = models.UpsertHandler(node, aid, pgDb)
        if err != nil {
            log.Println("main 03", err)
            w.WriteHeader(500)
            w.Write([]byte("Internal error."))
        }
	}

	// foreign key update
    fkNodes := []models.ForeignKeyUpdater{
        &accountNodes,
        &collectionNodes,
        &contentNodes,
        &itemNodes,
        &imageNodes,
        &personNodes,
        &placeNodes,
    }
	for _, node := range fkNodes {
        err = models.ForeignKeyUpdateHandler(node, pgDb)
        if err != nil {
            log.Println("main 04", err)
            w.WriteHeader(500)
            w.Write([]byte("Internal error."))
        }
	}

	// related table upsert
    rtNodes := []models.RelatedTableUpserter{
        &accountNodes,
        &collectionNodes,
        &contentNodes,
        &itemNodes,
        &imageNodes,
        &personNodes,
        &placeNodes,
    }
	for _, node := range rtNodes {
        err = models.RelatedTableUpsertHandler(node, aid, pgDb)
        if err != nil {
            log.Println("main 05", err)
            w.WriteHeader(500)
            w.Write([]byte("Internal error."))
        }
	}

	pgDb.Close()
    elapsed := time.Since(start)
	w.WriteHeader(200)
    timeTook := fmt.Sprintf("Gets it done fast. Upsert time %s", elapsed)
    w.Write([]byte(timeTook))
}

func DeleteDataHandler(w http.ResponseWriter, r *http.Request) {

    start := time.Now()

    aid, err := AuthenticateRequest(w, r)
    if err != nil {
        log.Printf("AuthenticationRequest %s %s", err, aid)
        w.WriteHeader(500)
        fmtErr := fmt.Sprintf("%s", err)
        w.Write([]byte(fmtErr))
        return
    }

	// read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
        log.Printf("Reading request body. %s", err)
        w.WriteHeader(500)
        w.Write([]byte("Internal error."))
        return
	}

	// open the db connections
	devPostgres := postgres.PostgresDb{}
	pgDb, err := postgres.Open(&devPostgres)

	// instantiate the structs
	accountNodes := models.AccountNodes{}
	collectionNodes := models.CollectionNodes{}
    contentNodes := models.ContentNodes{}
    itemNodes := models.ItemNodes{}
    imageNodes := models.ImageNodes{}
    personNodes := models.PersonNodes{}
    placeNodes := models.PlaceNodes{}

	// loader validator
    lvNodes := []models.LoaderValidator{
        &accountNodes,
        &collectionNodes,
        &contentNodes,
        &itemNodes,
        &imageNodes,
        &personNodes,
        &placeNodes,
    }

	for _, node := range lvNodes {
        err = models.LoadValidateHandler(node, &body)
        if err != nil {
            log.Println("main 05", err)
            w.WriteHeader(500)
            w.Write([]byte("Internal error."))
        }
	}

    dNodes := []models.Deleter{
        &accountNodes,
        &collectionNodes,
        &contentNodes,
        &itemNodes,
        //&imageNodes,
        &personNodes,
        &placeNodes,
    }

	for _, node := range dNodes {
        err = models.DeleteHandler(node, pgDb)
        if err != nil {
            log.Println("main 06", err)
            w.WriteHeader(500)
            w.Write([]byte("Internal error."))
        }
	}

	pgDb.Close()
    elapsed := time.Since(start)
	w.WriteHeader(200)
    timeTook := fmt.Sprintf("Gets it done fast. Delete time %s", elapsed)
    w.Write([]byte(timeTook))
}


func GqlHandler(w http.ResponseWriter, r *http.Request) {

    var queryType = graphql.NewObject(graphql.ObjectConfig{
        Name: "Query",
        Fields: graphql.Fields{
            "latestPost": &graphql.Field{
                Type: graphql.String,
                Resolve: func(p types.ResolveParams) interface{} {
                    return "Hello World!"
                },
            },
        },
    })

    var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
        Query: queryType,
    })

    start := time.Now()
    // code here
    w.Write([]byte("Graphql not implemented."))
    elapsed := time.Since(start)


	w.WriteHeader(200)
    timeTook := fmt.Sprintf("Gets it done fast. Delete time %s", elapsed)
    w.Write([]byte(timeTook))
}

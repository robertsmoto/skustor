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
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	//"net/http/httputil"

	"github.com/robertsmoto/skustor/cmd"
	"github.com/robertsmoto/skustor/internal/configs"
	"github.com/robertsmoto/skustor/internal/models"
	"github.com/robertsmoto/skustor/internal/postgres"
)

func main() {
	// for the cli
	cmd.Execute()

	configs.Load(&configs.Config{})
	fmt.Println("## Loaded configs env variables ...")

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/upsert", UpsertDataHandler).Methods("POST")
	r.HandleFunc("/delete", DeleteDataHandler).Methods("DELETE")
	fmt.Println("Listening on port 10000 ...")
	// change to only listen on https in production
	log.Fatal(http.ListenAndServe(":10000", r))

}

func UpsertDataHandler(w http.ResponseWriter, r *http.Request) {

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
	uid := header.Get("uid")
	ckUid := "f8b0f997-1dcc-4e56-915c-9f62f52345ee"
	sub := "CURRENT"
	ckSub := "CURRENT"

	// this will come from sv-user table / redis eventually
	if auth != ckAuth || prefix != ckPrefix || uid != ckUid || sub != ckSub {
        log.Fatal("User not authorized to access the api.")
        w.WriteHeader(500)
        w.Write([]byte("User not authorized to access the api."))
		return
	}

	fmt.Println("## header", header)
	contentType := header.Get("content-type")
	if contentType != "application/json" {
        log.Fatal("Content-Type application/json is required in the request header.")
        w.WriteHeader(500)
        w.Write([]byte("Content-Type application/json is required."))
		return
	}

	// read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
        log.Fatal("Reading request body.")
        w.WriteHeader(500)
        w.Write([]byte("User not authorized to access the api."))
	}

	// instantiate the structs
	collectionNodes := models.CollectionNodes{}
    priceClassNodes := models.PriceClassNodes{}
    unitNodes := models.UnitNodes{}
    itemNodes := models.ItemNodes{}
	contentNodes := models.ContentNodes{}

	// add them to slice of interface
	loaderNodes := []models.LoaderProcesserUpserter{
		&collectionNodes,
        &priceClassNodes,
        &unitNodes,
        &itemNodes,
		&contentNodes,
	}

	// open the db connections
	devPostgres := postgres.PostgresDb{}
	pgDb, err := postgres.Open(&devPostgres)

	for _, node := range loaderNodes {
        
        fmt.Printf("## node ", node, " type %T ", node, "\n")
		err = models.JsonLoaderUpserterHandler(
			node, uid, &body, pgDb,
		)
		if err != nil {
            // return error
            log.Printf("main.Main() json loader.", err)
            w.WriteHeader(500)
            w.Write([]byte("Internal error."))
		}
	}
	pgDb.Close()
	w.WriteHeader(200)
}

func DeleteDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Fatal("Not implemented.")
}

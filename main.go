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
    "log"
    "net/http"
    "github.com/gorilla/mux"

    "github.com/robertsmoto/skustor/internal/configs"
    "github.com/robertsmoto/skustor/cmd"
)

func main() {
    // for the cli
	cmd.Execute()

    configs.Load(&configs.Config{})
    fmt.Print("Loaded configs env variables ...")

    r := mux.NewRouter().StrictSlash(true)
    //r.HandleFunc("/allgroceries", AllGroceries) // ----> To request all groceries
    //r.HandleFunc("/groceries/{name}", SingleGrocery) // ----> To request a specific grocery
    r.HandleFunc("/upsert", upsertData).Methods("POST") // ----> To add  new grocery to buy
    //r.HandleFunc("/groceries/{name}", UpdateGrocery).Methods("PUT")// ----> To update a grocery
    //r.HandleFunc("/groceries/{name}", DeleteGrocery).Methods("DELETE") // ----> Delete a grocery
    log.Fatal(http.ListenAndServe(":10000", r))
}





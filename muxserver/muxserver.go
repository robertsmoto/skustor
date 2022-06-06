package muxserver

import (
    "fmt"
    //"log"
    "net/http"

    //"github.com/gorilla/mux"
)

func HomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}


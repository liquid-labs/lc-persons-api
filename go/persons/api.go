package persons

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/Liquid-Labs/catalyst-core-api/go/restserv"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "/persons alive\n")
}

func createHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "TODO")
}

func listHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  contextType := vars["contextType"]

  if contextType == "" {
    fmt.Fprintf(w, "TODO: %+v\n", r.Context().Value(restserv.FireauthKey))
  } else {
    contextId := vars["contextId"]
    fmt.Fprintf(w, "TODO: in context %s/%s\n", contextType, contextId)
  }
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  pubId := vars["pubId"]
  fmt.Fprintf(w, "TODO: GET %s\n", pubId)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  pubId := vars["pubId"]
  fmt.Fprintf(w, "TODO: POST %s\n", pubId)
}

func InitAPI(r *mux.Router) {
  r.HandleFunc("/persons/", pingHandler).Methods("PING")
  r.HandleFunc("/persons/", createHandler).Methods("POST")
  r.HandleFunc("/persons/", listHandler).Methods("GET")
  r.HandleFunc("/{contextType:[a-z-]*[a-z]}/{contextId:[0-9]+}/persons/", listHandler).Methods("GET")
  r.HandleFunc("/persons/{pubId}", detailHandler).Methods("GET")
  r.HandleFunc("/persons/{pubId}", updateHandler).Methods("PUT")
}

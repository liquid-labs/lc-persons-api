package persons

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"
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
    fmt.Fprint(w, "TODO: listing persons")
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

const uuidRe = `[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}`

func InitAPI(r *mux.Router) {
  r.HandleFunc("/persons/", pingHandler).Methods("PING")
  r.HandleFunc("/persons/", createHandler).Methods("POST")
  r.HandleFunc("/persons/", listHandler).Methods("GET")
  r.HandleFunc("/{contextType:[a-z-]*[a-z]}/{contextId:" + uuidRe + "}/persons/", listHandler).Methods("GET")
  r.HandleFunc("/persons/{pubId:" + uuidRe + "}/", detailHandler).Methods("GET")
  r.HandleFunc("/persons/{pubId:" + uuidRe + "}/", updateHandler).Methods("PUT")
}

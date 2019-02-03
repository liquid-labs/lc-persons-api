package persons

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "/persons alive\n")
}

func InitAPI(r *mux.Router) {
  // r.HandleFunc("/persons/", pingHandler).Methods("GET")

  rUsers := r.PathPrefix("/persons/").Subrouter()

  rUsers.HandleFunc("/", pingHandler).Methods("GET")
}

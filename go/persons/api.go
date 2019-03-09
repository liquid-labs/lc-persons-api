package persons

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"

  "github.com/Liquid-Labs/catalyst-core-api/go/restserv"
  "github.com/Liquid-Labs/catalyst-firewrap/go/fireauth"
  "github.com/Liquid-Labs/go-rest/rest"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "/persons alive\n")
}

func createHandler(w http.ResponseWriter, r *http.Request) {
  authClient := r.Context().Value(restserv.FireauthKey).(*fireauth.ScopedClient)
  authToken, restErr := authClient.GetToken() // effectively checks if user authorized
  if restErr != nil {
    rest.HandleError(w, restErr)
    return
  }

  var person Person
  if err := rest.ExtractJson(w, r, &person, `Person`); err != nil {
    // HTTP response is already set by 'ExtractJson'
    return
  }
  if authToken.UID != person.AuthId.String {
    rest.HandleError(w, rest.AuthorizationError("Create record for yourself.", nil))
    return
  }

  newPerson, err := CreatePerson(&person, r.Context())
  if err != nil {
    rest.HandleError(w, err)
    return
  }

  rest.StandardResponse(w, newPerson, `Person created.`, nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  contextType := vars["contextType"]

  if contextType == "" {
    fmt.Fprint(w, "TODO: listing persons")
  } else {
    contextType := vars["contextType"]
    contextId := vars["contextId"]
    // TODO: distribute 'join' defs as common includes to all resources in a
    // given system; e.g., list is compiled at app.
    // TODO: make internal REST call to get the 'JOIN' info for unknowns?
    fmt.Fprintf(w, "TODO: in context %s/%s\n", contextType, contextId)
  }
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
  authClient := r.Context().Value(restserv.FireauthKey).(*fireauth.ScopedClient)
  authToken, restErr := authClient.GetToken() // effectively checks if user authorized
  if restErr != nil {
    rest.HandleError(w, restErr)
    return
  }

  vars := mux.Vars(r)
  authId := vars["authId"]
  pubId := vars["pubId"]
  if authId == `` {
    if pubId != `self` {
      rest.HandleError(w, rest.AuthorizationError("May only request your own data. Try '/persons/self'. (1)", nil))
      return
    }
    authId = authToken.UID
  } else if authId != authToken.UID {
    rest.HandleError(w, rest.AuthorizationError("May only request your own data. Try '/persons/self'. (2)", nil))
    return
  }

  var person *Person
  var err rest.RestError
  if authId != `` {
    person, err = GetPersonByAuthId(authToken.UID, r.Context())
  } else {
    // not currently used, but will do once general authorization system in place
    person, err = GetPerson(pubId, r.Context())
  }
  if err != nil {
    rest.HandleError(w, err)
    return
  } else {
    rest.StandardResponse(w, person, `Person retrieved.`, nil)
  }
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  pubId := vars["pubId"]
  fmt.Fprintf(w, "TODO: POST %s\n", pubId)
}

const uuidRe = `[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}`
const personIdRe = uuidRe + `|self`

func InitAPI(r *mux.Router) {
  r.HandleFunc("/persons/", pingHandler).Methods("PING")
  r.HandleFunc("/persons/", createHandler).Methods("POST")
  r.HandleFunc("/persons/", listHandler).Methods("GET")
  r.HandleFunc("/{contextType:[a-z-]*[a-z]}/{contextId:" + uuidRe + "}/persons/", listHandler).Methods("GET")
  r.HandleFunc("/persons/{pubId:" + personIdRe + "}/", detailHandler).Methods("GET")
  r.HandleFunc("/persons/{pubId:" + personIdRe + "}/", updateHandler).Methods("PUT")
  // special auth-id fetcher
  r.HandleFunc("/persons/auth-id-{authId}/", detailHandler).Methods("GET")
}

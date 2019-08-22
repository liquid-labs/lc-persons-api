package persons

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"

  "github.com/Liquid-Labs/lc-authentication-api/go/auth"
  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  "github.com/Liquid-Labs/go-rest/rest"
  model "github.com/Liquid-Labs/lc-persons-model/go/persons"
  . "github.com/Liquid-Labs/terror/go/terror"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "/persons is alive\n")
}

func requireAuthentication(w http.ResponseWriter, r *http.Request) (bool, string) {
  authOracle := auth.GetAuthOracleFromContext(r.Context())
  if authOracle == nil || !authOracle.IsRequestAuthenticated() {
    rest.HandleError(w, UnauthenticatedError("Request must be authenticated."))
    return false, ``
  } else { return true, authOracle.GetAuthID() }
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
  var person *model.Person = &model.Person{}
  ok, authID := requireAuthentication(w, r)
  if ok {
    if err := rest.ExtractJson(w, r, &person, `Person`); err != nil {
      rest.HandleError(w, err); return
    } else if authID != person.GetAuthID() {
      rest.HandleError(w, ForbiddenError("You can create a record only for yourself. (Auth IDs must match.)"))
      return
    }

    if cErr := person.CreateSelf(r.Context()); cErr != nil {
      rest.HandleError(w, ServerError("Could not create person.", cErr))
    } else {
      rest.StandardResponse(w, person, `Person created.`, nil)
    }
  }
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
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

func DetailHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  reqAuthID := vars["authID"]
  pubID := vars["pubID"]

  ok, authID := requireAuthentication(w, r)
  if ok {
    if authID == reqAuthID || pubID == `self` {
      p, err := model.RetrievePersonSelf(rdb.ConnectWithContext(r.Context()))
      if (err != nil) {
        rest.HandleError(w, ServerError("Error retrieving person.", err))
      } else {
        rest.StandardResponse(w, p, `Person retrieved.`, nil)
      }
    } else {
      // We do not currenty support non-self Person details
      rest.HandleError(w, ForbiddenError("May only request your own data. Try '/persons/self'. (1)"))
      return
    }
  }
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
  ok, _ := requireAuthentication(w, r)
  if ok {
    newData := &model.Person{}
    if err := rest.ExtractJson(w, r, &newData, `Person`); err != nil {
      // HTTP response is already set by 'ExtractJson'
      return
    }

    vars := mux.Vars(r)
    pubID := vars["pubID"]

    if string(newData.ID) != pubID {
      rest.HandleError(w, ForbiddenError("You can only update your own data."))
      return
    }

    err := newData.UpdateSelf(rdb.Connect())
    if err != nil {
      rest.HandleError(w, ServerError("Error updating person.", err))
    } else {
      rest.StandardResponse(w, newData, `Person updated.`, nil)
    }
  }
}

const uuidRe = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}`
const personIdRe = `(?:` + uuidRe + `|self)`

func InitAPI(r *mux.Router) {
  r.HandleFunc("/", pingHandler).Methods("PING")
  r.HandleFunc("/persons/", CreateHandler).Methods("POST")
  r.HandleFunc("/persons/", ListHandler).Methods("GET")
  r.HandleFunc("/{contextType:[a-z-]*[a-z]}/{contextId:" + uuidRe + "}/persons/", ListHandler).Methods("GET")
  r.HandleFunc("/persons/{pubID:" + personIdRe + "}/", DetailHandler).Methods("GET")
  r.HandleFunc("/persons/{pubID:" + personIdRe + "}/", UpdateHandler).Methods("PUT")
  // special auth-id fetcher
  r.HandleFunc("/{foo:persons}/auth-id-{authID:.+}/", DetailHandler).Methods("GET")
}

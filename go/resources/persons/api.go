package persons

import (
  "fmt"
  "log"
  "net/http"

  "github.com/gorilla/mux"

  "github.com/Liquid-Labs/catalyst-core-api/go/handlers"
  "github.com/Liquid-Labs/go-rest/rest"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "/persons is alive\n")
}

func createHandler(w http.ResponseWriter, r *http.Request) {
  log.Print("Creating person.")
  var person *Person = &Person{}
  if authToken, restErr := handlers.CheckAndExtract(w, r, person, `Person`); restErr != nil {
    return // response handled by CheckAndExtract
  } else {
    if authToken.UID != person.AuthId.String {
      rest.HandleError(w, rest.AuthorizationError("You can create a record only for yourself.", nil))
      return
    }

    handlers.DoCreate(w, r, CreatePerson, person, `Person`)
  }
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
  if authToken, restErr := handlers.BasicAuthCheck(w, r); restErr != nil {
    return // response handled by BasicAuthCheck
  } else {
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

    if authId != `` {
      handlers.DoGetDetail(w, r, GetPersonByAuthId, authToken.UID, `Person`)
    } else {
      // not currently used, but will do once general authorization system in place
      handlers.DoGetDetail(w, r, GetPerson, pubId, `Person`)
    }
  }
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
  var newData *Person = &Person{}
  if authToken, restErr := handlers.CheckAndExtract(w, r, newData, `Person`); restErr != nil {
    return // response handled by CheckAndExtract
  } else {
    vars := mux.Vars(r)
    pubID := vars["pubId"]

    // TODO: This is essentially to do the auth check. Once we have a proper
    // authorization infrasatructure, this can go away.
    var user *Person
    var err rest.RestError
    user, err = GetPersonByAuthId(authToken.UID, r.Context())
    if err != nil {
      rest.HandleError(w, err)
      return
    }

    if newData.PubId.String != user.PubId.String {
      rest.HandleError(w, rest.AuthorizationError("You can only update your own data.", nil))
      return
    }

    handlers.DoUpdate(w, r, UpdatePerson, newData, pubID, `Person`)
  }
}

const uuidRe = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}`
const personIdRe = `(?:` + uuidRe + `|self)`

func InitAPI(r *mux.Router) {
  r.HandleFunc("/", pingHandler).Methods("PING")
  r.HandleFunc("/persons/", createHandler).Methods("POST")
  r.HandleFunc("/persons/", listHandler).Methods("GET")
  r.HandleFunc("/{contextType:[a-z-]*[a-z]}/{contextId:" + uuidRe + "}/persons/", listHandler).Methods("GET")
  r.HandleFunc("/persons/{pubId:" + personIdRe + "}/", detailHandler).Methods("GET")
  r.HandleFunc("/persons/{pubId:" + personIdRe + "}/", updateHandler).Methods("PUT")
  // special auth-id fetcher
  r.HandleFunc("/persons/auth-id-{authId}/", detailHandler).Methods("GET")
}

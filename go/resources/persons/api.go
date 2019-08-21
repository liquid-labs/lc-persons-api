package persons

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"

  "github.com/Liquid-Labs/lc-authentication-api/go/auth"
  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  "github.com/Liquid-Labs/go-rest/rest"
  model "github.com/Liquid-Labs/lc-persons-model/go/persons"
  . "github.com/Liquid-Labs/terror/go/terror"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "/persons is alive\n")
}

func checkAuthentication(w http.ResponseWriter, r *http.Request) (bool, string) {
  authenticator, authID, err := auth.CheckAuthentication(r.Context())
  if err != nil {
    rest.HandleError(w, ServerError("Error checking authentication.", err))
    return false, ``
  } else if !authenticator.IsRequestAuthenticated() {
    rest.HandleError(w, UnauthenticatedError("Request must be authenticated."))
    return false, ``
  } else { return true, authID }
}

func createHandler(w http.ResponseWriter, r *http.Request) {
  var person *model.Person = &model.Person{}
  ok, authID := checkAuthentication(w, r)
  if ok {
    if authID != person.GetAuthID() {
      rest.HandleError(w, ForbiddenError("You can create a record only for yourself."))
      return
    }

    im := ConnectItemManager()
    cErr := im.CreateRaw(person)
    if (cErr != nil) {
      rest.HandleError(w, ServerError("Could not create person.", cErr))
    } else {
      rest.StandardResponse(w, person, `Person created.`, nil)
    }
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
  ok, authID := checkAuthentication(w, r)
  if ok {
    vars := mux.Vars(r)
    reqAuthID := vars["authId"]
    pubId := vars["pubId"]
    if reqAuthID == `` {
      if pubId != `self` {
        rest.HandleError(w, ForbiddenError("May only request your own data. Try '/persons/self'. (1)"))
        return
      }
    } else if reqAuthID != authID {
      rest.HandleError(w, ForbiddenError("May only request your own data. Try '/persons/self'. (2)"))
      return
    }

    if authID != `` {
      p, err := model.RetrievePersonSelf(rdb.ConnectWithContext(r.Context()))
      if (err != nil) {
        rest.HandleError(w, ServerError("Error retrieving person.", err))
      } else {
        rest.StandardResponse(w, p, `Person retrieved.`, nil)
      }
    } else {
      // not currently used, but will do once general authorization system in place; at the moment we exit after checkAuthenticaiton
      // handlers.DoGetDetail(w, r, GetPerson, pubId, `Person`)
      return
    }
  }
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
  ok, _ := checkAuthentication(w, r)
  if ok {
    newData := &model.Person{}
    if err := rest.ExtractJson(w, r, &newData, `Person`); err != nil {
      // HTTP response is already set by 'ExtractJson'
      return
    }

    vars := mux.Vars(r)
    pubID := vars["pubId"]

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
  r.HandleFunc("/persons/", createHandler).Methods("POST")
  r.HandleFunc("/persons/", listHandler).Methods("GET")
  r.HandleFunc("/{contextType:[a-z-]*[a-z]}/{contextId:" + uuidRe + "}/persons/", listHandler).Methods("GET")
  r.HandleFunc("/persons/{pubId:" + personIdRe + "}/", detailHandler).Methods("GET")
  r.HandleFunc("/persons/{pubId:" + personIdRe + "}/", updateHandler).Methods("PUT")
  // special auth-id fetcher
  r.HandleFunc("/persons/auth-id-{authId}/", detailHandler).Methods("GET")
}

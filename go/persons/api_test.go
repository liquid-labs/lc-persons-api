package persons_test

import (

  "bytes"
  "context"
  "fmt"
  "io"
  "net/http"
  "net/http/httptest"
  "os"
  "testing"

  . "github.com/golang/mock/gomock"
  // "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"
  "github.com/gorilla/mux"

  "github.com/Liquid-Labs/lc-authentication-api/go/auth"
  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  . "github.com/Liquid-Labs/lc-locations-model/go/locations"
  authmock "github.com/Liquid-Labs/lc-authentication-api/go/mock"
  api "github.com/Liquid-Labs/lc-persons-api/go/persons"
  . "github.com/Liquid-Labs/lc-persons-model/go/persons"
  "github.com/Liquid-Labs/strkit/go/strkit"
  "github.com/Liquid-Labs/terror/go/terror"
)

func init() {
  terror.EchoErrorLog()
}

func joeBobJSON(authID string) []byte {
  return []byte(`{
    "authId": "` + authID + `",
    "name": "Joe Bob",
    "givenName": "Joe",
    "familyName": "Bob",
    "email": "jbob@foo.com",
    "phone": "555-565-383",
    "backupPhone": "555-384-2832",
    "avatarUrl": "https://avatars.com/joeBob",
    "addresses": [
      {
      "address1": "100 Main Str",
      "city": "Anwhere",
      "state": "TX",
      "zip": "78383-4833",
      "label": "home"
    }]
  }`)
}

func joeBobStruct(authID string) *Person {
  a1 := NewAddress(`Camelot`, `a house`, EID(``), false, `100 Main Str`, `#B`, `Paris`, `TX`, `78383`, EID(``), `home`)
  as := Addresses{a1}
  return NewPerson(`Joe Bob`,
    `A man`,
    authID,
   `444-53-3838`,
   `SSN`,
   true,
   `Joe`,
   `Bob`,
   `jBob@test.com`,
   `555-565-3838`,
   ``,
   `555-282-9878`,
   `https://avatars.com/joeBob`,
   as)
}

type reqHandler func (http.ResponseWriter, *http.Request)

func requestAndCheck(t *testing.T, req *http.Request, h reqHandler, expectedCode int) *http.Response {
  rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedCode)
	}
  rr.Flush()
  return rr.Result()
}

func (s *PersonAPIIntegrationSuite) requestAndCheck2(req *http.Request, expectedCode int) *http.Response {
  rr := httptest.NewRecorder()
  s.Router.ServeHTTP(rr, req)

  require.Equal(s.T(), rr.Code, expectedCode,
    "handler returned wrong status code: got %v want %v",
    rr.Code, expectedCode)

  rr.Flush()
  return rr.Result()
}

func init() {
  terror.EchoErrorLog()
  os.Setenv(`ALLOW_UNSAFE_STATE_CHANGES`, `true`)
}

type PersonAPIIntegrationSuite struct {
  suite.Suite
  Router *mux.Router
}
func (s *PersonAPIIntegrationSuite) SetupTest() {
  s.Router = mux.NewRouter()
  api.InitAPI(s.Router)
}
func TestPersonAPIIntegrationSuite(t *testing.T) {
  if os.Getenv(`SKIP_INTEGRATION`) == `true` {
    t.Skip()
  } else {
    suite.Run(t, new(PersonAPIIntegrationSuite))
  }
}

func TestCreatePersonNoAuthentication(t *testing.T) {
	req, err := http.NewRequest("POST", "/persons", nil)
	if err != nil { t.Fatal(err) }

  requestAndCheck(t, req, api.CreateHandler, http.StatusUnauthorized)
}

func TestCreatePersonValid(t *testing.T) {
  authID := strkit.RandString(strkit.LettersAndNumbers, 16)

  controller := NewController(t)
  defer controller.Finish()
  authOracle := authmock.NewMockAuthOracle(controller)
  authOracle.EXPECT().GetAuthID().Return(authID).AnyTimes()
  authOracle.EXPECT().IsRequestAuthenticated().Return(true).AnyTimes()

  ctx := auth.SetAuthOracleOnContext(authOracle, context.Background())

  payload := joeBobJSON(authID)

	req, err := http.NewRequest("POST", "/persons", bytes.NewBuffer(payload))
	if err != nil { t.Fatal(err) }
  req = req.WithContext(ctx)

  requestAndCheck(t, req, api.CreateHandler, http.StatusOK)
}

func TestCreatePersonNonSelf(t *testing.T) {
  authID1 := strkit.RandString(strkit.LettersAndNumbers, 16)
  authID2 := strkit.RandString(strkit.LettersAndNumbers, 16)

  controller := NewController(t)
  defer controller.Finish()
  authOracle := authmock.NewMockAuthOracle(controller)
  authOracle.EXPECT().GetAuthID().Return(authID1).AnyTimes()
  authOracle.EXPECT().IsRequestAuthenticated().Return(true).AnyTimes()

  ctx := auth.SetAuthOracleOnContext(authOracle, context.Background())

  payload := joeBobJSON(authID2)

	req, err := http.NewRequest("POST", "/persons", bytes.NewBuffer(payload))
	if err != nil { t.Fatal(err) }
  req = req.WithContext(ctx)

  requestAndCheck(t, req, api.CreateHandler, http.StatusForbidden)
}

func TestGetSelfNonAuthenticated(t *testing.T) {
  req, err := http.NewRequest("GET", "/persons/auth-id-abcd1234", nil)
  if err != nil { t.Fatal(err) }

  requestAndCheck(t, req, api.CreateHandler, http.StatusUnauthorized)
}

func (s *PersonAPIIntegrationSuite) TestGetSelfValid() {
  authID := strkit.RandString(strkit.LettersAndNumbers, 16)
  joeBob := joeBobStruct(authID)

  controller := NewController(s.T())
  defer controller.Finish()
  authOracle := authmock.NewMockAuthOracle(controller)
  authOracle.EXPECT().GetAuthID().Return(authID).AnyTimes()
  authOracle.EXPECT().IsRequestAuthenticated().Return(true).AnyTimes()

  ctx := auth.SetAuthOracleOnContext(authOracle, context.Background())
  joeBob.CreateRaw(ctx)

  url := fmt.Sprintf("/persons/auth-id-%s/", authID)
  req := httptest.NewRequest("GET", url, nil)
  req = req.WithContext(ctx)

  res := s.requestAndCheck2(req, http.StatusOK)

  bodyJSON := make([]byte, 0, 2048)
  buffer := make([]byte, 2048, 2048)
  for _, err := res.Body.Read(buffer); err != nil && err != io.EOF; _, err = res.Body.Read(buffer) {
    require.NoErrorf(s.T(), err, `Reader failed with error: %s`, err)
    bodyJSON = append(bodyJSON, buffer...)
  }
  // log.Printf("JSON:\n\n%s\n\n", string(bodyJSON))
  // TODO: verify match
}

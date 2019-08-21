package persons_test

import (
  "bytes"
  "context"
  // "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"

  . "github.com/golang/mock/gomock"

  "github.com/Liquid-Labs/lc-authentication-api/go/auth"
  // . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  // . "github.com/Liquid-Labs/lc-locations-model/go/locations"
  authmock "github.com/Liquid-Labs/lc-authentication-api/go/mock"
  api "github.com/Liquid-Labs/lc-persons-api/go/persons"
  //. "github.com/Liquid-Labs/lc-persons-model/go/persons"
  "github.com/Liquid-Labs/strkit/go/strkit"
  "github.com/Liquid-Labs/terror/go/terror"
)

func init() {
  terror.EchoErrorLog()
}

func TestCreatePersonNoAuthentication(t *testing.T) {
	req, err := http.NewRequest("CREATE", "/persons", nil)
	if err != nil { t.Fatal(err) }

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.CreateHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

// func TestCreatePersonAuthenticatedDifferentUser

func TestCreatePersonValid(t *testing.T) {
  authID := strkit.RandString(strkit.LettersAndNumbers, 16)

  controller := NewController(t)
  defer controller.Finish()
  authOracle := authmock.NewMockAuthOracle(controller)
  authOracle.EXPECT().GetAuthID().Return(authID).AnyTimes()
  authOracle.EXPECT().IsRequestAuthenticated().Return(true).AnyTimes()

  ctx := auth.SetAuthOracleOnContext(authOracle, context.Background())


  /*a1 := NewAddress(`Camelot`, `a house`, EID(``), false, `100 Main Str`, `#B`, `Paris`, `TX`, `78383`, EID(``), `home`)
  as := Addresses{a1}
  p := NewPerson(`Joe Bob`,
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
  payload, err := json.Marshal(&p)
  log.Printf("\n\n%s\n\n", payload)
  if err != nil { t.Fatalf(`Failed to marshal test data: %s`, err) }*/

  payload := []byte(`{
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

	req, err := http.NewRequest("POST", "/persons", bytes.NewBuffer(payload))
	if err != nil { t.Fatal(err) }
  req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.CreateHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

/*
	// Check the response body is what we expect.
	expected := `[{"id":1,"first_name":"Krish","last_name":"Bhanushali","email_address":"krishsb@g.com","phone_number":"0987654321"},{"id":2,"first_name":"xyz","last_name":"pqr","email_address":"xyz@pqr.com","phone_number":"1234567890"},{"id":6,"first_name":"FirstNameSample","last_name":"LastNameSample","email_address":"lr@gmail.com","phone_number":"1111111111"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}*/

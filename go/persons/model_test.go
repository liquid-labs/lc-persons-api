package persons_test

import (
  "encoding/json"
  "strconv"
  "strings"
  "testing"

  . "github.com/Liquid-Labs/catalyst-persons-api/go/persons"
  "github.com/Liquid-Labs/go-nullable-mysql/nulls"
  "github.com/stretchr/testify/assert"
)

const jdDisplayName = "John Doe"
const jdEmail = "johndoe@test.com"
const jdPhone = "555-555-0000"
const jdActive = true

var johnDoeJson string = `
  {
    "displayName": "` + jdDisplayName + `",
    "email": "` + jdEmail + `",
    "phone": "` + jdPhone + `",
    "active": ` + strconv.FormatBool(jdActive) + `
  }`

var decoder *json.Decoder = json.NewDecoder(strings.NewReader(johnDoeJson))
var johnDoePerson = &Person{}
var decodeErr = decoder.Decode(johnDoePerson)

func TestPersonsDecode(t *testing.T) {
  assert.NoError(t, decodeErr, "Unexpected error decoding person JSON.")
  assert.Equal(t, jdDisplayName, johnDoePerson.DisplayName.String, "Unexpected display name.")
  assert.Equal(t, jdEmail, johnDoePerson.Email.String, "Unexpected email.")
  assert.Equal(t, jdPhone, johnDoePerson.Phone.String, "Unexpected phone.")
  assert.Equal(t, jdActive, johnDoePerson.Active.Bool, "Unexpected active value.")
}

func TestPersonFormatter(t *testing.T) {
  testP := &Person{PersonSummary: PersonSummary{
    Phone: nulls.NewString(`5555555555`),
    PhoneBackup: nulls.NewString(`1234567890`),
  }}
  testP.FormatOut()
  assert.Equal(t, `555-555-5555`, testP.Phone.String)
  assert.Equal(t, `123-456-7890`, testP.PhoneBackup.String)
}

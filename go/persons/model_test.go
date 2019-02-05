package persons_test

import (
  "encoding/json"
  "strings"
  "testing"

  . "github.com/Liquid-Labs/catalyst-persons-api/go/persons"
  "github.com/stretchr/testify/assert"
)

const johnDoeDisplayName = "John Doe"
const johnDoeEmail = "johndoe@test.com"
const johnDoePhone = "555-555-0000"
const johnDoeJson = `
  {
    "displayName": "` + johnDoeDisplayName + `",
    "email": "` + johnDoeEmail + `",
    "phone": "` + johnDoePhone + `"
  }`

var decoder *json.Decoder = json.NewDecoder(strings.NewReader(johnDoeJson))
var johnDoePerson = &Person{}
var decodeErr = decoder.Decode(johnDoePerson)

func TestPersonsDecode(t *testing.T) {
  assert.NoError(t, decodeErr, "Unexpected error decoding person JSON.")
}

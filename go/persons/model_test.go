package persons_test

import (
  "encoding/json"
  "reflect"
  "strconv"
  "strings"
  "testing"

  . "github.com/Liquid-Labs/catalyst-persons-api/go/persons"
  "github.com/Liquid-Labs/catalyst-core-api/go/entities"
  "github.com/Liquid-Labs/catalyst-core-api/go/locations"
  "github.com/Liquid-Labs/catalyst-core-api/go/users"
  "github.com/Liquid-Labs/go-nullable-mysql/nulls"
  "github.com/stretchr/testify/assert"
)

var trivialPersonSummary = &PersonSummary{
  users.User{
    entities.Entity{
      nulls.NewInt64(1),
      nulls.NewString(`a`),
      nulls.NewInt64(2),
    },
    nulls.NewBool(false),
  },
  nulls.NewString(`displayName`),
  nulls.NewString(`foo@test.com`),
  nulls.NewString(`555-555-9999`),
  nulls.NewString(`555-555-9998`),
  nulls.NewString(`http://foo.com/avatar`),
}

func TestPersonSummaryClone(t *testing.T) {
  clone := trivialPersonSummary.Clone()
  assert.Equal(t, trivialPersonSummary, clone, `Original does not match clone.`)
  clone.Id = nulls.NewInt64(3)
  clone.PubId = nulls.NewString(`b`)
  clone.LastUpdated = nulls.NewInt64(4)
  clone.Active = nulls.NewBool(true)
  clone.DisplayName = nulls.NewString(`different name`)
  clone.Email = nulls.NewString(`blah@test.com`)
  clone.Phone = nulls.NewString(`555-555-9997`)
  clone.PhoneBackup = nulls.NewString(`555-555-9996`)
  clone.PhotoURL = nulls.NewString(`http://bar.com/image`)

  oReflection := reflect.ValueOf(trivialPersonSummary).Elem()
  cReflection := reflect.ValueOf(clone).Elem()
  for i := 0; i < oReflection.NumField(); i++ {
    assert.NotEqualf(
      t,
      oReflection.Field(i).Interface(),
      cReflection.Field(i).Interface(),
      `Fields '%s' unexpectedly match.`,
      oReflection.Type().Field(i),
    )
  }
}

var trivialPerson = &Person{
  *trivialPersonSummary,
  locations.Addresses{
    &locations.Address{
      locations.Location{
        nulls.NewInt64(1),
        nulls.NewString(`a`),
        nulls.NewString(`b`),
        nulls.NewString(`c`),
        nulls.NewString(`d`),
        nulls.NewString(`e`),
        nulls.NewFloat64(2.0),
        nulls.NewFloat64(3.0),
        []string{`f`, `g`},
      },
      nulls.NewInt64(1),
      nulls.NewString(`label a`),
    },
  },
  []string{`h`, `i`},
}

func TestPersonClone(t *testing.T) {
  clone := trivialPerson.Clone()
  assert.Equal(t, trivialPerson, clone, `Original does not match clone.`)
  clone.Id = nulls.NewInt64(3)
  clone.PubId = nulls.NewString(`b`)
  clone.LastUpdated = nulls.NewInt64(4)
  clone.Active = nulls.NewBool(true)
  clone.DisplayName = nulls.NewString(`different name`)
  clone.Email = nulls.NewString(`blah@test.com`)
  clone.Phone = nulls.NewString(`555-555-9997`)
  clone.PhoneBackup = nulls.NewString(`555-555-9996`)
  clone.PhotoURL = nulls.NewString(`http://bar.com/image`)
  clone.Addresses = locations.Addresses{
    &locations.Address{
      locations.Location{
        nulls.NewInt64(2),
        nulls.NewString(`z`),
        nulls.NewString(`y`),
        nulls.NewString(`x`),
        nulls.NewString(`w`),
        nulls.NewString(`u`),
        nulls.NewFloat64(4.0),
        nulls.NewFloat64(5.0),
        []string{`i`},
      },
      nulls.NewInt64(2),
      nulls.NewString(`label b`),
    },
  }
  clone.ChangeDesc = []string{`j`}

  assert.NotEqual(t, trivialPerson.Addresses, clone.Addresses, `Addresses unexpectedly equal.`)
  aoReflection := reflect.ValueOf(trivialPerson.Addresses[0]).Elem()
  acReflection := reflect.ValueOf(clone.Addresses[0]).Elem()
  for i := 0; i < aoReflection.NumField(); i++ {
    assert.NotEqualf(
      t,
      aoReflection.Field(i).Interface(),
      acReflection.Field(i).Interface(),
      `Fields '%s' unexpectedly match.`,
      aoReflection.Type().Field(i),
    )
  }

  oReflection := reflect.ValueOf(trivialPerson).Elem()
  cReflection := reflect.ValueOf(clone).Elem()
  for i := 0; i < oReflection.NumField(); i++ {
    assert.NotEqualf(
      t,
      oReflection.Field(i).Interface(),
      cReflection.Field(i).Interface(),
      `Fields '%s' unexpectedly match.`,
      oReflection.Type().Field(i),
    )
  }
}

const jdDisplayName = "John Doe"
const jdEmail = "johndoe@test.com"
const jdPhone = "555-555-0000"
const jdActive = false

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

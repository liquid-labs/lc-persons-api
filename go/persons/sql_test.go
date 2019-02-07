package persons_test

import (
  "context"
  "os"
  "testing"

  // the package we're testing
  . "github.com/Liquid-Labs/catalyst-persons-api/go/persons"
  "github.com/Liquid-Labs/catalyst-core-api/go/entities"
  "github.com/Liquid-Labs/catalyst-core-api/go/locations"
  "github.com/Liquid-Labs/catalyst-core-api/go/users"
  "github.com/Liquid-Labs/go-api/sqldb"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
)

func TestPersonsDBIntegration(t *testing.T) {
  if os.Getenv(`SKIP_INTEGRATION`) == `true` {
    t.Skip()
  }

  if johnDoePerson == nil {
    t.Error(`Person struct not define; can't continue. This probbaly indicates a setup failure in 'model_test.go'.`)
  } else {
    if t.Run(`PersonsDBSetup`, testPersonDBSetup) {
      if sqldb.DB == nil { // test was skipped, but we still need to setup
        setupDB()
      }
      t.Run(`PersonGet`, testPersonGet)
      t.Run(`PersonCreate`, testPersonCreate)
      t.Run(`PersonUpdate`, testPersonUpdate)
      // t.Run(`PersonGetInTxn`, testPersonGetInTxn)
      t.Run(`PersonCreateInTxn`, testPersonCreateInTxn)
      t.Run(`PersonUpdateInTxn`, testPersonUpdateInTxn)
    }
  }
}

const janeDoeId=`4BE66BE5-2A62-11E9-B987-42010A8003FF`

func setupDB() {
  sqldb.RegisterSetup(entities.SetupDB, locations.SetupDB, users.SetupDB, /*persons.*/SetupDB)
  sqldb.InitDB() // panics if unable to initialize
}

func testPersonDBSetup(t *testing.T) {
  setupDB()
}

func testPersonGet(t *testing.T) {
  person, err := GetPerson(janeDoeId, context.Background())
  require.NoError(t, err, `Unexpected error getting Person.`)
  require.NotNil(t, person, `Unexpected nil Person on create (with no error).`)
  assert.Equal(t, `Jane Doe`, person.DisplayName.String, `Unexpected display name.`)
  assert.Equal(t, `janedoe@test.com`, person.Email.String, `Unexpected email.`)
  assert.Equal(t, `555-555-1111`, person.Phone.String, `Unexpected phone.`)
  assert.Equal(t, false, person.Active.Bool, `Unexpected active value.`)
  assert.NotEmpty(t, person.Id, `Unexpected empty ID.`)
  assert.Equal(t, janeDoeId, person.PubId.String, `Unexpected public id.`)
}

func testPersonCreate(t *testing.T) {
  person, err := CreatePerson(johnDoePerson, context.Background())
  require.NoError(t, err, `Unexpected error creating Person.`)
  require.NotNil(t, person, `Unexpected nil Person on create (with no error).`)
  assert.Equal(t, johnDoePerson.DisplayName, person.DisplayName, `Unexpected display name.`)
  assert.Equal(t, johnDoePerson.Email, person.Email, `Unexpected email.`)
  assert.Equal(t, johnDoePerson.Phone, person.Phone, `Unexpected phone.`)
  assert.Equal(t, johnDoePerson.Active, person.Active, `Unexpected active value.`)
  assert.NotEmpty(t, person.Id, `Unexpected empty ID.`)
  assert.NotEmpty(t, person.PubId, `Unexpected empty public id.`)
}

func testPersonUpdate(t *testing.T) {
  janeDoePerson, err := GetPerson(janeDoeId, context.Background())
  require.NoError(t, err, `Unexpected error getting Person.`)
  janeDoePerson.SetActive(true)
  janeDoePerson.SetDisplayName(`Jane P. Doe`)
  janeDoePerson.SetEmail(`janepdoe@test.com`)
  janeDoePerson.SetPhone(`555-555-0001`)
  janeDoePerson.SetPhoneBackup(`555-555-0002`)
  person, err := UpdatePerson(janeDoePerson, context.Background())
  require.NoError(t, err, `Unexpected error updating Person.`)
  require.NotNil(t, person, `Unexpected nil Person on create (with no error).`)
  assert.Equal(t, janeDoePerson.DisplayName, person.DisplayName, `Unexpected display name.`)
  assert.Equal(t, janeDoePerson.Email, person.Email, `Unexpected email.`)
  assert.Equal(t, janeDoePerson.Phone, person.Phone, `Unexpected phone.`)
  assert.Equal(t, janeDoePerson.Active, person.Active, `Unexpected active value.`)
  assert.NotEmpty(t, person.Id, `Unexpected empty ID.`)
  assert.NotEmpty(t, person.PubId, `Unexpected empty public id.`)
}
/*
func testPersonGetInTxn(t *testing.T) {
  janeDoePerson, err := GetPerson(janeDoeId, context.Background)
  assert.NoError(t, err, `Unexpected error getting person.`)
  txn, err := sqldb.DB.Begin()
  // if we get in a txn, we should see the changes
  janeDoePerson.SetPhone(`555-555-0003`)
  person, err := UpdatePersonInTxn(janeDoePerson, context.Background(), txn)

  assert.NoError(t, err, `Unexpected error opening transaction.`)
  janeDoePerson, err := GetPersonInTxn(janeDoeId, context.Background, txn)
  assert.NoError(t, err, `Unexpected error getting person.`)

  assert.NotEqual(t, janeDoePerson.Phone, person.Phone, `Phone number unexpectedly the same.`)
  assert.NoError(t, txn.Rollback())
}
*/
func testPersonCreateInTxn(t *testing.T) {
  /*txn, err := sqldb.DB.Begin()
  assert.NoError(t, err, `Unexpected error opening transaction.`)*/
}

func testPersonUpdateInTxn(t *testing.T) {
  /*txn, err := sqldb.DB.Begin()
  assert.NoError(t, err, `Unexpected error opening transaction.`)*/
}

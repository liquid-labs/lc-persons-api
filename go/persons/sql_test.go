package persons_test

import (
  "context"
  "os"
  "testing"

  // the package we're testing
  . "github.com/Liquid-Labs/catalyst-persons-api/go/persons"
  "github.com/Liquid-Labs/catalyst-core-api/go/entities"
  "github.com/Liquid-Labs/catalyst-core-api/go/users"
  "github.com/Liquid-Labs/go-api/sqldb"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
)

func testPersonCreate(t *testing.T) {
  person, err := CreatePerson(johnDoePerson, context.Background())
  require.NoError(t, err, "Unexpected error creating Person.")
  require.NotNil(t, person, "Unexpected nil Person on create (with no error).")
  assert.Equal(t, johnDoePerson.DisplayName, person.DisplayName, "Unexpected display name.")
  assert.Equal(t, johnDoePerson.Email, person.Email, "Unexpected email.")
  assert.Equal(t, johnDoePerson.Phone, person.Phone, "Unexpected phone.")
  assert.NotEmpty(t, person.Id, "Unexpected empty ID.")
  assert.NotEmpty(t, person.PubId, "Unexpected empty public id.")
}

func testPersonDBSetup(t *testing.T) {
  sqldb.RegisterSetup(entities.SetupDB, users.SetupDB, /*persons.*/SetupDB)
  sqldb.InitDB() // panics if unable to initialize
}

func TestPersonsDBIntegration(t *testing.T) {
  if os.Getenv("SKIP_INTEGRATION") == "true" {
    t.Skip()
  }

  if johnDoePerson == nil {
    t.Error("Person struct not define; can't continue. This probbaly indicates a setup failure in 'model_test.go'.")
  } else {
    if t.Run("PersonsDBSetup", testPersonDBSetup) {
      t.Run("PersonCreate", testPersonCreate)
    }
  }
}

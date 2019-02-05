package persons_test

import (
  "context"
  "strings"
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
  assert.Equal(t, johnDoeDisplayName, person.DisplayName.String, "Unexpected display name.")
  assert.Equal(t, johnDoeEmail, person.Email.String, "Unexpected email.")
  phoneFormatter := strings.NewReplacer("-", "", ".", "", "(", "", ")", "")
  assert.Equal(t, phoneFormatter.Replace(johnDoePhone), person.Phone.String, "Unexpected phone.")
  assert.NotEmpty(t, person.Id, "Unexpected empty ID.")
  assert.NotEmpty(t, person.PubId, "Unexpected empty public id.")
}

func TestPersonsSqlSuite(t *testing.T) {
  sqldb.RegisterSetup(entities.SetupDB, users.SetupDB, /*persons.*/SetupDB)
  sqldb.InitDB() // panics if unabel to initialize
  t.Run("PersonCreate", testPersonCreate)
}

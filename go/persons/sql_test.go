package persons

import (
  "encoding/json"
  "log"
  "strings"
  "testing"

  "github.com/Liquid-Labs/go-api/sqldb"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/suite"
)

type PersonsSqlTestSuite struct {
  suite.Suite
  JohnDoe *Person
}

const JohnDoeJson = `
  {
    "displayName": "John Doe",
    "email": "johndoe@test.com",
    "phone": "555-555-5555"
  }`

func (suite *PersonsSqlTestSuite) SetupSuite() {
  sqldb.InitDb()

  suite.JohnDoe = &Person{}
  decoder := json.NewDecoder(strings.NewReader(JohnDoeJson))
  if err := decoder.Decode(suite.JohnDoe); err != nil {
    log.Panicf("Could not create John Doe: %s", err)
  }
}

func (suite *PersonsSqlTestSuite) TestPersonCreate() {
  assert.NotNil(suite.T(), suite.JohnDoe)
}

func TestPersonsSqlTestSuite(t *testing.T) {
    suite.Run(t, &PersonsSqlTestSuite{})
}

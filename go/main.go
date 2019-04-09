package main

import (
  "github.com/Liquid-Labs/catalyst-core-api/go/restserv"
  // core resources
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/entities"
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/locations"
  "github.com/Liquid-Labs/catalyst-core-api/go/users"

  "github.com/Liquid-Labs/catalyst-persons-api/go/persons"
  "github.com/Liquid-Labs/go-api/sqldb"
)

func main() {
  sqldb.RegisterSetup(entities.SetupDB)
  sqldb.RegisterSetup(locations.SetupDB)
  sqldb.RegisterSetup(users.SetupDB)
  sqldb.RegisterSetup(persons.SetupDB)
  sqldb.InitDB()
  restserv.RegisterResource(persons.InitAPI)
  restserv.Init()
}

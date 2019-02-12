package main

import (
  "github.com/Liquid-Labs/catalyst-core-api/go/restserv"
  "github.com/Liquid-Labs/catalyst-persons-api/go/persons"
  "github.com/Liquid-Labs/go-api/sqldb"
)

func main() {
  sqldb.RegisterSetup(persons.SetupDB)
  restserv.RegisterResource(persons.InitAPI)
  restserv.Init()
}

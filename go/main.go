package main

import (
  "github.com/Liquid-Labs/catalyst-core-api/go/restserv"
  "github.com/Liquid-Labs/catalyst-persons-api/go/persons"
)

func main() {
  restserv.RegisterResource(persons.InitDB, persons.InitAPI)
  restserv.Init()
}

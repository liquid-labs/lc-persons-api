package main

import (
  "github.com/Liquid-Labs/catalyst-core-api/go/restserv"

  "github.com/Liquid-Labs/lc-persons-api/go/resources/persons"
)

func main() {
  restserv.RegisterResource(persons.InitAPI)
  restserv.Init()
}

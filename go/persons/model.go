package persons

import (
  "github.com/Liquid-Labs/catalyst-core-api/go/entities"
  "github.com/Liquid-Labs/catalyst-core-api/go/locations"
  "github.com/Liquid-Labs/go-nullable-mysql/nulls"
)

// On summary, we don't include address. Note leaving it empty and using
// 'omitempty' on the Person struct won't work because then Persons without
// an address will appear 'incomplete' in the front-end model and never resolve.
type PersonSummary struct {
  entities.Entity
  DisplayName  nulls.String `json:"name"`
  Email         nulls.String `json:"email"`
  Phone         nulls.String `json:"phone,string"`
  PhoneBackup   nulls.String `json:"phoneBackup,string"`
  PhotoURL      nulls.String `json:"photoUrl"`
}

// We expect an empty address array if no addresses on detail
type Person struct {
  PersonSummary
  Addresses     locations.Addresses  `json:"addresses"`
  Active        nulls.Bool           `json:"active"`
  ChangeDesc    []string             `json:"changeDesc,omitempty"`
}

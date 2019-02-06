package persons

import (
  "regexp"

  "github.com/Liquid-Labs/catalyst-core-api/go/users"
  "github.com/Liquid-Labs/catalyst-core-api/go/locations"
  "github.com/Liquid-Labs/go-nullable-mysql/nulls"
)

// On summary, we don't include address. Note leaving it empty and using
// 'omitempty' on the Person struct won't work because then Persons without
// an address will appear 'incomplete' in the front-end model and never resolve.
type PersonSummary struct {
  users.User
  DisplayName   nulls.String `json:"displayName"`
  Email         nulls.String `json:"email"`
  Phone         nulls.String `json:"phone,string"`
  PhoneBackup   nulls.String `json:"phoneBackup,string"`
  PhotoURL      nulls.String `json:"photoUrl"`
}

// We expect an empty address array if no addresses on detail
type Person struct {
  PersonSummary
  Addresses     locations.Addresses  `json:"addresses"`
  ChangeDesc    []string             `json:"changeDesc,omitempty"`
}

var phoneOutFormatter *regexp.Regexp = regexp.MustCompile(`^(\d{3})(\d{3})(\d{4})$`)

func (p *PersonSummary) FormatOut() {
  p.Phone.String = phoneOutFormatter.ReplaceAllString(p.Phone.String, `$1-$2-$3`)
  p.PhoneBackup.String = phoneOutFormatter.ReplaceAllString(p.PhoneBackup.String, `$1-$2-$3`)
}

func (p *Person) FormatOut() {
  p.PersonSummary.FormatOut()
}

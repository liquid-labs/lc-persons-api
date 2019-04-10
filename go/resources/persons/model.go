package persons

import (
  "regexp"

  "github.com/Liquid-Labs/catalyst-core-api/go/resources/users"
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/locations"
  "github.com/Liquid-Labs/catalyst-core-api/go/resources"
  "github.com/Liquid-Labs/go-nullable-mysql/nulls"
)

var phoneOutFormatter *regexp.Regexp = regexp.MustCompile(`^(\d{3})(\d{3})(\d{4})$`)

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

func (p *PersonSummary) FormatOut() {
  p.Phone.String = phoneOutFormatter.ReplaceAllString(p.Phone.String, `$1-$2-$3`)
  p.PhoneBackup.String = phoneOutFormatter.ReplaceAllString(p.PhoneBackup.String, `$1-$2-$3`)
}

func (p *PersonSummary) SetDisplayName(val string) {
  p.DisplayName = nulls.NewString(val)
}

func (p *PersonSummary) SetEmail(val string) {
  p.Email = nulls.NewString(val)
}

func (p *PersonSummary) SetPhone(val string) {
  p.Phone = nulls.NewString(val)
}

func (p *PersonSummary) SetPhoneBackup(val string) {
  p.PhoneBackup = nulls.NewString(val)
}

func (p *PersonSummary) SetPhotoURL(val string) {
  p.PhotoURL = nulls.NewString(val)
}

func (p *PersonSummary) Clone() *PersonSummary {
  return &PersonSummary{
    *p.User.Clone(),
    p.DisplayName,
    p.Email,
    p.Phone,
    p.PhoneBackup,
    p.PhotoURL,
  }
}

// We expect an empty address array if no addresses on detail
type Person struct {
  PersonSummary
  Addresses     locations.Addresses  `json:"addresses"`
  ChangeDesc    []string             `json:"changeDesc,omitempty"`
}

func (p *Person) Clone() *Person {
  newChangeDesc := make([]string, len(p.ChangeDesc))
  copy(newChangeDesc, p.ChangeDesc)

  return &Person{
    *p.PersonSummary.Clone(),
    *p.Addresses.Clone(),
    newChangeDesc,
  }
}

func (p *Person) PromoteChanges() {
  p.ChangeDesc = resources.PromoteChanges(p.Addresses, p.ChangeDesc)
}

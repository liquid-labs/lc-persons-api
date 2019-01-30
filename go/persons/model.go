package persons

type Person struct {
  UID           string   `json:"uid"`
  DisplayName   string   `json:"displayName"`
  Email         string   `json:"email"`
  PhoneNumber   string   `json:"phoneNumber"`
  PhotoURL      string   `json:"photoUrl"`
}

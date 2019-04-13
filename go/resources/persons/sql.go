package persons

import (
  "context"
  "database/sql"
  "fmt"
  "log"
  "strconv"

  "github.com/Liquid-Labs/go-api/sqldb"
  "github.com/Liquid-Labs/go-nullable-mysql/nulls"
  "github.com/Liquid-Labs/go-rest/rest"
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/users"
  "github.com/Liquid-Labs/catalyst-core-api/go/resources/locations"
)

var PersonsSorts = map[string]string{
  "": `p.display_name ASC `,
  `name-asc`: `p.display_name ASC `,
  `name-desc`: `p.display_name DESC `,
}

func ScanPersonSummary(row *sql.Rows) (*PersonSummary, error) {
	var p PersonSummary

	if err := row.Scan(&p.PubId, &p.LastUpdated, &p.DisplayName, &p.Phone, &p.Email, &p.PhoneBackup); err != nil {
		return nil, err
	}

	return &p, nil
}

func ScanPersonDetail(row *sql.Rows) (*Person, *locations.Address, error) {
  var p Person
  var a locations.Address

	if err := row.Scan(&p.PubId, &p.LastUpdated, &p.DisplayName, &p.Phone,
      &p.Email, &p.PhoneBackup, &p.Active, &p.AuthId, &p.LegalID, &p.LegalIDType,
      &p.Id, &a.LocationId, &a.Idx, &a.Label, &a.Address1, &a.Address2, &a.City,
      &a.State, &a.Zip, &a.Lat, &a.Lng); err != nil {
		return nil, nil, err
	}

  // Negative locationIds are used by the UID for temporary identification.
  if a.LocationId.Int64 < 0 {
    a.LocationId = nulls.NewNullInt64()
  }

	return &p, &a, nil
}

// implement rest.ResultBuilder
func BuildPersonResults(rows *sql.Rows) (interface{}, error) {
  results := make([]*PersonSummary, 0)
  for rows.Next() {
    person, err := ScanPersonSummary(rows)
    if err != nil {
      return nil, err
    }

    results = append(results, person)
  }

  return results, nil
}

// Implements rest.GeneralSearchWhereBit
func PersonsGeneralWhereGenerator(term string, params []interface{}) (string, []interface{}, error) {
  likeTerm := `%`+term+`%`
  var whereBit string
  if _, err := strconv.ParseInt(term,10,64); err == nil {
    whereBit += "AND (p.phone LIKE ? OR p.phone_backup LIKE ?) "
    params = append(params, likeTerm, likeTerm)
  } else {
    whereBit += "AND (p.display_name LIKE ? OR p.email LIKE ?) "
    params = append(params, likeTerm, likeTerm)
  }

  return whereBit, params, nil
}

const CommonPersonFields = `e.pub_id, e.last_updated, p.display_name, p.phone, p.email, p.phone_backup, u.active, u.auth_id, u.legal_id, u.legal_id_type `
const CommonPersonsFrom = `FROM persons p JOIN users u ON p.id=u.id JOIN entities e ON p.id=e.id `

const createPersonStatement = `INSERT INTO persons (id, display_name, phone, email, phone_backup) VALUES(?,?,?,?,?)`
func CreatePerson(p *Person, ctx context.Context) (*Person, rest.RestError) {
  txn, err := sqldb.DB.Begin()
  if err != nil {
    defer txn.Rollback()
    return nil, rest.ServerError("Could not create person record. (txn error)", err)
  }
  newP, restErr := CreatePersonInTxn(p, ctx, txn)
  // txn already rolled back if in error, so we only need to commit if no error
  if err == nil {
    defer txn.Commit()
  }
  return newP, restErr
}

func CreatePersonInTxn(p *Person, ctx context.Context, txn *sql.Tx) (*Person, rest.RestError) {
  p.Addresses.CompleteAddresses(ctx)

  var err error
  newId, restErr := users.CreateUserInTxn(&p.User, txn)
  if restErr != nil {
    defer txn.Rollback()
		return nil, restErr
  }

  p.Id = nulls.NewInt64(newId)

	_, err = txn.Stmt(createPersonQuery).Exec(newId, p.DisplayName, p.Phone, p.Email, p.PhoneBackup)
	if err != nil {
    // TODO: can we do more to tell the cause of the failure? We assume it's due to malformed data with the HTTP code
    defer txn.Rollback()
    log.Print(err)
		return nil, rest.UnprocessableEntityError("Failure creating person.", err)
	}

  if restErr := p.Addresses.CreateAddresses(nulls.NewInt64(newId), ctx, txn); restErr != nil {
    defer txn.Rollback()
    return nil, restErr
  }

  newPerson, err := GetPersonByIDInTxn(p.Id.Int64, ctx, txn)
  if err != nil {
    return nil, rest.ServerError("Problem retrieving newly updated person.", err)
  }
  // Carry any 'ChangeDesc' made by the geocoding out.
  p.PromoteChanges()
  newPerson.ChangeDesc = p.ChangeDesc

  return newPerson, nil
}

const CommonPersonGet string = `SELECT ` + CommonPersonFields + `, p.id, loc.id, ea.idx, ea.label, loc.address1, loc.address2, loc.city, loc.state, loc.zip, loc.lat, loc.lng ` + CommonPersonsFrom + ` LEFT JOIN entity_addresses ea ON p.id=ea.entity_id AND ea.idx >= 0 LEFT JOIN locations loc ON ea.location_id=loc.id `
const getPersonStatement string = CommonPersonGet + `WHERE e.pub_id=? `

// GetPerson retrieves a Person from a public ID string (UUID). Attempting to
// retrieve a non-existent Person results in a rest.NotFoundError. This is used
// primarily to retrieve a Person in response to an API request.
//
// Consider using GetPersonByID to retrieve a Person from another backend/DB
// function. TODO: reference discussion of internal vs public IDs.
func GetPerson(pubId string, ctx context.Context) (*Person, rest.RestError) {
  return getPersonHelper(getPersonQuery, pubId, ctx, nil)
}

// GetPersonInTxn retrieves a Person by public ID string (UUID) in the context
// of an existing transaction. See GetPerson.
func GetPersonInTxn(pubId string, ctx context.Context, txn *sql.Tx) (*Person, rest.RestError) {
  return getPersonHelper(getPersonQuery, pubId, ctx, txn)
}

const getPersonByAuthIdStatement string = CommonPersonGet + ` WHERE u.auth_id=? `
// GetPersonByAuthId retrieves a Person from a public authentication ID string
// provided by the authentication provider (firebase). Attempting to retrieve a
// non-existent Person results in a rest.NotFoundError. This is used primarily
// to retrieve a Person in response to an API request, especially
// '/persons/self'.
func GetPersonByAuthId(authId string, ctx context.Context) (*Person, rest.RestError) {
  return getPersonHelper(getPersonByAuthIdQuery, authId, ctx, nil)
}

// GetPersonByAuthIdInTxn retrieves a Person by public authentication ID string
// in the context of an existing transaction. See GetPersonByAuthId.
func GetPersonByAuthIdInTxn(authId string, ctx context.Context, txn *sql.Tx) (*Person, rest.RestError) {
  return getPersonHelper(getPersonByAuthIdQuery, authId, ctx, txn)
}

const getPersonByIdStatement string = CommonPersonGet + ` WHERE p.id=? `
// GetPersonByID retrieves a Person by internal ID. As the internal ID must
// never be exposed to users, this method is exclusively for internal/backend
// use. Specifically, since Persons are associated with other Entities through
// the internal ID (i.e., foreign keys use the internal ID), this function is
// most often used to retrieve a Person which is to be bundled in a response.
//
// Use GetPerson to retrieve a Person in response to an API request. TODO:
// reference discussion of internal vs public IDs.
func GetPersonByID(id int64, ctx context.Context) (*Person, rest.RestError) {
  return getPersonHelper(getPersonByIdQuery, id, ctx, nil)
}

// GetPersonByIDInTxn retrieves a Person by internal ID in the context of an
// existing transaction. See GetPersonByID.
func GetPersonByIDInTxn(id int64, ctx context.Context, txn *sql.Tx) (*Person, rest.RestError) {
  return getPersonHelper(getPersonByIdQuery, id, ctx, txn)
}

func getPersonHelper(stmt *sql.Stmt, id interface{}, ctx context.Context, txn *sql.Tx) (*Person, rest.RestError) {
  if txn != nil {
    stmt = txn.Stmt(stmt)
  }
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, rest.ServerError("Error retrieving person.", err)
	}
	defer rows.Close()

	var person *Person
  var address *locations.Address
  var addresses locations.Addresses = make(locations.Addresses, 0)
	for rows.Next() {
    var err error
    // The way the scanner works, it processes all the data each time. :(
    // 'person' gets updated with an equivalent structure while we gather up
    // the addresses.
    if person, address, err = ScanPersonDetail(rows); err != nil {
      return nil, rest.ServerError(fmt.Sprintf("Problem getting data for person: '%v'", id), err)
    }

    if address.LocationId.Valid {
	    addresses = append(addresses, address)
    }
	}
  if person != nil {
    person.Addresses = addresses
    person.FormatOut()
  } else {
    return nil, rest.NotFoundError(fmt.Sprintf(`Person '%s' not found.`, id), nil)
  }

	return person, nil
}

// BUG(zane@liquid-labs.com): UpdatePerson should use internal IDs if available
// on the Person struct. (I'm assuming this is slightly more efficient, though
// we should test.)

// UpdatesPerson updates the canonical Person record. Attempting to update a
// non-existent Person results in a rest.NotFoundError.
func UpdatePerson(p *Person, ctx context.Context) (*Person, rest.RestError) {
  txn, err := sqldb.DB.Begin()
  if err != nil {
    defer txn.Rollback()
    return nil, rest.ServerError("Could not update person record.", err)
  }

  newP, restErr := UpdatePersonInTxn(p, ctx, txn)
  // txn already rolled back if in error, so we only need to commit if no error
  if restErr == nil {
    defer txn.Commit()
  }

  return newP, restErr
}

// UpdatesPersonInTxn updates the canonical Person record within an existing
// transaction. See UpdatePerson.
func UpdatePersonInTxn(p *Person, ctx context.Context, txn *sql.Tx) (*Person, rest.RestError) {
  if p.Addresses != nil {
    p.Addresses.CompleteAddresses(ctx)
  }
  var err error
  if p.Addresses != nil {
    if restErr := p.Addresses.Update(p.PubId.String, ctx, txn); restErr != nil {
      defer txn.Rollback()
      // TODO: this message could be misleading; like the person was updated, and just the addresses not
      return nil, restErr
    }
  }

  var updateStmt *sql.Stmt = txn.Stmt(updatePersonQuery)
  _, err = updateStmt.Exec(p.Active, p.LegalID, p.LegalIDType, p.DisplayName, p.Phone, p.Email, p.PhoneBackup, p.PubId)
  if err != nil {
    if txn != nil {
      defer txn.Rollback()
    }
    return nil, rest.ServerError("Could not update person record.", err)
  }

  newPerson, err := GetPersonInTxn(p.PubId.String, ctx, txn)
  if err != nil {
    return nil, rest.ServerError("Problem retrieving newly updated person.", err)
  }
  // Carry any 'ChangeDesc' made by the geocoding out.
  p.PromoteChanges()
  newPerson.ChangeDesc = p.ChangeDesc

  return newPerson, nil
}

// TODO: enable update of AuthID
const updatePersonStatement = `UPDATE persons p JOIN users u ON u.id=p.id JOIN entities e ON p.id=e.id SET u.active=?, u.legal_id=?, u.legal_id_type=?, p.display_name=?, p.phone=?, p.email=?, p.phone_backup=?, e.last_updated=0 WHERE e.pub_id=?`
var createPersonQuery, updatePersonQuery, getPersonQuery, getPersonByAuthIdQuery, getPersonByIdQuery *sql.Stmt
func SetupDB(db *sql.DB) {
  var err error
  if createPersonQuery, err = db.Prepare(createPersonStatement); err != nil {
    log.Fatalf("mysql: prepare create person stmt:\n%v\n%s", err, createPersonStatement)
  }
  if getPersonQuery, err = db.Prepare(getPersonStatement); err != nil {
    log.Fatalf("mysql: prepare get person stmt: %v", err)
  }
  if getPersonByAuthIdQuery, err = db.Prepare(getPersonByAuthIdStatement); err != nil {
    log.Fatalf("mysql: prepare get person by auth ID stmt: %v", err)
  }
  if getPersonByIdQuery, err = db.Prepare(getPersonByIdStatement); err != nil {
    log.Fatalf("mysql: prepare get person by ID stmt: %v", err)
  }
  if updatePersonQuery, err = db.Prepare(updatePersonStatement); err != nil {
    log.Fatalf("mysql: prepare update person stmt: %v", err)
  }
}

package user

import "github.com/lnk.by/shared/service"

type User struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	OrganizationID string `json:"organizationId"`
}

func (u *User) FieldsPtrs() []any {
	return []any{&u.ID, &u.Email, &u.Name, &u.OrganizationID}
}

func (u *User) FieldsVals() []any {
	return []any{u.ID, u.Email, u.Name, u.OrganizationID}
}

var Ops = service.Ops[*User]{
	Create:   "INSERT INTO users (id, email, name, organization_id) VALUES ($1, $2, $3, $4)",
	Retrieve: "SELECT id, email, name, organization_id FROM users WHERE id = $1",
	Update:   "UPDATE users SET email = $2, name = $3, organization_id = $4 WHERE id = $1",
	Delete:   "DELETE FROM users WHERE id = $1",
}

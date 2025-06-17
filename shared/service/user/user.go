package user

import (
	"github.com/lnk.by/shared/service"
)

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

const UserIdParam = "userId"

var (
	CreateSQL   service.CreateSQL[*User]   = "INSERT INTO user (id, email, name, organization_id) VALUES ($1, $2, $3, $4)"
	RetrieveSQL service.RetrieveSQL[*User] = "SELECT id, email, name, organization_id FROM user WHERE id = $1"
	UpdateSQL   service.UpdateSQL[*User]   = "UPDATE user SET email = $2, name = $3, organization_id = $4 WHERE id = $1"
	DeleteSQL   service.DeleteSQL[*User]   = "DELETE FROM user WHERE id = $1"
	ListSQL     service.ListSQL[*User]     = "SELECT id, email, name, organization_id FROM user OFFSET $1 LIMIT $2"
)

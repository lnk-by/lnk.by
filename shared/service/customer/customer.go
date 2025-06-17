package customer

import (
	"errors"

	"github.com/lnk.by/shared/service"
)

type Customer struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	OrganizationID string `json:"organizationId"`
}

func (u *Customer) FieldsPtrs() []any {
	return []any{&u.ID, &u.Email, &u.Name, &u.OrganizationID}
}

func (u *Customer) FieldsVals() []any {
	return []any{u.ID, u.Email, u.Name, u.OrganizationID}
}

func (u *Customer) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}

	if u.Name == "" {
		return errors.New("name is required")
	}

	if u.ID != "" {
		return errors.New("user ID is managed by the server")
	}

	return nil
}

const IdParam = "customerId"

var (
	CreateSQL   service.CreateSQL[*Customer]   = "INSERT INTO customer (id, email, name, organization_id) VALUES ($1, $2, $3, $4)"
	RetrieveSQL service.RetrieveSQL[*Customer] = "SELECT id, email, name, organization_id FROM customer WHERE id = $1"
	UpdateSQL   service.UpdateSQL[*Customer]   = "UPDATE customer SET email = $2, name = $3, organization_id = $4 WHERE id = $1"
	DeleteSQL   service.DeleteSQL[*Customer]   = "DELETE FROM customer WHERE id = $1"
	ListSQL     service.ListSQL[*Customer]     = "SELECT id, email, name, organization_id FROM customer OFFSET $1 LIMIT $2"
)

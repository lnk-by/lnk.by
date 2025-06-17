package customer

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/utils"
)

type Customer struct {
	ID             string       `json:"id"`
	Email          string       `json:"email"`
	Name           string       `json:"name"`
	OrganizationID string       `json:"organizationId"`
	Status         utils.Status `json:"status"`
}

func (u *Customer) FieldsPtrs() []any {
	return []any{&u.ID, &u.Email, &u.Name, &u.OrganizationID, &u.Status}
}

func (u *Customer) FieldsVals() []any {
	return []any{u.ID, u.Email, u.Name, u.OrganizationID, u.Status}
}

func (u *Customer) WithId(id string) {
	u.ID = id
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

	id, err := uuid.NewV1()
	if err != nil {
		log.Fatalf("failed to generate UUIDv1: %v", err)
		return errors.New(fmt.Sprintf("failed to generate UUIDv1: %v", err))
	}
	u.ID = id.String()
	if u.Status == "" {
		u.Status = utils.StatusActive
	}

	return nil
}

const IdParam = "customerId"

var (
	CreateSQL   service.CreateSQL[*Customer]   = "INSERT INTO customer (id, email, name, organization_id, status) VALUES ($1, $2, $3, NULLIF($4, ''), $5)"
	RetrieveSQL service.RetrieveSQL[*Customer] = "SELECT id, email, name, COALESCE(organization_id, ''), status FROM customer WHERE id = $1 AND status='active'"
	UpdateSQL   service.UpdateSQL[*Customer]   = "UPDATE customer SET email = $2, name = $3, organization_id = NULLIF($4, ''), status = $5 WHERE id = $1"
	DeleteSQL   service.DeleteSQL[*Customer]   = "DELETE FROM customer WHERE id = $1"
	ListSQL     service.ListSQL[*Customer]     = "SELECT id, email, name, COALESCE(organization_id, ''), status FROM customer WHERE status='active' OFFSET $1 LIMIT $2"
)

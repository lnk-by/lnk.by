package customer

import (
	"errors"

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

func (c *Customer) FieldsPtrs() []any {
	return []any{&c.ID, &c.Email, &c.Name, &c.OrganizationID, &c.Status}
}

func (c *Customer) FieldsVals() []any {
	return []any{c.ID, c.Email, c.Name, c.OrganizationID, c.Status}
}

func (c *Customer) WithId(id string) {
	c.ID = id
}

func (c *Customer) Validate() error {
	switch {
	case c.Name == "":
		return service.ErrNameRequired
	case c.ID != "":
		return service.ErrIDManagedByServer
	case c.Email == "":
		return errors.New("email is required")
	default:
		return nil
	}
}

func (c *Customer) Generate() {
	c.ID = service.UUID()

	if c.Status == "" {
		c.Status = utils.StatusActive
	}
}

var (
	CreateSQL   service.CreateSQL[*Customer]   = "INSERT INTO customer (id, email, name, organization_id, status) VALUES ($1, $2, $3, NULLIF($4, ''), $5)"
	RetrieveSQL service.RetrieveSQL[*Customer] = "SELECT id, email, name, COALESCE(organization_id, ''), status FROM customer WHERE id = $1 AND status='active'"
	UpdateSQL   service.UpdateSQL[*Customer]   = "UPDATE customer SET email = $2, name = $3, organization_id = NULLIF($4, ''), status = $5 WHERE id = $1"
	DeleteSQL   service.DeleteSQL[*Customer]   = "DELETE FROM customer WHERE id = $1"
	ListSQL     service.ListSQL[*Customer]     = "SELECT id, email, name, COALESCE(organization_id, ''), status FROM customer WHERE status='active' OFFSET $1 LIMIT $2"
)

package campaign

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/utils"
)

type Campaign struct {
	ID             uuid.UUID    `json:"id"`
	Name           string       `json:"name"`
	ValidFrom      time.Time    `json:"validFrom"`
	ValidUntil     time.Time    `json:"validUntil"`
	OrganizationID *uuid.UUID   `json:"organizationId"`
	CustomerID     *uuid.UUID   `json:"customerId"`
	Status         utils.Status `json:"status"`
}

func (c *Campaign) FieldsPtrs() []any {
	return []any{&c.ID, &c.Name, &c.OrganizationID, &c.CustomerID, &c.Status}
}

func (c *Campaign) FieldsVals() []any {
	return []any{c.ID, c.Name, c.OrganizationID, c.CustomerID, c.Status}
}

func (c *Campaign) ParseID(idString string) (uuid.UUID, error) {
	return uuid.FromString(idString)
}

func (c *Campaign) WithID(id uuid.UUID) {
	c.ID = id
}

func (c *Campaign) Validate() error {
	switch {
	case c.Name == "":
		return service.ErrNameRequired
	case c.ID != uuid.Nil:
		return service.ErrIDManagedByServer
	default:
		return nil
	}
}

func (c *Campaign) Generate() {
	c.ID = service.UUID()

	if c.Status == "" {
		c.Status = utils.StatusActive
	}
}

var (
	CreateSQL   service.CreateSQL[*Campaign]   = "INSERT INTO campaign (id, name, organization_id, customer_id, status) VALUES ($1, $2, $3, $4, $5)"
	RetrieveSQL service.RetrieveSQL[*Campaign] = "SELECT id, name, organization_id, customer_id, status FROM campaign WHERE id = $1 AND status='active' AND now() BETWEEN valid_from AND valid_until"
	UpdateSQL   service.UpdateSQL[*Campaign]   = "UPDATE campaign SET name = $2, organization_id = $3, customer_id = $4, status = $5 WHERE id = $1"
	DeleteSQL   service.DeleteSQL[*Campaign]   = "DELETE FROM campaign WHERE id = $1"
	ListSQL     service.ListSQL[*Campaign]     = "SELECT id, name, organization_id, customer_id, status FROM campaign WHERE status='active' AND customer_id=$1 OFFSET $2 LIMIT $3"
)

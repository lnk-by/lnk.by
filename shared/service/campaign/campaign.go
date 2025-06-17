package campaign

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/utils"
)

type Campaign struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	OrganizationID string       `json:"organizationId"`
	CustomerID     string       `json:"customerId"`
	Status         utils.Status `json:"status"`
}

func (c *Campaign) FieldsPtrs() []any {
	return []any{&c.ID, &c.Name, &c.OrganizationID, &c.CustomerID, &c.Status}
}

func (c *Campaign) FieldsVals() []any {
	return []any{c.ID, c.Name, c.OrganizationID, c.CustomerID, c.Status}
}
func (u *Campaign) WithId(id string) {
	u.ID = id
}

func (c *Campaign) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}

	if c.ID != "" {
		return errors.New("user ID is managed by the server")
	}

	id, err := uuid.NewV1()
	if err != nil {
		log.Fatalf("failed to generate UUIDv1: %v", err)
		return errors.New(fmt.Sprintf("failed to generate UUIDv1: %v", err))
	}
	c.ID = id.String()
	if c.Status == "" {
		c.Status = utils.StatusActive
	}

	return nil
}

const IdParam = "campaignId"

var (
	CreateSQL   service.CreateSQL[*Campaign]   = "INSERT INTO campaign (id, name, organization_id, customer_id, status) VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), $5)"
	RetrieveSQL service.RetrieveSQL[*Campaign] = "SELECT id, name, COALESCE(organization_id, ''), COALESCE(customer_id, ''), status FROM campaign WHERE id = $1 AND status='active'"
	UpdateSQL   service.UpdateSQL[*Campaign]   = "UPDATE campaign SET name = $2, organization_id = NULLIF($3, ''), customer_id = NULLIF($4, ''), status = $5 WHERE id = $1"
	DeleteSQL   service.DeleteSQL[*Campaign]   = "DELETE FROM campaign WHERE id = $1"
	ListSQL     service.ListSQL[*Campaign]     = "SELECT id, name, COALESCE(organization_id, ''), COALESCE(customer_id, ''), status FROM campaign WHERE status='active' OFFSET $1 LIMIT $2"
)

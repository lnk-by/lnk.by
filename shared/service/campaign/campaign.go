package campaign

import (
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

func (c *Campaign) WithId(id string) {
	c.ID = id
}

func (c *Campaign) Validate() error {
	switch {
	case c.Name == "":
		return service.ErrNameRequired
	case c.ID != "":
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

const IdParam = "campaignId"

var (
	CreateSQL   service.CreateSQL[*Campaign]   = "INSERT INTO campaign (id, name, organization_id, customer_id, status) VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), $5)"
	RetrieveSQL service.RetrieveSQL[*Campaign] = "SELECT id, name, COALESCE(organization_id, ''), COALESCE(customer_id, ''), status FROM campaign WHERE id = $1 AND status='active'"
	UpdateSQL   service.UpdateSQL[*Campaign]   = "UPDATE campaign SET name = $2, organization_id = NULLIF($3, ''), customer_id = NULLIF($4, ''), status = $5 WHERE id = $1"
	DeleteSQL   service.DeleteSQL[*Campaign]   = "DELETE FROM campaign WHERE id = $1"
	ListSQL     service.ListSQL[*Campaign]     = "SELECT id, name, COALESCE(organization_id, ''), COALESCE(customer_id, ''), status FROM campaign WHERE status='active' OFFSET $1 LIMIT $2"
)

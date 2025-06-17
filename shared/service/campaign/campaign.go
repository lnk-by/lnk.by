package campaign

import "github.com/lnk.by/shared/service"

type Campaign struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Campaign) FieldsPtrs() []any {
	return []any{&c.ID, &c.Name}
}

func (c *Campaign) FieldsVals() []any {
	return []any{c.ID, c.Name}
}

var (
	CreateSQL   service.CreateSQL[*Campaign]   = "INSERT INTO organization (id, name) VALUES ($1, $2)"
	RetrieveSQL service.RetrieveSQL[*Campaign] = "SELECT id, name FROM organization WHERE id = $1"
	UpdateSQL   service.UpdateSQL[*Campaign]   = "UPDATE organization SET name = $2 WHERE id = $1"
	DeleteSQL   service.DeleteSQL[*Campaign]   = "DELETE FROM user WHERE id = $1"
	ListSQL     service.ListSQL[*Campaign]     = "SELECT id, name FROM organization OFFSET $1 LIMIT $2"
)

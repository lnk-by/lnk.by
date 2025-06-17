package organization

import "github.com/lnk.by/shared/service"

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (o *Organization) FieldsPtrs() []any {
	return []any{&o.ID, &o.Name}
}

func (o *Organization) FieldsVals() []any {
	return []any{o.ID, o.Name}
}

var (
	CreateSQL   service.CreateSQL[*Organization]   = "INSERT INTO organization (id, name) VALUES ($1, $2)"
	RetrieveSQL service.RetrieveSQL[*Organization] = "SELECT id, name FROM organization WHERE id = $1"
	UpdateSQL   service.UpdateSQL[*Organization]   = "UPDATE organization SET name = $2 WHERE id = $1"
	DeleteSQL   service.DeleteSQL[*Organization]   = "DELETE FROM organization WHERE id = $1"
	ListSQL     service.ListSQL[*Organization]     = "SELECT id, name FROM organization OFFSET $1 LIMIT $2"
)

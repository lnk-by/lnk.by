package organization

import (
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/utils"
)

type Organization struct {
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Status utils.Status `json:"status"`
}

func (o *Organization) FieldsPtrs() []any {
	return []any{&o.ID, &o.Name, &o.Status}
}

func (o *Organization) FieldsVals() []any {
	return []any{o.ID, o.Name, o.Status}
}

func (o *Organization) WithId(id string) {
	o.ID = id
}

func (o *Organization) Validate() error {
	switch {
	case o.Name == "":
		return service.ErrNameRequired
	case o.ID != "":
		return service.ErrIDManagedByServer
	default:
		return nil
	}
}

func (o *Organization) Generate() {
	o.ID = service.UUID()

	if o.Status == "" {
		o.Status = utils.StatusActive
	}
}

const IdParam = "organizationId"

var (
	CreateSQL   service.CreateSQL[*Organization]   = "INSERT INTO organization (id, name, status) VALUES ($1, $2, $3)"
	RetrieveSQL service.RetrieveSQL[*Organization] = "SELECT id, name, status FROM organization WHERE id = $1 AND status='active'"
	UpdateSQL   service.UpdateSQL[*Organization]   = "UPDATE organization SET name = $2, status=$3 WHERE id = $1"
	DeleteSQL   service.DeleteSQL[*Organization]   = "DELETE FROM organization WHERE id = $1"
	ListSQL     service.ListSQL[*Organization]     = "SELECT id, name, status FROM organization WHERE status='active' OFFSET $1 LIMIT $2"
)

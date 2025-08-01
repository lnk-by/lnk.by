package landingpage

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/lnk.by/aws/s3client"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/utils"
)

type LandingPage struct {
	ID             uuid.UUID              `json:"id"`
	Name           string                 `json:"name"`
	Template       string                 `json:"template"`      // the HTML file name
	Style          string                 `json:"style"`         // the CSS file name
	Configuration  map[string]interface{} `json:"configuration"` // stored as a JSON file
	OrganizationID *uuid.UUID             `json:"organizationId"`
	CustomerID     *uuid.UUID             `json:"customerId"`
	Status         utils.Status           `json:"status"`
}

func (p *LandingPage) FieldsPtrs() []any {
	return []any{&p.ID, &p.Name, &p.Template, &p.Style, &p.OrganizationID, &p.CustomerID, &p.Status}
}

func (p *LandingPage) FieldsVals() []any {
	return []any{p.ID, p.Name, p.Template, p.Style, p.OrganizationID, p.CustomerID, p.Status}
}

func (p *LandingPage) ParseID(idString string) (uuid.UUID, error) {
	return uuid.FromString(idString)
}

func (p *LandingPage) WithID(id uuid.UUID) {
	p.ID = id
}

func (p *LandingPage) Validate() error {
	switch {
	case p.Name == "":
		return service.ErrNameRequired
	case p.ID != uuid.Nil:
		return service.ErrIDManagedByServer
	case p.Template == "":
		return errors.New("email is required")
	default:
		return nil
	}
}

func (p *LandingPage) Generate() {
	p.ID = service.UUID()

	if p.Status == "" {
		p.Status = utils.StatusActive
	}
}

var (
	CreateSQL   service.CreateSQL[*LandingPage]   = "INSERT INTO landingpage (id, name, template, style, organization_id, customer_id, status) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	RetrieveSQL service.RetrieveSQL[*LandingPage] = "SELECT id, name, template, style, organization_id, customer_id, status FROM landingpage WHERE id = $1 AND status='active'"
	UpdateSQL   service.UpdateSQL[*LandingPage]   = "UPDATE landingpage SET name = $2, template=$3, style=$4, customer_id=$5, organization_id = $6, status = $7 WHERE id = $1"
	DeleteSQL   service.DeleteSQL[*LandingPage]   = "DELETE FROM landingpage WHERE id = $1"
	// Right now select all landing pages that belong to the same organization together with the currently logged in customer.
	ListSQL service.ListSQL[*LandingPage] = `
		SELECT p.id, p.name, p.template, p.style, p.organization_id, p.customer_id, p.status
		FROM landingpage p
		JOIN customer me ON me.id = $1
		WHERE p.status = 'active'
		AND (p.customer_id = me.id OR p.organization_id = me.organization_id)
		AND now() BETWEEN valid_from AND valid_until
		OFFSET $2 LIMIT $3
	`
)

func SetConfiguration(ctx context.Context, page *LandingPage) (*LandingPage, error) {
	if err := ReadConfiguration(ctx, getConfigurationPath(page.ID), &page.Configuration); err != nil {
		return page, err
	}
	return page, nil
}

func StroreConfiguration(ctx context.Context, page *LandingPage) error {
	return WriteConfiguration(ctx, getConfigurationPath(page.ID), page.Configuration)
}

func DeleteConfiguration(ctx context.Context, id uuid.UUID) error {
	return s3client.Delete(ctx, getConfigurationPath(id))
}

func getConfigurationPath(id uuid.UUID) string {
	return fmt.Sprintf("landingpages/conf/%s.json", id)
}

func CreateLandingPage(ctx context.Context, body string, userID *uuid.UUID) (int, string) {
	page, err := service.Parse[*LandingPage](ctx, []byte(body))
	if err != nil {
		return http.StatusBadRequest, "faied to parse LandingPage data"
	}
	if page.CustomerID == nil {
		page.CustomerID = userID
	}
	status, body := service.CreateRecord(ctx, CreateSQL, page, 0)
	if status >= http.StatusMultipleChoices {
		return http.StatusInternalServerError, body
	}

	err = WriteConfiguration(ctx, getConfigurationPath(page.ID), page.Configuration)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusCreated, body
}

package customer

import (
	"context"
	"os"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/lnk.by/shared/test/db"
	"github.com/lnk.by/shared/test/service"
	"github.com/lnk.by/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(
		func() int {
			stop := db.Start(context.Background())
			defer stop()

			return m.Run() // Run tests
		}(),
	)
}

func TestListEmpty(t *testing.T) {
	db.WithTable(t, "customer", func() {
		customers := service.List(t, ListSQL, 0, 10)
		assert.Equal(t, 0, len(customers))
	})
}

func TestCRUDL(t *testing.T) {
	email := "adam@human.net"
	email2 := "adam@robot.net"
	name := "Adam"
	adam := Customer{Email: email, Name: name}

	db.WithTable(t, "customer", func() {
		created := service.Create(t, CreateSQL, &adam)
		assert.Equal(t, email, created.Email)
		assert.Equal(t, name, created.Name)
		assert.Nil(t, created.OrganizationID)
		assert.Equal(t, utils.StatusActive, created.Status)

		retrieved := service.Retrieve(t, RetrieveSQL, created.ID.String())
		assert.Equal(t, created, retrieved)

		id := retrieved.ID
		retrieved.ID = uuid.Nil
		retrieved.Email = email2

		updated := service.Update(t, UpdateSQL, id.String(), retrieved)
		assert.Equal(t, id, updated.ID)
		assert.Equal(t, email2, updated.Email)
		assert.Equal(t, name, updated.Name)
		assert.Nil(t, updated.OrganizationID)
		assert.Equal(t, utils.StatusActive, updated.Status)

		listed := service.List(t, ListSQL, 0, 10)
		assert.Len(t, listed, 1)
		assert.Equal(t, updated, listed[0])

		service.Delete(t, DeleteSQL, id.String())

		listed = service.List(t, ListSQL, 0, 10)
		assert.Len(t, listed, 0)
	})
}

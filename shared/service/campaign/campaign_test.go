package campaign

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
	db.WithTable(t, "campaign", func() {
		campaigns := service.List(t, ListSQL, 0, 10)
		assert.Equal(t, 0, len(campaigns))
	})
}

func TestCRUDL(t *testing.T) {
	name := "NewYear"
	name2 := "Hanukah"
	adam := Campaign{Name: name}

	db.WithTable(t, "campaign", func() {
		created := service.Create(t, CreateSQL, &adam)
		assert.Equal(t, name, created.Name)
		assert.Nil(t, created.OrganizationID)
		assert.Nil(t, created.CustomerID)
		assert.Equal(t, utils.StatusActive, created.Status)

		retrieved := service.Retrieve(t, RetrieveSQL, created.ID.String())
		assert.Equal(t, created, retrieved)

		id := retrieved.ID
		retrieved.ID = uuid.Nil
		retrieved.Name = name2

		updated := service.Update(t, UpdateSQL, id.String(), retrieved)
		assert.Equal(t, id, updated.ID)
		assert.Equal(t, name2, updated.Name)
		assert.Nil(t, updated.OrganizationID)
		assert.Nil(t, updated.CustomerID)
		assert.Equal(t, utils.StatusActive, updated.Status)

		listed := service.List(t, ListSQL, 0, 10)
		assert.Len(t, listed, 1)
		assert.Equal(t, updated, listed[0])

		service.Delete(t, DeleteSQL, id.String())

		listed = service.List(t, ListSQL, 0, 10)
		assert.Len(t, listed, 0)
	})
}

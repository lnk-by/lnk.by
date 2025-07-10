package customer

import (
	"context"
	"os"
	"testing"

	"github.com/lnk.by/shared/test/db"
	"github.com/lnk.by/shared/test/service"
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

func TestCreateAndGet(t *testing.T) {
	db.WithTable(t, "customer", func() {
		adam := Customer{Email: "adam@human.net", Name: "Adam"}
		created := service.Create(t, CreateSQL, &adam)

		retrieved := service.Retrieve(t, RetrieveSQL, created.ID.String())
		assert.Equal(t, created, retrieved)

		listed := service.List(t, ListSQL, 0, 10)
		assert.Equal(t, 1, len(listed))
		assert.Equal(t, created, listed[0])
	})
}

// main.go - placeholder for cmd/server/main.go

package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lnk.by/shared/db"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/campaign"
	"github.com/lnk.by/shared/service/customer"
	"github.com/lnk.by/shared/service/organization"
	"github.com/lnk.by/shared/service/short_url"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

func init() {
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			slog.Error("Failed to load .env", "error", err)
			os.Exit(1)
		}
		slog.Info(".env file not found, continuing...")
	}

	if err := db.Init(context.Background(), os.Getenv("DB_URL"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
}

func list[T service.FieldsPtrsAware](c *gin.Context, sql service.ListSQL[T]) {
	offset, err := parseQueryInt(c, "offset", 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offest"})
		return
	}
	limit, err := parseQueryInt(c, "limit", math.MaxInt32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	status, body := service.List(c.Request.Context(), sql, offset, limit)
	respondWithJSON(c, status, body)
}

func parseQueryInt(c *gin.Context, key string, defaultValue int) (int, error) {
	valStr := c.Query(key)
	if valStr == "" {
		return defaultValue, nil
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("faild to parse value %q: must be an integer", key)
	}

	return val, nil
}

func listHandler[T service.FieldsPtrsAware](sql service.ListSQL[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		list(c, sql)
	}
}

func retrieveHandler[T service.FieldsPtrsAware](sql service.RetrieveSQL[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieve(c, sql)
	}
}

func retrieve[T service.FieldsPtrsAware](c *gin.Context, sql service.RetrieveSQL[T]) {
	id := c.Param("id")
	status, body := service.Retrieve(c.Request.Context(), sql, id)
	respondWithJSON(c, status, body)
}

func createCustomer(c *gin.Context) {
	var u customer.Customer
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	status, body := service.Create(c.Request.Context(), customer.CreateSQL, &u)
	respondWithJSON(c, status, body)
}

func updateCustomer(c *gin.Context) {
	var u customer.Customer
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	status, body := service.Update(c.Request.Context(), customer.UpdateSQL, c.Param("id"), &u)
	respondWithJSON(c, status, body)
}

////////////////////////////////

func createOrganization(c *gin.Context) {
	var o organization.Organization
	if err := c.BindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	status, body := service.Create(c.Request.Context(), organization.CreateSQL, &o)
	respondWithJSON(c, status, body)
}

func updateOrganization(c *gin.Context) {
	var o organization.Organization
	if err := c.BindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	status, body := service.Update(c.Request.Context(), organization.UpdateSQL, c.Param("id"), &o)
	respondWithJSON(c, status, body)
}

func createCampaign(c *gin.Context) {
	var e campaign.Campaign
	if err := c.BindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	status, body := service.Create(c.Request.Context(), campaign.CreateSQL, &e)
	respondWithJSON(c, status, body)
}

func updateCampaign(c *gin.Context) {
	var e campaign.Campaign
	if err := c.BindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	status, body := service.Update(c.Request.Context(), campaign.UpdateSQL, c.Param("id"), &e)
	respondWithJSON(c, status, body)
}

func createShortURL(c *gin.Context) {
	var e short_url.ShortURL
	if err := c.BindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	status, body := service.Create(c.Request.Context(), short_url.CreateSQL, &e)
	respondWithJSON(c, status, body)
}

func updateShortURL(c *gin.Context) {
	var e short_url.ShortURL
	if err := c.BindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	status, body := service.Update(c.Request.Context(), short_url.UpdateSQL, c.Param("id"), &e)
	respondWithJSON(c, status, body)
}

///////////////////////////////

func deleteHandler[T service.FieldsPtrsAware](sql service.DeleteSQL[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		deleteEntity(c, sql)
	}
}

func deleteEntity[T service.FieldsPtrsAware](c *gin.Context, sql service.DeleteSQL[T]) {
	id := c.Param("id")
	status, body := service.Delete(c.Request.Context(), sql, id)
	respondWithJSON(c, status, body)
}

func respondWithJSON(c *gin.Context, statusCode int, jsonStr string) {
	c.Header(contentTypeHeader, contentTypeJSON)
	c.String(statusCode, jsonStr)
}

func jsonErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors[0].Err
		c.JSON(-1, gin.H{"error": err.Error()})
		return
	}

	if c.Writer.Status() >= 400 {
		c.JSON(c.Writer.Status(), gin.H{"error": http.StatusText(c.Writer.Status())})
	}
}

func main() {
	router := gin.Default()
	router.Use(gin.Recovery(), jsonErrorHandler)
	router.RemoveExtraSlash = true

	router.POST("/customers", createCustomer)
	router.PUT("/customers/:id", updateCustomer)
	router.GET("/customers", listHandler(customer.ListSQL))
	router.GET("/customers/:id", retrieveHandler(customer.RetrieveSQL))
	router.DELETE("/customers/:id", deleteHandler(customer.DeleteSQL))

	router.POST("/organizations", createOrganization)
	router.PUT("/organizations/:id", updateOrganization)
	router.GET("/organizations", listHandler(organization.ListSQL))
	router.GET("/organizations/:id", retrieveHandler(organization.RetrieveSQL))
	router.DELETE("/organizations/:id", deleteHandler(organization.DeleteSQL))

	router.POST("/campaigns", createCampaign)
	router.PUT("/campaigns/:id", updateCampaign)
	router.GET("/campaigns", listHandler(campaign.ListSQL))
	router.GET("/campaigns/:id", retrieveHandler(campaign.RetrieveSQL))
	router.DELETE("/campaigns/:id", deleteHandler(campaign.DeleteSQL))

	router.POST("/shorturls", createShortURL)
	router.PUT("/shorturls/:id", updateShortURL)
	router.GET("/shorturls", listHandler(short_url.ListSQL))
	router.GET("/shorturls/:id", retrieveHandler(short_url.RetrieveSQL))
	router.DELETE("/shorturls/:id", deleteHandler(short_url.DeleteSQL))

	if err := router.Run("localhost:8080"); err != nil {
		slog.Error("Failed to start server", "error", err.Error())
	}
}

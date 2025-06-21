// main.go - placeholder for cmd/server/main.go

package main

import (
	"context"
	"fmt"
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
	"github.com/lnk.by/shared/service/shortURL"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

func init() {
	if err := godotenv.Load(); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(".env file not found, continuing...")
		} else {
			fmt.Fprintf(os.Stderr, "Failed to load .env: %v\n", err)
			os.Exit(1)
		}
	}

	if err := db.Init(context.Background(), os.Getenv("DB_URL"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func getEntities[T service.FieldsPtrsAware](c *gin.Context, sql service.ListSQL[T]) {
	offset, err := parseQueryInt(c, "offset", 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	limit, err := parseQueryInt(c, "limit", math.MaxInt32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
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

func getEntitiesHandler[T service.FieldsPtrsAware](sql service.ListSQL[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		getEntities(c, sql)
	}
}

func getEntityHandler[T service.FieldsPtrsAware](sql service.RetrieveSQL[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		getEntity(c, sql)
	}
}

func getEntity[T service.FieldsPtrsAware](c *gin.Context, sql service.RetrieveSQL[T]) {
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
	var e shortURL.ShortURL
	if err := c.BindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	status, body := service.Create(c.Request.Context(), shortURL.CreateSQL, &e)
	respondWithJSON(c, status, body)
}

func updateShortURL(c *gin.Context) {
	var e shortURL.ShortURL
	if err := c.BindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	status, body := service.Update(c.Request.Context(), shortURL.UpdateSQL, c.Param("id"), &e)
	respondWithJSON(c, status, body)
}

///////////////////////////////

func deleteEntityHandler[T service.FieldsPtrsAware](sql service.DeleteSQL[T]) gin.HandlerFunc {
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
	router.GET("/customers", getEntitiesHandler(customer.ListSQL))
	router.GET("/customers/:id", getEntityHandler(customer.RetrieveSQL))
	router.DELETE("/customers/:id", deleteEntityHandler(customer.DeleteSQL))

	router.POST("/organizations", createOrganization)
	router.PUT("/organizations/:id", updateOrganization)
	router.GET("/organizations", getEntitiesHandler(organization.ListSQL))
	router.GET("/organizations/:id", getEntityHandler(organization.RetrieveSQL))
	router.DELETE("/organizations/:id", deleteEntityHandler(organization.DeleteSQL))

	router.POST("/campaigns", createCampaign)
	router.PUT("/campaigns/:id", updateCampaign)
	router.GET("/campaigns", getEntitiesHandler(campaign.ListSQL))
	router.GET("/campaigns/:id", getEntityHandler(campaign.RetrieveSQL))
	router.DELETE("/campaigns/:id", deleteEntityHandler(campaign.DeleteSQL))

	router.POST("/shorturls", createShortURL)
	router.PUT("/shorturls/:id", updateShortURL)
	router.GET("/shorturls", getEntitiesHandler(shortURL.ListSQL))
	router.GET("/shorturls/:id", getEntityHandler(shortURL.RetrieveSQL))
	router.DELETE("/shorturls/:id", deleteEntityHandler(shortURL.DeleteSQL))

	router.Run("localhost:8080")
	fmt.Println("Server started")
}

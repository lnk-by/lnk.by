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

func retrieve[T service.FieldsPtrsAware](c *gin.Context, sql service.RetrieveSQL[T]) {
	status, body := service.Retrieve(c.Request.Context(), sql, c.Param("id"))
	respondWithJSON(c, status, body)
}

func create[T service.Creatable](c *gin.Context, sql service.CreateSQL[T]) {
	status, body := service.CreateFromReqBody(c.Request.Context(), sql, c.Request.Body)
	respondWithJSON(c, status, body)
}

func update[T service.Updatable](c *gin.Context, sql service.UpdateSQL[T]) {
	status, body := service.UpdateFromReqBody(c.Request.Context(), sql, c.Param("id"), c.Request.Body)
	respondWithJSON(c, status, body)
}

func deleteEntity[T service.FieldsPtrsAware](c *gin.Context, sql service.DeleteSQL[T]) {
	status, body := service.Delete(c.Request.Context(), sql, c.Param("id"))
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

	if c.Writer.Status() >= http.StatusBadRequest {
		c.JSON(c.Writer.Status(), gin.H{"error": http.StatusText(c.Writer.Status())})
	}
}

func redirect(c *gin.Context) {
	status, url, errStr := service.RetrieveValueAndMarshalError(c.Request.Context(), short_url.RetrieveSQL, c.Param("id"))
	if errStr != "" {
		//c.JSON(status, gin.H{"error": errStr})
		respondWithJSON(c, status, errStr)
		return
	}

	c.Redirect(http.StatusFound, url.Target)
}

func main() {
	router := gin.Default()
	router.Use(gin.Recovery(), jsonErrorHandler)
	router.RemoveExtraSlash = true

	router.POST("/customers", func(c *gin.Context) { create(c, customer.CreateSQL) })
	router.PUT("/customers/:id", func(c *gin.Context) { update(c, customer.UpdateSQL) })
	router.GET("/customers", func(c *gin.Context) { list(c, customer.ListSQL) })
	router.GET("/customers/:id", func(c *gin.Context) { retrieve(c, customer.RetrieveSQL) })
	router.DELETE("/customers/:id", func(c *gin.Context) { deleteEntity(c, customer.DeleteSQL) })

	router.POST("/organizations", func(c *gin.Context) { create(c, organization.CreateSQL) })
	router.PUT("/organizations/:id", func(c *gin.Context) { update(c, organization.UpdateSQL) })
	router.GET("/organizations", func(c *gin.Context) { list(c, organization.ListSQL) })
	router.GET("/organizations/:id", func(c *gin.Context) { retrieve(c, organization.RetrieveSQL) })
	router.DELETE("/organizations/:id", func(c *gin.Context) { deleteEntity(c, organization.DeleteSQL) })

	router.POST("/campaigns", func(c *gin.Context) { create(c, campaign.CreateSQL) })
	router.PUT("/campaigns/:id", func(c *gin.Context) { update(c, campaign.UpdateSQL) })
	router.GET("/campaigns", func(c *gin.Context) { list(c, campaign.ListSQL) })
	router.GET("/campaigns/:id", func(c *gin.Context) { retrieve(c, campaign.RetrieveSQL) })
	router.DELETE("/campaigns/:id", func(c *gin.Context) { deleteEntity(c, campaign.DeleteSQL) })

	router.POST("/shorturls", func(c *gin.Context) { create(c, short_url.CreateSQL) })
	router.PUT("/shorturls/:id", func(c *gin.Context) { update(c, short_url.UpdateSQL) })
	router.GET("/shorturls", func(c *gin.Context) { list(c, short_url.ListSQL) })
	router.GET("/shorturls/:id", func(c *gin.Context) { retrieve(c, short_url.RetrieveSQL) })
	router.DELETE("/shorturls/:id", func(c *gin.Context) { deleteEntity(c, short_url.DeleteSQL) })

	router.GET("/go/:id", redirect)

	if err := router.Run("localhost:8080"); err != nil {
		slog.Error("Failed to start server", "error", err.Error())
	}
}

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
	"github.com/lnk.by/shared/service/shorturl"
)

const (
	contentTypeHeader   = "Content-Type"
	contentTypeJSON     = "application/json"
	authorizationHeader = "Authorization"

	accessControlAllowOrigin   = "Access-Control-Allow-Origin"
	accessControlAllowMethods  = "Access-Control-Allow-Methods"
	accessControlAllowHeaders  = "Access-Control-Allow-Headers"
	accessControlExposeHeaders = "Access-Control-Expose-Headers"

	any = "*"
)

func initDbConnection() error {
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			slog.Error("Failed to load .env", "error", err)
			return err
		}
		slog.Info(".env file not found, continuing...")
	}

	if err := db.InitFromEnvironement(context.Background()); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return err
	}
	return nil
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
	c.Header(accessControlAllowOrigin, any)
	c.String(statusCode, jsonStr)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(accessControlAllowOrigin, any)
		c.Header(accessControlAllowMethods, "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header(accessControlAllowHeaders, "Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
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
	status, url, errStr := service.RetrieveValueAndMarshalError(c.Request.Context(), shorturl.RetrieveSQL, c.Param("id"))
	if errStr != "" {
		respondWithJSON(c, status, errStr)
		return
	}

	c.Redirect(http.StatusFound, url.Target)
}

func main() {
	router := gin.Default()
	router.Use(gin.Recovery(), jsonErrorHandler, CORSMiddleware())
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

	router.POST("/shorturls", func(c *gin.Context) { create(c, shorturl.CreateSQL) })
	router.PUT("/shorturls/:id", func(c *gin.Context) { update(c, shorturl.UpdateSQL) })
	router.GET("/shorturls", func(c *gin.Context) { list(c, shorturl.ListSQL) })
	router.GET("/shorturls/:id", func(c *gin.Context) { retrieve(c, shorturl.RetrieveSQL) })
	router.DELETE("/shorturls/:id", func(c *gin.Context) { deleteEntity(c, shorturl.DeleteSQL) })

	router.GET("/go/:id", redirect)

	if err := initDbConnection(); err != nil {
		slog.Error("Failed to start server", "error", err.Error())
		os.Exit(1)
	}

	if err := router.Run("localhost:8080"); err != nil {
		slog.Error("Failed to start server", "error", err.Error())
		os.Exit(1)
	}
}

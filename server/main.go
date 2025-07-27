package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/joho/godotenv"
	"github.com/lnk.by/shared/db"
	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/campaign"
	"github.com/lnk.by/shared/service/customer"
	"github.com/lnk.by/shared/service/landingpage"
	"github.com/lnk.by/shared/service/organization"
	"github.com/lnk.by/shared/service/shorturl"
	"github.com/lnk.by/shared/service/stats"
	"github.com/lnk.by/shared/service/stats/maxmind"
)

const (
	accessControlAllowOriginHeader  = "Access-Control-Allow-Origin"
	accessControlAllowMethodsHeader = "Access-Control-Allow-Methods"
	accessControlAllowHeadersHeader = "Access-Control-Allow-Headers"
	authorizationHeader             = "Authorization"
	contentTypeHeader               = "Content-Type"
	retryAfterHeader                = "Retry-After"
)

const contentTypeJSON = "application/json"
const allowAnyOrigin = "*"

func initDbConnection() error {
	if err := godotenv.Load(); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("failed to load .env: %w", err)
		}
		slog.Info(".env file not found, continuing...")
	}

	ctx := context.Background()

	if err := db.InitFromEnvironment(ctx); err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}

	initScript := os.Getenv("DB_INIT_SCRIPT")
	if initScript != "" {
		if err := db.RunScript(ctx, initScript); err != nil {
			return fmt.Errorf("failed to run SQL script: %w", err)
		}
	}

	return nil
}

func list[K any, T service.Retrievable[K]](c *gin.Context, sql service.ListSQL[T]) {
	listAndTransform(c, sql, func(t T) (T, error) { return t, nil })
}

func listAndTransform[K any, T service.Retrievable[K]](c *gin.Context, sql service.ListSQL[T], transformer func(t T) (T, error)) {
	userID := service.GetUUIDFromAuthorization(c.GetHeader("Authorization"))
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
	status, body := service.List(c.Request.Context(), sql, userID, offset, limit, transformer)
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

func retrieve[K any, T service.Retrievable[K]](c *gin.Context, sql service.RetrieveSQL[T]) {
	retrieveAndTransform(c, sql, func(t T) (T, error) { return t, nil })
}

func retrieveAndTransform[K any, T service.Retrievable[K]](c *gin.Context, sql service.RetrieveSQL[T], transformer func(t T) (T, error)) {
	status, body := service.Retrieve(c.Request.Context(), sql, c.Param("id"), transformer)
	respondWithJSON(c, status, body)
}

func create[T service.Creatable](c *gin.Context, sql service.CreateSQL[T]) {
	status, body := service.CreateFromReqBody(c.Request.Context(), sql, c.Request.Body)
	respondWithJSON(c, status, body)
}

func update[K any, T service.Updatable[K]](c *gin.Context, sql service.UpdateSQL[T]) {
	updateAndFinalize(c, sql, func(id K, t T) error { return nil })
}

func updateAndFinalize[K any, T service.Updatable[K]](c *gin.Context, sql service.UpdateSQL[T], finalizer func(id K, t T) error) {
	status, body := service.UpdateFromReqBody(c.Request.Context(), sql, c.Param("id"), c.Request.Body, finalizer)
	respondWithJSON(c, status, body)
}
func deleteEntity[K any, T service.Retrievable[K]](c *gin.Context, sql service.DeleteSQL[T]) {
	deleteEntityAndFinalize(c, sql, func(id K) error { return nil })
}

func deleteEntityAndFinalize[K any, T service.Retrievable[K]](c *gin.Context, sql service.DeleteSQL[T], finalizer func(id K) error) {
	status, body := service.Delete(c.Request.Context(), sql, c.Param("id"), finalizer)
	respondWithJSON(c, status, body)
}

func respondWithJSON(c *gin.Context, statusCode int, jsonStr string) {
	c.Header(contentTypeHeader, contentTypeJSON)
	c.Header(accessControlAllowOriginHeader, allowAnyOrigin)
	c.String(statusCode, jsonStr)
}

func createLandingPage(c *gin.Context) {
	content, err := io.ReadAll(c.Request.Body)
	if err != nil {
		respondWithJSON(c, http.StatusInternalServerError, fmt.Sprintf("{\"error\": %s}", fmt.Errorf("failed to read request body: %w", err)))
		return
	}
	status, body := landingpage.CreateLandingPage(c.Request.Context(), string(content), nil)
	respondWithJSON(c, status, body)

}

func createShortURL(c *gin.Context) {
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		respondWithJSON(c, http.StatusInternalServerError, fmt.Sprintf("{\"error\": %s}", fmt.Errorf("failed to read request body: %w", err)))
		return
	}
	status, responseBody := shorturl.CreateShortURL(c.Request.Context(), requestBody, nil)
	respondWithJSON(c, status, responseBody)
}

var (
	allowedMethods = strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}, ",")
	allowedHeaders = strings.Join([]string{authorizationHeader, contentTypeHeader}, ",")
)

func corsMiddleware(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
		return
	}

	c.Header(accessControlAllowOriginHeader, allowAnyOrigin)
	c.Header(accessControlAllowMethodsHeader, allowedMethods)
	c.Header(accessControlAllowHeadersHeader, allowedHeaders)
	c.AbortWithStatus(http.StatusNoContent)
}

func jsonErrorHandler(c *gin.Context) {
	c.Next()

	if c.Writer.Written() {
		return
	}

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
	key := c.Param("id")
	now := time.Now()
	day := fmt.Sprintf("day%03d", now.YearDay())
	hour := fmt.Sprintf("hour%02d", now.Hour())
	status, url, errStr := service.RetrieveValueAndMarshalError(c.Request.Context(), shorturl.RetrieveValidSQL, key, day, hour)
	if errStr != "" {
		respondWithJSON(c, status, errStr)
		return
	}

	if limitExceeded, retryAfter := shorturl.GetLimitExceededMessage(url); limitExceeded != "" {
		c.Header(retryAfterHeader, string(retryAfter))
		c.JSON(http.StatusTooManyRequests, gin.H{"error": limitExceeded})
		c.Writer.Written()
		return
	}

	if err := sendStatistics(c, key); err != nil {
		slog.Warn("Failed to send stats", "error", err)
	}

	c.Redirect(http.StatusFound, url.Target)
}

func sendStatistics(c *gin.Context, key string) error {
	header := c.Request.Header
	event := stats.Event{
		Key:       key,
		IP:        c.ClientIP(),
		UserAgent: header.Get("user-agent"),
		Referer:   header.Get("referer"),
		Timestamp: time.Now().UTC(),
		Language:  header.Get("accept-language"),
	}
	return stats.Process(c.Request.Context(), event)
}

func run() error {
	router := gin.Default()
	router.Use(gin.Recovery(), jsonErrorHandler, corsMiddleware)
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

	router.POST("/shorturls", func(c *gin.Context) { createShortURL(c) })
	router.PUT("/shorturls/:id", func(c *gin.Context) { update(c, shorturl.UpdateSQL) })
	router.GET("/shorturls", func(c *gin.Context) { list(c, shorturl.ListSQL) })
	router.GET("/shorturls/:id", func(c *gin.Context) { retrieve(c, shorturl.RetrieveSQL) })
	router.DELETE("/shorturls/:id", func(c *gin.Context) { deleteEntity(c, shorturl.DeleteSQL) })

	router.POST("/landingpages", func(c *gin.Context) { createLandingPage(c) })
	router.PUT("/landingpages/:id", func(c *gin.Context) {
		updateAndFinalize(c, landingpage.UpdateSQL, func(id uuid.UUID, p *landingpage.LandingPage) error {
			_, err := landingpage.SetConfiguration(c.Request.Context(), p)
			return err
		})
	})
	router.GET("/landingpages", func(c *gin.Context) {
		listAndTransform(c, landingpage.ListSQL,
			func(p *landingpage.LandingPage) (*landingpage.LandingPage, error) {
				return landingpage.SetConfiguration(c.Request.Context(), p)
			})
	})
	router.GET("/landingpages/:id", func(c *gin.Context) {
		retrieveAndTransform(c, landingpage.RetrieveSQL, func(p *landingpage.LandingPage) (*landingpage.LandingPage, error) {
			return landingpage.SetConfiguration(c.Request.Context(), p)
		})
	})
	router.DELETE("/landingpages/:id", func(c *gin.Context) {
		deleteEntityAndFinalize(c, landingpage.DeleteSQL, func(id uuid.UUID) error { return landingpage.DeleteConfiguration(c.Request.Context(), id) })
	})

	router.GET("/go/:id", redirect)

	router.Static("/ui", "../ui")
	router.Static("/landingpages", "../landingpages")

	if err := initDbConnection(); err != nil {
		return fmt.Errorf("failed to init DB connnection: %w", err)
	}
	if err := maxmind.Init(); err != nil {
		slog.Error("Failed to intialize mixmind", "error", err)
	}

	if err := router.Run(":8080"); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to run server", "error", err.Error())
		os.Exit(1)
	}
}

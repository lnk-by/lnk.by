package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lnk.by/aws/adapter"
	"github.com/lnk.by/aws/s3client"
	"github.com/lnk.by/shared/service/landingpage"
)

func listLandingPageTemplate(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	templatePaths, err := s3client.List(ctx, landingpage.Path, "html")
	if err != nil {
		slog.Error("Failed to retrieve list of landing page templates", "error", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: toJSON(err), Headers: adapter.StandardHeaders}, nil
	}

	var templates []landingpage.LandingPageTemplate
	for _, templatePath := range templatePaths {
		template, err := landingpage.RetrieveLandingPageTemplate(ctx, fileNameFromPath(templatePath))
		if err != nil {
			slog.Error("Failed to retrieve list of landing page template", "path", templatePath, "error", err)
			return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: toJSON(err), Headers: adapter.StandardHeaders}, nil
		}
		templates = append(templates, template)
	}
	bytes, err := json.Marshal(templates)
	if err != nil {
		slog.Error("Failed to marshal landing page templates", "error", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: toJSON(err), Headers: adapter.StandardHeaders}, nil
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusOK, Body: string(bytes), Headers: adapter.StandardHeaders}, nil
}

func fileNameFromPath(filepath string) string {
	filename := path.Base(filepath)
	return strings.TrimSuffix(filename, path.Ext(filename))
}

func toJSON(err error) string {
	body, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	return string(body)
}

func main() {
	s3client.InitializeFromEnvironment(context.Background())
	adapter.LambdaMain(listLandingPageTemplate)
}

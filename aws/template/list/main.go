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
	slog.Info("listLandingPageTemplate 1")
	templatePaths, err := s3client.List(ctx, landingpage.Path, "html")
	slog.Info("listLandingPageTemplate 2")
	if err != nil {
		slog.Info("listLandingPageTemplate 2.1", "error", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: toJSON(err), Headers: adapter.StandardHeaders}, nil
	}

	slog.Info("listLandingPageTemplate 3")
	var templates []landingpage.LandingPageTemplate
	for _, templatePath := range templatePaths {
		slog.Info("listLandingPageTemplate 4", "templatePath", templatePath)
		template, err := landingpage.RetrieveLandingPageTemplate(ctx, fileNameFromPath(templatePath))
		slog.Info("listLandingPageTemplate 5")
		if err != nil {
			slog.Info("listLandingPageTemplate 5.1", "error", err)
			return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: toJSON(err), Headers: adapter.StandardHeaders}, nil
		}
		templates = append(templates, template)
	}
	slog.Info("listLandingPageTemplate 6")
	bytes, err := json.Marshal(templates)
	slog.Info("listLandingPageTemplate 7")
	if err != nil {
		slog.Info("listLandingPageTemplate 7.1", "error", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError, Body: toJSON(err), Headers: adapter.StandardHeaders}, nil
	}

	slog.Info("listLandingPageTemplate 8")
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

package landingpage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lnk.by/aws/s3client"
)

type LandingPageTemplate struct {
	Name          string                 `json:"name"`
	Style         []string               `json:"style"`
	Configuration map[string]interface{} `json:"configuration"`
}

const Path = "landingpages/templates/"

func RetrieveLandingPageTemplate(ctx context.Context, name string) (LandingPageTemplate, error) {
	var template LandingPageTemplate

	template.Name = name

	conf, err := s3client.GetString(ctx, getFilePath(name, "json"))
	if err != nil {
		return template, fmt.Errorf("loading template %q config: %w", name, err)
	}

	if err := json.Unmarshal([]byte(conf), &template.Configuration); err != nil {
		return template, fmt.Errorf("parsing template %q config: %w", name, err)
	}

	styles, err := s3client.List(ctx, Path, "css")
	if err != nil {
		return template, fmt.Errorf("listing styles for template %q: %w", name, err)
	}
	template.Style = styles
	return template, nil
}

func getFilePath(name string, extension string) string {
	return fmt.Sprintf("%s%s.%s", Path, name, extension)
}

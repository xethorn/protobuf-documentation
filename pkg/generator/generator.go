package generator

import (
	"bytes"
	"encoding/json"
	"html/template"

	"github.com/user/proto-doc/pkg/parser"
)

type Config struct {
	Title             string `json:"title"`
	PrimaryColor      string `json:"primary_color"`
	PrimaryHoverColor string `json:"primary_hover_color"`
	FontFamily        string `json:"font_family"`
	Theme             string `json:"theme"`
	LogoURL           string `json:"logo_url"`
}

var DefaultConfig = Config{
	Title:             "API Documentation",
	PrimaryColor:      "#006aff",
	PrimaryHoverColor: "#0051c7",
	FontFamily:        "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
	Theme:             "light",
	LogoURL:           "",
}

type RenderData struct {
	Config      Config
	DocData     *parser.RootResult
	DocDataJSON template.JS
}

type HTMLGenerator struct {
	config Config
}

func NewHTMLGenerator(cfg Config) *HTMLGenerator {
	merged := DefaultConfig
	if cfg.Title != "" {
		merged.Title = cfg.Title
	}
	if cfg.PrimaryColor != "" {
		merged.PrimaryColor = cfg.PrimaryColor
	}
	if cfg.PrimaryHoverColor != "" {
		merged.PrimaryHoverColor = cfg.PrimaryHoverColor
	}
	if cfg.FontFamily != "" {
		merged.FontFamily = cfg.FontFamily
	}
	if cfg.Theme != "" {
		merged.Theme = cfg.Theme
	}
	if cfg.LogoURL != "" {
		merged.LogoURL = cfg.LogoURL
	}
	return &HTMLGenerator{config: merged}
}

func (g *HTMLGenerator) Render(docData *parser.RootResult) (string, error) {
	tmpl, err := template.New("doc").Parse(HTMLTemplate)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(docData)
	if err != nil {
		return "", err
	}

	data := RenderData{
		Config:      g.config,
		DocData:     docData,
		DocDataJSON: template.JS(string(jsonBytes)),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

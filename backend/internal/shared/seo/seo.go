package seo

import (
	"bytes"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// SEOConfig berisi konfigurasi yang di-inject dari AppConfig saat startup.
type SEOConfig struct {
	BaseURL            string
	SiteName           string
	DefaultTitle       string
	DefaultDescription string
	DefaultImage       string
}

// MetaTagData berisi data yang akan di-render ke template HTML Open Graph.
type MetaTagData struct {
	Title       string
	Description string
	ImageURL    string
	PageURL     string
	SiteName    string
}

var seoTemplate = template.Must(template.New("seo").Parse(`<!DOCTYPE html>
<html lang="id">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.Title}}</title>
	<meta name="description" content="{{.Description}}">
	<meta property="og:type" content="article">
	<meta property="og:title" content="{{.Title}}">
	<meta property="og:description" content="{{.Description}}">
	<meta property="og:image" content="{{.ImageURL}}">
	<meta property="og:url" content="{{.PageURL}}">
	<meta property="og:site_name" content="{{.SiteName}}">
	<meta name="twitter:card" content="summary_large_image">
	<meta name="twitter:title" content="{{.Title}}">
	<meta name="twitter:description" content="{{.Description}}">
	<meta name="twitter:image" content="{{.ImageURL}}">
</head>
<body>
	<h1>{{.Title}}</h1>
	<p>{{.Description}}</p>
</body>
</html>`))

func RenderMetaHTML(c *gin.Context, data MetaTagData) {
	var buf bytes.Buffer
	if err := seoTemplate.Execute(&buf, data); err != nil {
		c.Data(http.StatusInternalServerError, "text/html; charset=utf-8",
			[]byte("<!DOCTYPE html><html><head><title>Error</title></head><body><h1>Internal Server Error</h1></body></html>"))
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
}

func DefaultMetaTagData(cfg SEOConfig, pathSegment string) MetaTagData {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	return MetaTagData{
		Title:       cfg.DefaultTitle,
		Description: cfg.DefaultDescription,
		ImageURL:    ResolveAbsoluteImageURL(baseURL, cfg.DefaultImage),
		PageURL:     baseURL + "/" + strings.TrimLeft(pathSegment, "/"),
		SiteName:    cfg.SiteName,
	}
}

func ResolveAbsoluteImageURL(baseURL, imagePath string) string {
	if imagePath == "" {
		return ""
	}
	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		return imagePath
	}
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(imagePath, "/")
}

var reHTMLTag = regexp.MustCompile(`<[^>]*>`)

func TruncateDescription(text string, maxLen int) string {
	cleaned := reHTMLTag.ReplaceAllString(text, " ")
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	if utf8.RuneCountInString(cleaned) <= maxLen {
		return cleaned
	}
	runes := []rune(cleaned)
	truncated := string(runes[:maxLen])
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > maxLen/2 {
		truncated = truncated[:lastSpace]
	}
	return truncated + "..."
}

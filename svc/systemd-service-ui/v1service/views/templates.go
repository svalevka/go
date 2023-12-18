package views

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

// embed template files within the binary at build time.
//
//go:embed templates
var html embed.FS

// pre-compile HTML templates at startup using the embedded filesystem.
var tpl = template.Must(template.New("").ParseFS(html, "templates/*.html"))

// Render executes a pre-compiled HTML template sourced from the given View,
// which special handling for overwriting the HTTP 200 Status Code.
func Render(w http.ResponseWriter, r *http.Request, view View) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	status := http.StatusOK
	if sc, ok := view.(interface {
		StatusCode() int
	}); ok {
		status = sc.StatusCode()
	}

	w.WriteHeader(status)

	err := tpl.ExecuteTemplate(w, view.TemplateName(), view)
	if err != nil {
		return fmt.Errorf("template: %w", err)
	}

	return nil
}

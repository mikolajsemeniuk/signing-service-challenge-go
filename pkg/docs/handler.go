package docs

import (
	"html/template"
	"net/http"
)

func NewHandler() *Handler {
	router := http.NewServeMux()

	handler := &Handler{router: router}
	handler.router.HandleFunc("GET /", handler.Elements)
	handler.router.HandleFunc("GET /docs", handler.OpenAPI)

	return handler
}

// Handler provides API compatible with HTTP and REST standards.
type Handler struct {
	router *http.ServeMux
}

// ServeHTTP is used for joining handlers to HTTP server.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// Elements serves elements ui.
func (h *Handler) Elements(w http.ResponseWriter, r *http.Request) {
	template.Must(template.New("ui").Parse(elements)).Execute(w, "./docs")
}

// Elements serves specification in OpenAPI standard.
func (h *Handler) OpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Write(docs)
}

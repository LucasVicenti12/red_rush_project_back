package module

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Module interface {
	Name() string
	Setup(r *mux.Router) *mux.Router
}

type Router struct {
	Url     string
	Methods []string
	Handler http.HandlerFunc
}

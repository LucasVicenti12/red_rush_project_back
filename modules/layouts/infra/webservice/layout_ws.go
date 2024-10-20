package webservice

import (
	"AwTV/config/module"
	"AwTV/modules/layouts/domain/use_case"
	"AwTV/modules/layouts/infra/repository"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type LayoutWebService struct {
	name    string
	base    string
	usecase use_case.LayoutUseCase
}

func NewDeviceWebService(conn *sql.DB) module.Module {
	rep := repository.NewLayoutRepositoryImpl(conn)

	return &LayoutWebService{
		name: "layout_web_service",
		base: "layout",
		usecase: &use_case.LayoutUseCaseImpl{
			Repository: rep,
		},
	}
}

func (lws *LayoutWebService) Name() string {
	return lws.name
}

func (lws LayoutWebService) Setup(r *mux.Router) *mux.Router {
	handlers := []module.Router{
		{
			Url:     "/getByName",
			Methods: []string{http.MethodGet},
			Handler: lws.getLayoutByName,
		},
		{
			Url:     "/getByUUID",
			Methods: []string{http.MethodGet},
			Handler: lws.getLayoutByUUID,
		},
		{
			Url:     "/getLayouts",
			Methods: []string{http.MethodGet},
			Handler: lws.getLayouts,
		},
		{
			Url:     "/saveLayout",
			Methods: []string{http.MethodPost, http.MethodPut},
			Handler: lws.saveLayout,
		},
		{
			Url:     "/deleteLayout",
			Methods: []string{http.MethodDelete},
			Handler: lws.deleteLayout,
		},
	}

	for _, h := range handlers {
		r.HandleFunc("/"+lws.base+h.Url, h.Handler).Methods(h.Methods...)
	}

	return r
}

func (lws *LayoutWebService) getLayouts(w http.ResponseWriter, r *http.Request) {
	layouts, err := lws.usecase.GetLayouts()

	if err != nil {
		_, _ = w.Write([]byte(err.Error()))

		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, _ := json.Marshal(layouts)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (lws *LayoutWebService) getLayoutByUUID(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")

	layout, err := lws.usecase.GetLayoutByUUID(uuid)

	if err != nil {
		_, _ = w.Write([]byte(err.Error()))

		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, _ := json.Marshal(layout)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (lws *LayoutWebService) getLayoutByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	layout, err := lws.usecase.GetLayoutByName(name)

	if err != nil {
		_, _ = w.Write([]byte(err.Error()))

		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, _ := json.Marshal(layout)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (lws *LayoutWebService) saveLayout(w http.ResponseWriter, r *http.Request) {

}

func (lws *LayoutWebService) deleteLayout(w http.ResponseWriter, r *http.Request) {

}

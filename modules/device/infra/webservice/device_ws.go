package webservice

import (
	"AwTV/config/module"
	"AwTV/modules/device/domain/entities"
	"AwTV/modules/device/domain/use_case"
	"AwTV/modules/device/infra/repository"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type DeviceWebService struct {
	name    string
	base    string
	usecase use_case.DeviceUseCase
}

func NewDeviceWebService(conn *sql.DB) module.Module {
	rep := repository.NewDeviceRepository(conn)

	return &DeviceWebService{
		name: "device_web_service",
		base: "device",
		usecase: &use_case.DeviceUseCaseImpl{
			Repository: rep,
		},
	}
}

func (dws DeviceWebService) Name() string {
	return dws.name
}

func (dws DeviceWebService) Setup(r *mux.Router) *mux.Router {
	handlers := []module.Router{
		{
			Url:     "/getByName",
			Methods: []string{http.MethodGet},
			Handler: dws.getDeviceByName,
		},
		{
			Url:     "/getByUUID",
			Methods: []string{http.MethodGet},
			Handler: dws.getDeviceByUUID,
		},
		{
			Url:     "/getDevices",
			Methods: []string{http.MethodGet},
			Handler: dws.getDevices,
		},
		{
			Url:     "/saveDevice",
			Methods: []string{http.MethodPost, http.MethodPut},
			Handler: dws.saveDevice,
		},
		{
			Url:     "/deleteDevice",
			Methods: []string{http.MethodDelete},
			Handler: dws.deleteDevice,
		},
	}

	for _, h := range handlers {
		r.HandleFunc("/"+dws.base+h.Url, h.Handler).Methods(h.Methods...)
	}

	return r
}

func (dws DeviceWebService) getDeviceByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]

	device, err := dws.usecase.GetDeviceByName(name)

	if err != nil {
		_, _ = w.Write([]byte(err.Error()))

		w.WriteHeader(http.StatusNoContent)
	}

	response, _ := json.Marshal(device)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (dws DeviceWebService) getDeviceByUUID(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")

	device, err := dws.usecase.GetDeviceByUUID(uuid)

	if err != nil {
		_, _ = w.Write([]byte(err.Error()))

		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, _ := json.Marshal(device)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (dws DeviceWebService) getDevices(w http.ResponseWriter, _ *http.Request) {
	device, err := dws.usecase.GetAllDevices()

	if err != nil {
		_, _ = w.Write([]byte(err.Error()))

		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, _ := json.Marshal(device)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (dws DeviceWebService) saveDevice(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		_, _ = w.Write([]byte("empty body"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var device entities.Device

	err = json.Unmarshal(body, &device)

	if err != nil {
		_, _ = w.Write([]byte("invalid json"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = dws.usecase.SaveDevice(&device)

	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	response, _ := json.Marshal(device)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(response)
}

func (dws DeviceWebService) deleteDevice(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")

	err := dws.usecase.DeleteDevice(uuid)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("device removed successfully"))
}

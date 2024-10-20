package webservice

import (
	"AwTV/config/module"
	"AwTV/config/user/domain/entities"
	"AwTV/config/user/domain/use_case"
	"AwTV/config/user/infra/repository"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type UserWebService struct {
	name    string
	base    string
	usecase use_case.UserUseCase
}

func NewUserWebService(conn *sql.DB) module.Module {
	rep := repository.NewUserRepository(conn)

	return &UserWebService{
		name: "user_web_service",
		base: "user",
		usecase: &use_case.UserUseCaseImpl{
			Repository: rep,
		},
	}
}

func (us UserWebService) Name() string {
	return us.name
}

func (us UserWebService) Setup(r *mux.Router) *mux.Router {
	handlers := []module.Router{
		{
			Url:     "/getByUUID/{uuid}",
			Methods: []string{http.MethodGet},
			Handler: us.getUserByUUID,
		},
		{
			Url:     "/getByNickname/{nickname}",
			Methods: []string{http.MethodGet},
			Handler: us.getByNickname,
		},
		{
			Url:     "/register",
			Methods: []string{http.MethodPost},
			Handler: us.registerUser,
		},
		{
			Url:     "/list",
			Methods: []string{http.MethodGet},
			Handler: us.list,
		},
	}

	for _, h := range handlers {
		r.HandleFunc("/"+us.base+h.Url, h.Handler).Methods(h.Methods...)
	}

	return r
}

func (us UserWebService) getUserByUUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]

	user, err := us.usecase.GetUserByUUID(uuid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	response, err := json.Marshal(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (us UserWebService) getByNickname(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	nickname := vars["nickname"]

	user, err := us.usecase.GetUserByNickname(nickname)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	response, err := json.Marshal(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (us UserWebService) registerUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	var user entities.User

	user.UserType = entities.Admin

	err = json.Unmarshal(body, &user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	err = us.usecase.RegisterUser(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("User registered successfully"))
}

func (us UserWebService) list(w http.ResponseWriter, r *http.Request) {
	users, err := us.usecase.ListUsers()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

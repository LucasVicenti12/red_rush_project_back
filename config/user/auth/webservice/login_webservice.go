package webservice

import (
	"AwTV/config/module"
	entities2 "AwTV/config/user/auth/entities"
	"AwTV/config/user/domain/use_case"
	"AwTV/config/user/infra/repository"
	"database/sql"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"time"
)

type authWebService struct {
	name    string
	base    string
	usecase use_case.UserUseCase
}

func NewAuthWebService(conn *sql.DB) module.Module {
	rep := repository.NewUserRepository(conn)

	return &authWebService{
		name: "auth_web_service",
		base: "auth",
		usecase: &use_case.UserUseCaseImpl{
			Repository: rep,
		},
	}
}

func (us authWebService) Name() string {
	return us.name
}

func (us authWebService) Setup(r *mux.Router) *mux.Router {
	handlers := []module.Router{
		{
			Url:     "/login",
			Methods: []string{http.MethodPost},
			Handler: us.login,
		},
	}

	for _, h := range handlers {
		r.HandleFunc("/"+us.base+h.Url, h.Handler).Methods(h.Methods...)
	}

	return r
}

func (us authWebService) login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	var loginReq entities2.Login

	err = json.Unmarshal(body, &loginReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	user, err := us.usecase.GetUserByNickname(loginReq.Nickname)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Invalid username or password"))
		if err != nil {
			return
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Invalid password"))
		if err != nil {
			return
		}
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Uuid,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(entities2.SecretKey))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			return
		}
		return
	}

	cookie := http.Cookie{
		Name:     entities2.CookieName,
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Secure:   true,
		Domain:   "localhost",
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Login successful"))
}

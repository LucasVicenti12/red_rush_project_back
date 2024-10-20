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
		{
			Url:     "/logout",
			Methods: []string{http.MethodPost},
			Handler: us.logout,
		},
		{
			Url:     "/user",
			Methods: []string{http.MethodGet},
			Handler: us.user,
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
	_, _ = w.Write([]byte("Login successful"))
}

func (us authWebService) logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     entities2.CookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
		Domain:   "localhost",
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/app/login", http.StatusFound)
}

func (us authWebService) user(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(entities2.CookieName)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Invalid cookie"))
		return
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(entities2.SecretKey), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Invalid token"))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Invalid token"))
		return
	}

	userUUID, _ := claims["sub"].(string)

	user, err := us.usecase.GetUserByUUID(userUUID)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("No user found"))
		return
	}

	response, err := json.Marshal(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Erro on marshalling response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

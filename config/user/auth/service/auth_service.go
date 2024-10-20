package service

import (
	"AwTV/config/user/auth/entities"
	"AwTV/shared/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func AuthorizationMiddleware(next http.Handler) http.Handler {
	requests := []string{
		"/auth/login",
		"/user/register",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		temp, err := mux.CurrentRoute(r).GetPathTemplate()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		allow := utils.Contains(requests, temp)

		if allow {
			next.ServeHTTP(w, r)
		} else {
			cookie, err := r.Cookie(entities.CookieName)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Invalid cookie"))
				return
			}

			if cookie.Value == "" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("Cookie is empty"))
				return
			}

			next.ServeHTTP(w, r)
		}
	})
}

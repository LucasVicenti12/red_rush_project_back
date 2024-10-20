package main

import (
	"AwTV/config/database"
	"AwTV/config/user/auth/service"
	webservice3 "AwTV/config/user/auth/webservice"
	webservice2 "AwTV/config/user/infra/webservice"
	"AwTV/modules/device/infra/webservice"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {
	cors := handlers.AllowedOriginValidator(func(o string) bool {
		return true
	})

	dbconn, err := database.ConnectDatabase()

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.Use(service.AuthorizationMiddleware)

	authWS := webservice3.NewAuthWebService(dbconn)
	authWS.Setup(r)

	userWS := webservice2.NewUserWebService(dbconn)
	userWS.Setup(r)

	deviceWS := webservice.NewDeviceWebService(dbconn)
	deviceWS.Setup(r)

	server := &http.Server{
		Addr: ":9000",
		Handler: handlers.CORS(
			cors,
			handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Accept"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodOptions,
			}),
			handlers.AllowCredentials(),
		)(r),
		ReadTimeout:       time.Second * 240,
		ReadHeaderTimeout: time.Second * 20,
		WriteTimeout:      time.Second * 420,
		IdleTimeout:       time.Second * 15,
	}

	_ = server.ListenAndServe()
}

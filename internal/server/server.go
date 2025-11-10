package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"booking.com/internal/config"
	"booking.com/internal/handlers/auth"
	"booking.com/internal/handlers/users"
	"booking.com/internal/svcs"
	jwtauth "booking.com/pkg/auth/jwt-auth"
	"booking.com/pkg/constants"
	"booking.com/pkg/utils"
)

func StartHttpTlsServer(cfg *config.AppConfig) error {
	tlsConfig, err := utils.CreateTlsConfig(cfg.HttpServer.CertPath, cfg.HttpServer.KeyPath, "")
	if err != nil {
		return err
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/health", commonMiddleWare(health, false, http.MethodGet))

	registerAuthApp(mux, cfg)
	registerUsersApp(mux, cfg)

	tlsConfig.ClientAuth = tls.NoClientCert

	httpTlsServer := &http.Server{
		Addr:      cfg.HttpServer.Address,
		TLSConfig: tlsConfig,
		Handler:   mux,
	}

	log.Printf("Server started on https://%v", httpTlsServer.Addr)

	if err := httpTlsServer.ListenAndServeTLS("", ""); err != nil {
		return err
	}
	return nil
}
func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set(constants.ContentType, constants.ContentTypeTextPlain)
	w.Write([]byte("Sever is Up and Running"))
}

func registerUsersApp(mux *http.ServeMux, cfg *config.AppConfig) {
	usrHandler := users.NewUserHandler(&svcs.UserSvc{AppCfg: cfg})

	mux.HandleFunc("/users/register", commonMiddleWare(usrHandler.Add, true, http.MethodPost))
	mux.HandleFunc("/users/update", commonMiddleWare(usrHandler.Put, true, http.MethodPost))
}
func registerAuthApp(mux *http.ServeMux, cfg *config.AppConfig) {
	authHandler := auth.NewAuthHandler(&svcs.AuthSvc{AppCfg: cfg})

	mux.HandleFunc("/login", commonMiddleWare(authHandler.Login, false, http.MethodPost))
	mux.HandleFunc("/refresh", commonMiddleWare(authHandler.Refresh, false, http.MethodPost))
}
func commonMiddleWare(handler http.HandlerFunc, authNeeded bool, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var found bool
		if slices.Contains(methods, r.Method) {
			found = true
		}
		if !found {
			http.Error(w, fmt.Sprintf("Method %s Not Allowed, Expected Methods %v", r.Method, strings.Join(methods, ",")), http.StatusMethodNotAllowed)
			return
		}
		if authNeeded {
			authHeader := r.Header.Get(constants.Authorization)
			if authHeader == "" || !strings.HasPrefix(authHeader, constants.Bearer) {
				http.Error(w, "token not provided", http.StatusUnauthorized)
				return
			}
			token := authHeader[len(constants.Bearer):]
			fmt.Println("==========", token)
			if token == "" {
				http.Error(w, "token not provided", http.StatusUnauthorized)
				return
			}
			_, _, err := jwtauth.IsTokenValid(token, false)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed validate token, error: %v", err), http.StatusUnauthorized)
				return
			}
		}
		handler(w, r)
	}
}

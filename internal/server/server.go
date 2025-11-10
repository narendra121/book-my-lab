package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

	"booking.com/internal/apis/auth"
	"booking.com/internal/apis/users"
	"booking.com/internal/config"
	jwtauth "booking.com/pkg/auth/jwt-auth"
	"booking.com/pkg/constants"
	"booking.com/pkg/utils"
)

func StartHttpTlsServer(cfg config.Server) error {
	tlsConfig, err := utils.CreateTlsConfig(cfg.CertPath, cfg.KeyPath, "")
	if err != nil {
		return err
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/health", commonMiddleWare(health, false, http.MethodGet))

	registerAuthApp(mux)
	registerUsersApp(mux)

	tlsConfig.ClientAuth = tls.NoClientCert

	httpTlsServer := &http.Server{
		Addr:      cfg.Address,
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

func registerUsersApp(mux *http.ServeMux) {
	mux.HandleFunc("/users/register", commonMiddleWare(users.Add, true, http.MethodPost))
	mux.HandleFunc("/users/update", commonMiddleWare(users.Put, true, http.MethodPost))
}
func registerAuthApp(mux *http.ServeMux) {
	mux.HandleFunc("/login", commonMiddleWare(auth.Login, false, http.MethodPost))
	mux.HandleFunc("/refresh", commonMiddleWare(auth.Refresh, false, http.MethodPost))
}
func commonMiddleWare(handler http.HandlerFunc, authNeede bool, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var found bool
		for _, method := range methods {
			if r.Method == method {
				found = true
				break
			}
		}
		if !found {
			http.Error(w, fmt.Sprintf("Method %s Not Allowed, Expected Methods %v", r.Method, strings.Join(methods, ",")), http.StatusMethodNotAllowed)
			return
		}
		token := r.Header.Get("Bearer")
		if token == "" {
			http.Error(w, "token not provided", http.StatusUnauthorized)
			return
		}
		if authNeede {
			_, _, err := jwtauth.IsTokenValid(token, true)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed validate token, error: %v", err), http.StatusUnauthorized)
				return
			}
		}
		handler(w, r)
	}
}

package server

// import (
// 	"crypto/tls"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"slices"
// 	"strings"

// 	"booking.com/internal/config"
// 	"booking.com/internal/db/postgresql/dao"
// 	"booking.com/internal/handlers/auth"
// 	"booking.com/internal/handlers/user"
// 	"booking.com/internal/svcs"
// 	"booking.com/internal/utils"
// 	jwtauth "booking.com/pkg/auth/jwt-auth"
// 	"booking.com/pkg/constants"
// )

// func StartHttpTlsServer(cfg *config.AppConfig) error {
// 	tlsConfig, err := utils.CreateTlsConfig(cfg.HttpServer.CertPath, cfg.HttpServer.KeyPath, "")
// 	if err != nil {
// 		return err
// 	}
// 	mux := http.NewServeMux()

// 	mux.HandleFunc("/health", commonMiddleWare(health, false, http.MethodGet))

// 	registerAuthApp(mux, cfg)
// 	registerUsersApp(mux, cfg)

// 	tlsConfig.ClientAuth = tls.NoClientCert

// 	httpTlsServer := &http.Server{
// 		Addr:      cfg.HttpServer.Address,
// 		TLSConfig: tlsConfig,
// 		Handler:   mux,
// 	}

// 	if err := httpTlsServer.ListenAndServeTLS("", ""); err != nil {
// 		return err
// 	}
// 	return nil
// }
// func health(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set(constants.ContentType, constants.ContentTypeTextPlain)
// 	w.Write([]byte("Sever is Up and Running"))
// }

// func registerUsersApp(mux *http.ServeMux, cfg *config.AppConfig) {
// 	usrHandler := user.NewUserHandler(&svcs.UserSvc{AppCfg: cfg})

// 	mux.HandleFunc("/user/register", commonMiddleWare(usrHandler.Add, false, http.MethodPost))
// 	mux.HandleFunc("/user/profile", commonMiddleWare(usrHandler.Get, true, http.MethodGet))
// 	mux.HandleFunc("/user/update", commonMiddleWare(usrHandler.Put, true, http.MethodPut))
// 	mux.HandleFunc("/user/list", commonMiddleWare(usrHandler.GetAll, true, http.MethodPut))
// 	mux.HandleFunc("/user/update-role", commonMiddleWare(usrHandler.UpdateRole, true, http.MethodPut))
// }
// func registerAuthApp(mux *http.ServeMux, cfg *config.AppConfig) {
// 	authHandler := auth.NewAuthHandler(&svcs.AuthSvc{AppCfg: cfg})

// 	mux.HandleFunc("/auth/login", commonMiddleWare(authHandler.Login, false, http.MethodPost))
// 	mux.HandleFunc("/auth/refresh", commonMiddleWare(authHandler.Refresh, false, http.MethodPost))
// 	// mux.HandleFunc("/auth/logout", commonMiddleWare(authHandler.Refresh, false, http.MethodPost))

// }
// func registerPropertiesApp(mux *http.ServeMux, cfg *config.AppConfig) {
// 	// mux.HandleFunc("/properties", commonMiddleWare(usrHandler.Put, true, http.MethodGet))
// 	// mux.HandleFunc("/properties", commonMiddleWare(usrHandler.Put, true, http.MethodGet))
// 	// mux.HandleFunc("/properties/{id}", commonMiddleWare(usrHandler.Put, true, http.MethodGet))
// 	// mux.HandleFunc("/properties/{id}", commonMiddleWare(usrHandler.Put, true, http.MethodGet))
// }
// func registerVisitsApp(mux *http.ServeMux, cfg *config.AppConfig) {
// 	// mux.HandleFunc("/visits", commonMiddleWare(usrHandler.Put, true, http.MethodGet))
// 	// mux.HandleFunc("/visits", commonMiddleWare(usrHandler.Put, true, http.MethodGet))
// 	// mux.HandleFunc("/visits/{id}", commonMiddleWare(usrHandler.Put, true, http.MethodGet))
// 	// mux.HandleFunc("/visits/{id}", commonMiddleWare(usrHandler.Put, true, http.MethodGet))
// }
// func commonMiddleWare(handler http.HandlerFunc, authNeeded bool, methods ...string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var found bool
// 		if slices.Contains(methods, r.Method) {
// 			found = true
// 		}
// 		if !found {
// 			http.Error(w, fmt.Sprintf("Method %s Not Allowed, Expected Methods %v", r.Method, strings.Join(methods, ",")), http.StatusMethodNotAllowed)
// 			return
// 		}
// 		handler(w, r)
// 	}
// }

// func Au

package router

import (
	"net/http"
	"os"

	"github.com/rysmaadit/go-template/handler"
	"github.com/rysmaadit/go-template/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(dependencies service.Dependencies) http.Handler {
	r := mux.NewRouter()

	setAuthRouter(r, dependencies.AuthService)
	setCheckRouter(r, dependencies.CheckService)
	setUserRouter(r, dependencies.UserService)
	setCustomerRouter(r, dependencies.CustomerService)
	setPetugasRouter(r, dependencies.PetugasService)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	return loggedRouter
}

func setAuthRouter(router *mux.Router, dependencies service.AuthServiceInterface) {
	router.Methods(http.MethodGet).Path("/auth/token").Handler(handler.GetToken(dependencies))
	router.Methods(http.MethodPost).Path("/auth/token/validate").Handler(handler.ValidateToken(dependencies))
}

func setCheckRouter(router *mux.Router, checkService service.CheckService) {
	router.Methods(http.MethodGet).Path("/check/mysql").Handler(handler.CheckMysql(checkService))
}

func setUserRouter(router *mux.Router, dependencies service.UserServiceInterface) {
	router.Methods(http.MethodPost).Path("/signup").Handler(handler.Create(dependencies))
	router.Methods(http.MethodPost).Path("/login").Handler(handler.Login(dependencies))
}

func setCustomerRouter(router *mux.Router, dependencies service.CustomerServiceInterface) {
	router.Methods(http.MethodGet).Path("/cekpengajuan").Handler(handler.GetCekPengajuan(dependencies))
	router.Methods(http.MethodPost).Path("/pengajuan").Handler(handler.CreatePengajuan(dependencies))
	router.Methods(http.MethodPost).Path("/createkelengkapan").Handler(handler.CreateKelengkapan(dependencies))
	router.Methods(http.MethodGet).Path("/kelengkapan").Handler(handler.GetByIdKelengkapan(dependencies))
	router.Methods(http.MethodGet).Path("/submission/getstatus").Handler(handler.GetStatusByIdKelengkapan(dependencies))
}

func setPetugasRouter(router *mux.Router, dependencies service.PetugasServiceInterface) {
	router.Methods(http.MethodGet).Path("/listpengajuan/{page}").Handler(handler.GetListPengajuan(dependencies))
	router.Methods(http.MethodGet).Path("/countpage").Handler(handler.GetCountPage(dependencies))
	router.Methods(http.MethodGet).Path("/searchbyname/{name}").Handler(handler.GetListByName(dependencies))
	router.Methods(http.MethodGet).Path("/kelengkapan/{id}").Handler(handler.GetSubmission(dependencies))
	router.Methods(http.MethodPost).Path("/kelengkapan/status/{id_cust}").Handler(handler.PostSubmissionStatus(dependencies))
	router.Methods(http.MethodPost).Path("/datadiri/status/{id_cust}").Handler(handler.PostIdentityStatus(dependencies))
}

// package router

// import (
// 	"net/http"
// 	"os"

// 	"github.com/rysmaadit/go-template/handler"
// 	"github.com/rysmaadit/go-template/service"

// 	"github.com/gorilla/handlers"
// 	"github.com/gorilla/mux"
// )

// func NewRouter(dependencies service.Dependencies) http.Handler {
// 	r := mux.NewRouter()

// 	setAuthRouter(r, dependencies.AuthService)
// 	setCheckRouter(r, dependencies.CheckService)

// 	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
// 	return loggedRouter
// }

// func setAuthRouter(router *mux.Router, dependencies service.AuthServiceInterface) {
// 	router.Methods(http.MethodGet).Path("/auth/token").Handler(handler.GetToken(dependencies))
// 	router.Methods(http.MethodPost).Path("/auth/token/validate").Handler(handler.ValidateToken(dependencies))
// }

// func setCheckRouter(router *mux.Router, checkService service.CheckService) {
// 	router.Methods(http.MethodGet).Path("/check/redis").Handler(handler.CheckRedis(checkService))
// 	router.Methods(http.MethodGet).Path("/check/mysql").Handler(handler.CheckMysql(checkService))
// 	router.Methods(http.MethodGet).Path("/check/minio").Handler(handler.CheckMinio(checkService))
// }

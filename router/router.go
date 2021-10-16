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
	setEmployeeRouter(r, dependencies.EmployeeService)

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
	router.Methods(http.MethodGet).Path("/checkapply").Handler(handler.GetCheckApply(dependencies))
	router.Methods(http.MethodPost).Path("/createidentity").Handler(handler.CreateIdentity(dependencies))
	router.Methods(http.MethodPost).Path("/createsubmission").Handler(handler.CreateSubmission(dependencies))
	router.Methods(http.MethodGet).Path("/submission/getstatus").Handler(handler.GetSubmissionStatus(dependencies))
	router.Methods(http.MethodGet).Path("/getsubmission").Handler(handler.GetSubmissionCustomer(dependencies))
}

func setEmployeeRouter(router *mux.Router, dependencies service.EmployeeServiceInterface) {
	router.Methods(http.MethodGet).Path("/numberofpage").Handler(handler.GetNumberOfPage(dependencies))
	router.Methods(http.MethodGet).Path("/listsubmission/{page}").Handler(handler.GetListSubmission(dependencies))
	router.Methods(http.MethodGet).Path("/searchbyname/{name}").Handler(handler.GetListByName(dependencies))
	router.Methods(http.MethodGet).Path("/submission/{id}").Handler(handler.GetSubmissionEmployee(dependencies))
	router.Methods(http.MethodPost).Path("submission/status/{id_cust}").Handler(handler.PostSubmissionStatus(dependencies))
	router.Methods(http.MethodPost).Path("/identity/status/{id_cust}").Handler(handler.PostIdentityStatus(dependencies))
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

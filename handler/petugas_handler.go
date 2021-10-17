package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/service"
	log "github.com/sirupsen/logrus"
)

func GetListSubmission(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		page := vars["page"]

		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		dataService, err := employeeService.SPGetListSubmission(page)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetNumberOfPage(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		dataService := employeeService.SPGetNumberOfPage()

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetListByName(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		dataService := employeeService.SPGetListByName(name)

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetSubmissionEmployee(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		custId := vars["id"]

		subIdint, err := strconv.Atoi(custId)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		dataService, err := employeeService.SPGetSubmission(uint(subIdint))

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func PostSubmissionStatus(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		subId := vars["id_cust"]

		subIdint, err := strconv.Atoi(subId)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		payloads, _ := ioutil.ReadAll(r.Body)
		var statusSubmission contract.Submission
		json.Unmarshal(payloads, &statusSubmission)

		dataService, err := employeeService.SPPostSubmissionStatus(&statusSubmission, uint(subIdint))

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetIdentityEmployee(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		subId := vars["id_cust"]

		subIdint, err := strconv.Atoi(subId)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		dataService, err := employeeService.SPGetIdentityEmployee(uint(subIdint))

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func PostIdentityStatus(petugasService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := petugasService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		custId := vars["id_cust"]

		custIdint, err := strconv.Atoi(custId)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		payloads, _ := ioutil.ReadAll(r.Body)
		var statusDataDIri contract.Identity
		json.Unmarshal(payloads, &statusDataDIri)

		dataService, err := petugasService.SPPostIdentityStatus(&statusDataDIri, uint(custIdint))

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetFileKtpEmployee(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		vars := mux.Vars(r)
		buktiKtp := vars["bukti_ktp"]

		dataService := employeeService.SPGetFileKtp(buktiKtp)

		data, readErr := ioutil.ReadAll(dataService)
		if readErr != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println("Can't read object ")
			return
		} else {
			w.Header().Set("Content-Type", "application/pdf")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}
		// responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func TotalIdentityUnconfirmed(employeeService service.EmployeeServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := employeeService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		if resp.LoginAs != 2 {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		dataService, err := employeeService.SPGetTotalIdentityUnconfirmed()

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

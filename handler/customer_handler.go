package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/service"
	log "github.com/sirupsen/logrus"
)

func GetCheckApply(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		dataService := customerService.SCGetCheckApply(resp.IdUser)

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func CreateIdentity(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}
		payloads, _ := ioutil.ReadAll(r.Body)

		var identity contract.Identity

		json.Unmarshal(payloads, &identity)

		validate := validator.New()
		error := validate.Struct(identity)

		if error != nil {
			log.Warning(error)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, error)
			return
		}

		dataService, err := customerService.SCCreateIdentity(&identity, resp.IdUser)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func CreateSubmission(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		payloads, _ := ioutil.ReadAll(r.Body)

		var submission contract.Submission
		json.Unmarshal(payloads, &submission)

		validate := validator.New()
		error := validate.Struct(submission)

		if error != nil {
			log.Warning(error)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, error)
			return
		}

		dataService := customerService.SCCreateSubmission(&submission, resp.IdUser)
		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
		return

	}
}

func GetSubmissionCustomer(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		dataService, err := customerService.SCGetSubmission(resp.IdUser)
		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetSubmissionStatus(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		dataService := customerService.SCGetSubmissionStatus(resp.IdUser)
		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func UploadFileKtp(customerService service.CustomerServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenC, err := contract.NewValidateTokenRequestViaCookie(r)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		resp, err := customerService.VerifyToken(tokenC)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("myFile")

		if err != nil {
			fmt.Println("error Retrieving file from form-data\n", err)
			return
		}
		defer file.Close()

		dataService := customerService.SCUploadFile(&file, handler, resp)
		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

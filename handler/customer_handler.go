package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/service"
	log "github.com/sirupsen/logrus"
)

func GetCekPengajuan(customerService service.CustomerServiceInterface) http.HandlerFunc {
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

		dataService := customerService.SCGetCekPengajuan(resp.IdUser)

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func CreatePengajuan(customerService service.CustomerServiceInterface) http.HandlerFunc {
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

		var pengajuan contract.Pengajuan

		json.Unmarshal(payloads, &pengajuan)

		validate := validator.New()
		error := validate.Struct(pengajuan)

		if error != nil {
			log.Warning(error)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, error)
			return
		}

		dataService, err := customerService.SCCreatePengajuan(&pengajuan, resp.IdUser)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func CreateKelengkapan(customerService service.CustomerServiceInterface) http.HandlerFunc {
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

		var kelengkapan contract.Kelengkapan
		json.Unmarshal(payloads, &kelengkapan)

		dataService := customerService.SCCreateKelengkapan(&kelengkapan, resp.IdUser)

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetByIdKelengkapan(customerService service.CustomerServiceInterface) http.HandlerFunc {
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

		dataService, err := customerService.SCGetByIdKelengkapan(resp.IdUser)
		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

		data_service := customerService.GetCekPengajuan(resp.Id_user)

		responder.NewHttpResponse(r, w, http.StatusOK, data_service, nil)
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

		data_service := customerService.CreatePengajuan(&pengajuan, resp.Id_user)

		responder.NewHttpResponse(r, w, http.StatusOK, data_service, nil)
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

		data_service := customerService.CreateKelengkapan(&kelengkapan, resp.Id_user)

		responder.NewHttpResponse(r, w, http.StatusOK, data_service, nil)
	}
}

func GetById_kelengkapan(customerService service.CustomerServiceInterface) http.HandlerFunc {
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

		data_service := customerService.GetById_kelengkapan(resp.Id_user)

		responder.NewHttpResponse(r, w, http.StatusOK, data_service, nil)
	}
}

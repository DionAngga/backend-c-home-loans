package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/service"
	log "github.com/sirupsen/logrus"
)

func GetListPengajuan(petugasService service.PetugasServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		page := vars["page"]

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

		dataService, err := petugasService.SPGetListPengajuan(page)

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetCountPage(petugasService service.PetugasServiceInterface) http.HandlerFunc {
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

		dataService := petugasService.SPGetCountPage()

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetListByName(petugasService service.PetugasServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
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

		dataService := petugasService.SPGetListByName(name)

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func GetSubmission(petugasService service.PetugasServiceInterface) http.HandlerFunc {
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
		subId := vars["id"]

		subIdint, err := strconv.Atoi(subId)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusBadRequest, nil, err)
			return
		}

		dataService, err := petugasService.SPGetSubmission(uint(subIdint))

		if err != nil {
			log.Error(err)
			responder.NewHttpResponse(r, w, http.StatusInternalServerError, nil, err)
			return
		}

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

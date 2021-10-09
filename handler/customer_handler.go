package handler

import (
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

package handler

import (
	"net/http"

	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/service"
)

func GetListPengajuan(pengajuanService service.PengajuanServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data_service := pengajuanService.GetListPengajuan(w, r)

		responder.NewHttpResponse(r, w, http.StatusOK, data_service, nil)
	}
}

func GetSebelumPengajuan(pengajuanService service.PengajuanServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data_service := pengajuanService.GetSebelumPengajuan(w, r)

		responder.NewHttpResponse(r, w, http.StatusOK, data_service, nil)
	}
}

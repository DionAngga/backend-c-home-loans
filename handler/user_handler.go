package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/service"
	log "github.com/sirupsen/logrus"
)

func Create(userService service.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data_service := userService.Create(w, r)

		responder.NewHttpResponse(r, w, http.StatusOK, data_service, nil)
	}
}

func Login(userService service.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payloads, _ := ioutil.ReadAll(r.Body)
		var user contract.User
		json.Unmarshal(payloads, &user)

		data_service, err := userService.Login(&user)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		expirationTime := time.Now().Add(time.Hour * 1)

		http.SetCookie(w,
			&http.Cookie{
				Name:    "token",
				Value:   data_service.Token,
				Expires: expirationTime,
			})

		responder.NewHttpResponse(r, w, http.StatusOK, data_service, nil)
	}
}

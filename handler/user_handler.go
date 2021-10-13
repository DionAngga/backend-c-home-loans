package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/rysmaadit/go-template/common/responder"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/service"
	log "github.com/sirupsen/logrus"
)

func Create(userService service.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payloads, _ := ioutil.ReadAll(r.Body)

		var user contract.User
		json.Unmarshal(payloads, &user)

		dataService := userService.SUCreate(&user)

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

func Login(userService service.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payloads, _ := ioutil.ReadAll(r.Body)
		var user contract.User
		json.Unmarshal(payloads, &user)

		validte := validator.New()
		err := validte.Struct(user)
		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		dataService, err := userService.SULogin(&user)

		if err != nil {
			log.Warning(err)
			responder.NewHttpResponse(r, w, http.StatusUnauthorized, nil, err)
			return
		}

		expirationTime := time.Now().Add(time.Hour * 1)

		http.SetCookie(w,
			&http.Cookie{
				Name:    "token",
				Value:   dataService.Token,
				Expires: expirationTime,
			})

		responder.NewHttpResponse(r, w, http.StatusOK, dataService, nil)
	}
}

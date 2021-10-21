package main

import (
	"net/http"
	"os"

	"github.com/rysmaadit/go-template/app"
	"github.com/rysmaadit/go-template/cli"
)

func main() {
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
	})
	c := cli.NewCli(os.Args)
	c.Run(app.Init())
}

package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var rend *render.Render = render.New()

type Handler struct {
	http.Handler
}

func CreateHandler() *Handler {
	mux := mux.NewRouter()
	handler := &Handler{
		Handler: mux,
	}

	mux.HandleFunc("/", handler.HomeHandler).Methods("GET")                              // HTML, CSS, JS 요청
	mux.HandleFunc("/ping", handler.PingHandler).Methods("GET")                          // Ping Check
	mux.HandleFunc("/get/model-list", handler.getModelListHandler).Methods("GET")        // Model List 반환
	mux.HandleFunc("/model/{name:[a-z-_]+}/infer", handler.inferHandler).Methods("POST") // Triton Server Inference 요청

	return handler
}

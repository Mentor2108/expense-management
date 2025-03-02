package rest

import (
	"backend-service/defn"
	"backend-service/rest/user"
	"backend-service/util"
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func AddRoutes(router *httprouter.Router) {
	log := util.GetGlobalLogger(context.Background())

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(time.Now().Format(time.RFC822Z), " 404 for - ", r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	})

	router.Handle(http.MethodGet, "/status", JwtProtectedRoutes(ServerStatus))

	user.AddRoutes(router)
}

func ServerStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set(defn.HTTPHeaderContentType, defn.ContentTypePlainText)
	w.Write([]byte("Server is up and running"))
}

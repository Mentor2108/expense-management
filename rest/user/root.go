package user

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AddRoutes(router *httprouter.Router) {
	router.Handle(http.MethodPost, "/auth/signup", Signup)
	router.Handle(http.MethodPost, "/auth/login", Login)
	// router.Handle(http.MethodPost, "/signup", Signup)
	// router.Handle(http.MethodGet, "/portfolio", RetrievePortfolio)
}

package rest

import (
	"backend-service/defn"
	"backend-service/util"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func ApplyMiddleware(router *httprouter.Router) http.Handler {
	return panicHandler(responseContentTypeJSON(router))
}

func responseContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(defn.HTTPHeaderContentType, defn.ContentTypeJSON)

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func panicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			log := util.GetGlobalLogger(r.Context())
			if panicVal := recover(); panicVal != nil {
				log.Printf("Recovered in middleware:\n%+v\n%s\n", panicVal, string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "something went wrong in server's side"}`))
			}
		}()

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func JwtProtectedRoutes(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if strings.HasPrefix(r.URL.Path, "/auth") {
			if next != nil {
				next(w, r, ps)
			}
			return
		}

		var tokenString string
		cookie, err := r.Cookie(defn.AuthorizationHeader)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				tokenString = r.Header.Get(defn.AuthorizationHeader)
				if strings.EqualFold(tokenString, "") {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"error": "missing token"}`))
					return
				}
				tokenString = strings.TrimPrefix(tokenString, defn.AuthorizationBearer)
			default:
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "something went wrong in server's side"}`))
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
		} else {
			tokenString = cookie.Value
		}

		ctxWithUser, cerr := util.VerifyTokenAndGetCurrentUserContext(r.Context(), tokenString)
		if cerr != nil {
			util.RespondWithError(r.Context(), w, http.StatusUnauthorized, cerr)
			return
		}

		if next != nil {
			next(w, r.WithContext(ctxWithUser), ps)
		}
	}
}

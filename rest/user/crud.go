package user

import (
	"backend-service/defn"
	"backend-service/service"
	"backend-service/util"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func addJWTCookieToResponse(w http.ResponseWriter, tokenString string) {
	cookie := http.Cookie{
		Name:     defn.AuthorizationHeader,
		Value:    tokenString,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   int(defn.JWTExpirationTime),
	}
	http.SetCookie(w, &cookie)
}

func Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	log := util.GetGlobalLogger(ctx)

	userProfileBody := defn.UserProfile{}
	if err := json.NewDecoder(r.Body).Decode(&userProfileBody); err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeFailedToParseRequestBody, defn.ErrFailedToParseRequestBody, map[string]string{
			"error": err.Error(),
		})
		log.Println("failed to parse request body", cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	//Checking Mandatory Fields
	if _, err := mail.ParseAddress(userProfileBody.Email); err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeInputInvalidFormat, defn.ErrInputInvalidFormat, map[string]string{
			"message": "invalid email provided",
		})
		log.Println("error while parsing email:", err)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	if len(userProfileBody.Password) < 8 || len(userProfileBody.Password) > defn.MaximumPasswordLength {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeInputInvalidFormat, defn.ErrInputInvalidFormat, map[string]string{
			"message": fmt.Sprintf("password should be between 8 and %d characters", defn.MaximumPasswordLength),
		})
		log.Println(cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	newUser, cerr := service.CreateUser(ctx, userProfileBody)
	if cerr != nil {
		if strings.EqualFold(cerr.Code, defn.ErrCodeUnexpectedError) {
			util.RespondWithError(ctx, w, http.StatusInternalServerError, cerr)
			return
		}
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	tokenString, cerr := util.CreateToken(ctx, newUser["id"].(string), defn.RoleUser)
	if cerr != nil {
		util.RespondWithError(ctx, w, http.StatusInternalServerError, cerr)
		return
	}

	addJWTCookieToResponse(w, tokenString)
	util.SendResponseMapWithStatus(ctx, w, http.StatusCreated, newUser)
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	log := util.GetGlobalLogger(ctx)

	loginRequestBody := defn.UserProfile{}
	if err := json.NewDecoder(r.Body).Decode(&loginRequestBody); err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeFailedToParseRequestBody, defn.ErrFailedToParseRequestBody, map[string]string{
			"error": err.Error(),
		})
		log.Println("failed to parse request body", cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	//Checking Mandatory Fields
	if strings.EqualFold(loginRequestBody.Email, "") {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeMissingRequiredField, defn.ErrMissingRequiredField, map[string]string{
			"field": "email",
		})
		log.Println(cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	if strings.EqualFold(loginRequestBody.Password, "") {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeMissingRequiredField, defn.ErrMissingRequiredField, map[string]string{
			"field": "password",
		})
		log.Println(cerr)
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	jwtResponse, cerr := service.Login(ctx, loginRequestBody.Email, loginRequestBody.Password)
	if cerr != nil {
		if strings.EqualFold(cerr.Code, defn.ErrCodeUnexpectedError) {
			util.RespondWithError(ctx, w, http.StatusInternalServerError, cerr)
			return
		}
		util.RespondWithError(ctx, w, http.StatusBadRequest, cerr)
		return
	}

	addJWTCookieToResponse(w, jwtResponse["token"].(string))
	util.SendResponseMapWithStatus(ctx, w, http.StatusOK, jwtResponse)
}

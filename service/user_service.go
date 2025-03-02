package service

import (
	"backend-service/data"
	"backend-service/defn"
	"backend-service/util"
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx context.Context, userData defn.UserProfile) (map[string]interface{}, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		cerr := util.NewCustomErrorWithKeys(ctx, defn.ErrCodeUnexpectedError, defn.ErrUnexpectedError, map[string]string{
			"error": err.Error(),
		})
		log.Println(cerr)
		return nil, cerr
	}

	userData.Password = string(hashedPassword)
	userRepo := data.NewUserRepository()

	newUser, cerr := userRepo.Create(ctx, userData)
	if cerr != nil {
		log.Println(cerr)
		return nil, cerr
	}
	return newUser, nil
}

func Login(ctx context.Context, email, password string) (map[string]interface{}, *util.CustomError) {
	log := util.GetGlobalLogger(ctx)

	userRepo := data.NewUserRepository()
	userData, cerr := userRepo.GetByEmail(ctx, email)
	if cerr != nil {
		if strings.EqualFold(cerr.Code, defn.ErrCodeNoDataFound) {
			cerr := util.NewCustomError(ctx, defn.ErrCodeLoginFailed, defn.ErrLoginFailed)
			log.Println(cerr)
			return nil, cerr
		}
		log.Println(cerr)
		return nil, cerr
	}

	err := bcrypt.CompareHashAndPassword([]byte(userData["password"].(string)), []byte(password))
	if err != nil {
		cerr := util.NewCustomError(ctx, defn.ErrCodeLoginFailed, defn.ErrLoginFailed)
		log.Println(cerr)
		return nil, cerr
	}

	stringToken, cerr := util.CreateToken(ctx, userData["id"].(string), defn.RoleUser)
	if cerr != nil {
		return nil, cerr
	}

	return map[string]interface{}{
		"token":      stringToken,
		"expires_in": defn.JWTExpirationTime / 1000000, //converting to seconds
	}, nil
}

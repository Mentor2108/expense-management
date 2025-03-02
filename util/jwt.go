package util

import (
	"backend-service/defn"
	"context"
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

// func init() {
// secretKey = generateSecretKey(context.Background(), 64)
// }

func generateSecretKey(ctx context.Context, length int) string {
	log := GetGlobalLogger(ctx)

	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatalf("Failed to generate secret key: %v", err)
	}
	return base64.StdEncoding.EncodeToString(key)
}

func CreateToken(ctx context.Context, id, role string) (string, *CustomError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": id,                //username
			"iss": "backend-service", //issuer
			"aud": role,              //role
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(defn.JWTExpirationTime).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		cerr := NewCustomErrorWithKeys(ctx, defn.ErrCodeUnexpectedError, defn.ErrUnexpectedError, map[string]string{
			"error": err.Error(),
		})
		GetGlobalLogger(ctx).Println(cerr)
		return "", cerr
	}
	return tokenString, nil
}

func VerifyTokenAndGetCurrentUserContext(ctx context.Context, tokenString string) (context.Context, *CustomError) {
	log := GetGlobalLogger(ctx)
	token, err := jwt.Parse((tokenString), func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		cerr := NewCustomError(ctx, defn.ErrCodeInvalidToken, defn.ErrInvalidToken)
		log.Println("an error occured while parsing token: ", err)
		return nil, cerr
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		cerr := NewCustomError(ctx, defn.ErrCodeInvalidToken, defn.ErrInvalidToken)
		log.Println(cerr)
		return nil, cerr
	} else {
		subject, err := claims.GetSubject()
		if err != nil {
			cerr := NewCustomError(ctx, defn.ErrCodeInvalidToken, defn.ErrInvalidToken)
			log.Println(cerr)
			return nil, cerr
		}
		return context.WithValue(ctx, defn.ContextUserKey, subject), nil
	}
}

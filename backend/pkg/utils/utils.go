package utils

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"tourism/pkg/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

var db *gorm.DB = config.GetDB()
var redisClient *redis.Client = config.GetRedis()
var jwtKey = []byte("mySecretKey")
type ContextKey string
const UserContextKey ContextKey = "user"

var Purposes = map[string]string{
	"login": "_1",
}

type CustomClaims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateJWT(id uuid.UUID) (string, error) {
    claims := &CustomClaims{
        UserID: id.String(),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed_token, err := token.SignedString(jwtKey)
    return signed_token, err
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
	
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		tokenStr := tokenParts[1]

		token, _ := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		claims, _ := token.Claims.(*CustomClaims)

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		var user gorm.Model
		if err := db.First(&user, "id = ?", userID).Error; err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		
		ctx := context.WithValue(r.Context(), UserContextKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func GenerateOTP(phoneNumber string, purpose string) string {
	otpInt := (rand.Intn(900000) + 100000)
	otp := fmt.Sprintf("%d", otpInt)

	ctx := context.Background()
	redisClient.Set(ctx, phoneNumber+ purpose, otp, 5*time.Minute)

	return otp
}

func ValidateOTP(phoneNumber, purpose, otp string) bool {
	ctx := context.Background()
	storedOTP, err := redisClient.Get(ctx, phoneNumber+purpose).Result()
	log.Print(storedOTP)
	if err != nil || otp != storedOTP {
		return false
	}
	return true
}
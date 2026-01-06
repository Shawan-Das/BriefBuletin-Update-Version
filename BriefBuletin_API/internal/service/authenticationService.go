package service

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rest/api/internal/model"
)

func (s *RESTService) getHashOf(password string) string {
	shaBytes := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", shaBytes)
}

func (s *RESTService) createJWTToken(userID int32, email, userName string, role string) string {
	if s.jwtSigningKey == nil {
		return ""
	}
	var expDate int64
	if role == "ADMIN" {
		expDate = time.Now().Add(1 * time.Hour).Unix()
	} else {
		expDate = time.Now().Add(24 * time.Hour).Unix()
	}
	claim := model.AuthorizationClaims{
		UserID:   userID,
		Email:    email,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expDate,
			Issuer:    "Auth Service",
			Id:        fmt.Sprintf("%d", userID),
		},
		Role: role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := token.SignedString(s.jwtSigningKey)
	if err != nil {
		_asLogger.Error("Error in generating token", err)
		return ""
	}
	_asLogger.Infof("Generated token for user %s", email)

	return tokenStr
}

func (s *RESTService) checkAuth(c *gin.Context) bool {
	url := c.Request.URL
	uri := url.RequestURI()

	// Allow bypass URLs
	if _, isFound := s.bypassAuth[uri]; isFound {
		return true
	}

	// If signing key is not set, allow all
	if s.jwtSigningKey == nil {
		return true
	}

	// Read Authorization header
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse + Validate
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSigningKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	// Build a struct (add more fields if needed)
	userInfo := model.UserRoleInfo{
		UserID:   int(claims["user_id"].(float64)),
		Email:    fmt.Sprintf("%v", claims["email"]),
		UserName: fmt.Sprintf("%v", claims["user_name"]),
		Role:     fmt.Sprintf("%v", claims["role"]),
		RoleMap: map[string]bool{
			fmt.Sprintf("%v", claims["role"]): true,
		},
	}

	// Store full claims in Gin context
	c.Set("__USER_INFO__", userInfo)

	return true
}

func GetLoggedInUser(c *gin.Context) (*model.UserRoleInfo, bool) {
	val, ok := c.Get("__USER_INFO__")
	if !ok {
		return nil, false
	}

	userInfo, ok := val.(model.UserRoleInfo)
	if !ok {
		return nil, false
	}

	return &userInfo, true
}

// user, ok := GetLoggedInUser(c)
// 	if !ok {
// 		return BuildResponse401("Unauthorized")
// 	}

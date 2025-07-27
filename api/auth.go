package api

import (
	"net/http"
	"time"

	db "github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Store       db.Querier
	JWTSecret   string
	TokenExpiry time.Duration
}

func NewAuthHandler(store db.Querier, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		Store:       store,
		JWTSecret:   jwtSecret,
		TokenExpiry: time.Hour * 24,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/auth")
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user, err := h.Store.CreateUser(c, db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hashed),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not register"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": user.ID,
		"email":   user.Email,
	})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Store.GetUserByEmail(c, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(h.TokenExpiry).Unix(),
	})

	signed, _ := token.SignedString([]byte(h.JWTSecret))

	c.JSON(http.StatusOK, gin.H{
		"access_token": signed,
	})
}

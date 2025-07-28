package api

import (
	"net/http"

	"github.com/Stylozone/go-ecom-api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	store sqlc.Querier
}

func NewUserHandler(store sqlc.Querier) *UserHandler {
	return &UserHandler{store: store}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc) {
	protected := r.Group("/user")
	protected.Use(auth)
	protected.GET("/me", h.Me)
}

func (h *UserHandler) Me(c *gin.Context) {
	// Get user_id from JWT middleware context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(int32)

	user, err := h.store.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"role":       user.Role.String,
		"created_at": user.CreatedAt,
	})
}

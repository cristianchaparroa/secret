package api

import (
	"fmt"
	"net/http"

	"github.com/cristianchaparroa/secret/core/response"
	"github.com/cristianchaparroa/secret/services"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SecretPostRequest contians the date related with the post request
type SecretPostRequest struct {

	// Secret is the plain text secret
	Secret string `binding:"required"`

	// The secret won't be available after the given number of views.
	// It must be greater than 0.
	ExpireAfterViews int32 `binding:"required"`
	// The secret won't be available after the given time.
	// The value is provided in minutes. 0 means never expires
	ExpireAfter int32 `binding:"required"`
}

// SecretHandler manges the request realted to secrets
type SecretHandler struct {
	db *gorm.DB
}

// NewSecretHandler returns a pointer to SecretHandler
func NewSecretHandler(db *gorm.DB) *SecretHandler {
	return &SecretHandler{db: db}
}

// CreateSecret add a new secret
func (h *SecretHandler) CreateSecret(c *gin.Context) {

	var req SecretPostRequest

	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusMethodNotAllowed, fmt.Sprintf("Invalid input:%v", err.Error()))
		return
	}

	ss := services.NewSecretService(h.db)
	m := ss.CreateSecret(req.Secret, req.ExpireAfterViews, req.ExpireAfter)

	builder := response.NewBuilder(c, http.StatusOK, m)
	builder.Render()
}

// FindSecret retrieve a secret according with a hash
func (h *SecretHandler) FindSecret(c *gin.Context) {
	hash := c.Param("hash")
	ss := services.NewSecretService(h.db)
	m := ss.GetSecret(hash)

	if m == nil {
		c.String(http.StatusNotFound, "Secret not found")
		return
	}
	builder := response.NewBuilder(c, http.StatusOK, m)
	builder.Render()
}

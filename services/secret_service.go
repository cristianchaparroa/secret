package services

import (
	"fmt"
	"time"

	"github.com/cristianchaparroa/secret/models"
	rep "github.com/cristianchaparroa/secret/repositories"
	"github.com/cristianchaparroa/secret/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// ISecretService defines the methods to handle the secrets
type ISecretService interface {

	// Create a new secret
	CreateSecret(secret string, remainingViews int32, expiresAt time.Time) *models.Secret

	// GetSecret retrieves a secret acording with hash if is still availabe.
	GetSecret(hash string) *models.Secret

	// IsSecretAvailble verifies if a secret should be available
	IsSecretAvailble(m *models.Secret) bool
}

// SecretService implements the methods to handle the secrets
type SecretService struct {
	db        *gorm.DB
	SecretRep rep.ISecretRepository
}

// NewSecretService returns a pointer to SecretService
func NewSecretService(db *gorm.DB) *SecretService {
	sr := rep.NewSecretRepository(db)
	return &SecretService{db: db, SecretRep: sr}
}

// CreateSecret add a new secret
func (s *SecretService) CreateSecret(secret string, remainingViews int32, expiresAt time.Time) *models.Secret {

	logrus.Debug("--> SecretService:CreateSecret")

	hash := util.CalculateSha512(secret)
	now := time.Now()
	m := &models.Secret{Hash: hash, SecretText: secret,
		CreatedAt: now, ExpiresAt: expiresAt, RemainingViews: remainingViews}

	err := s.db.Create(m)

	fmt.Println(err)

	logrus.Debug("<-- SecretService:CreateSecret")
	return m
}

// GetSecret retrieves a secret
func (s *SecretService) GetSecret(hash string) *models.Secret {
	logrus.Debug("--> SecretService:GetSecret")

	m := s.SecretRep.FindByHash(hash)

	fmt.Println(m)

	if !s.IsSecretAvailble(m) {
		logrus.Debug("<-- SecretService:GetSecret")
		return nil
	}

	m.RemainingViews = m.RemainingViews - 1
	fmt.Println(m)
	s.SecretRep.Update(m)
	logrus.Debug("<-- SecretService:GetSecret")
	return m
}

// IsSecretAvailble verifies if a secret should be available
func (s *SecretService) IsSecretAvailble(m *models.Secret) bool {

	if m == nil {
		return false
	}

	if m.RemainingViews == 0 {
		return false
	}

	now := time.Now()

	if now.After(m.ExpiresAt) {
		return false
	}

	return true
}

package repositories

import (
	"github.com/cristianchaparroa/secret/models"
	"github.com/jinzhu/gorm"
)

// ISecretRepository defines the methods to interact with secret model
type ISecretRepository interface {
	FindByHash(hash string) *models.Secret

	Update(m *models.Secret)
}

// SecretRepository implements the methods to interact with secret model
type SecretRepository struct {
	db *gorm.DB
}

// NewSecretRepository generates a pointer to  SecretRepository
func NewSecretRepository(db *gorm.DB) *SecretRepository {
	return &SecretRepository{db: db}
}

// FindByHash retrieves a secret using the hash.
func (r *SecretRepository) FindByHash(hash string) *models.Secret {
	var s models.Secret
	r.db.Table(SecretsTableName).Where("hash = ?", hash).First(&s)

	if s.Hash == "" {
		return nil
	}
	return &s
}

// Update makes an update of secret model
func (r *SecretRepository) Update(m *models.Secret) *models.Secret {
	r.db.Table(SecretsTableName).Save(m)
	return m
}

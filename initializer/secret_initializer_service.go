package initializer

import (
	"github.com/cristianchaparroa/secret/models"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// SecretInitializerService implements the functions to works the first time
type SecretInitializerService struct {
	db *gorm.DB
}

// NewSecretInitialzerService creates a pointer to SecretInitializerService
func NewSecretInitialzerService(db *gorm.DB) *SecretInitializerService {
	return &SecretInitializerService{db: db}
}

// Execute executes all the logic of this initializer
func (s *SecretInitializerService) Execute() error {

	if !s.shouldRun() {
		return nil
	}

	s.CreateTables()

	return nil
}

// ShouldInit determintes if the initializer was executed
func (s *SecretInitializerService) shouldRun() bool {

	if s.db.HasTable(&models.Secret{}) {
		return false
	}
	return true
}

// CreateTables create all the tables according with the models
func (s *SecretInitializerService) CreateTables() error {

	logrus.Debug("--> SecretInitializerService:CreateTables")

	s.db.CreateTable(&models.Secret{})

	logrus.Debug("<-- SecretInitializerService:CreateTables")
	return nil
}

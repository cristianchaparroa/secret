package initializer

import "github.com/jinzhu/gorm"

// InitService defines the operations that sould be performed the first time
// that start the application.
type InitService interface {

	//Execute
	Execute() error
}

// InitialzerManager is in charge to execute all initializer
type InitialzerManager struct {
}

// NewInitialzerManager creates a pointer to InitialzerManager
func NewInitialzerManager() *InitialzerManager {
	return &InitialzerManager{}
}

// Run  should register all initializer and execute thems
func (m *InitialzerManager) Run(db *gorm.DB) {
	initializers := make([]InitService, 0)
	initializers = append(initializers, NewSecretInitialzerService(db))

	for _, init := range initializers {
		init.Execute()
	}
}

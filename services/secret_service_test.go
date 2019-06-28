package services

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cristianchaparroa/secret/models"
	"github.com/cristianchaparroa/secret/util"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSecretServiceCreateSecret(t *testing.T) {

	db, _, _ := sqlmock.New()

	mockDB, err := gorm.Open("postgres", db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	now := time.Now()
	var test = []struct {
		secret         string
		remainingViews int32
		createdAt      time.Time
		expiresAt      int32
	}{
		{"test-secret", 1, now, 10},
		{"please save it as your life", 5, now, 10},
	}

	s := NewSecretService(mockDB)

	for _, tc := range test {

		hash := util.CalculateSha512(tc.secret)

		m := s.CreateSecret(tc.secret, tc.remainingViews, tc.expiresAt)

		assert.NotNil(t, m)
		assert.NotNil(t, m.SecretText)
		assert.NotNil(t, m.Hash)
		assert.NotNil(t, m.CreatedAt)
		assert.NotNil(t, m.ExpiresAt)
		assert.NotNil(t, m.RemainingViews)

		assert.Equal(t, tc.secret, m.SecretText, "The secret's should be equals")
		assert.Equal(t, tc.remainingViews, m.RemainingViews, "The remaining views should be equals")
		assert.Equal(t, hash, m.Hash, "The hash's should be equals")
	}
}

func TestSecretServiceGetSecret(t *testing.T) {
	db, mock, _ := sqlmock.New()

	mockDB, err := gorm.Open("postgres", db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	availableExpiresAt := time.Now().Local().Add(time.Hour*time.Duration(100) +
		time.Minute*time.Duration(100) +
		time.Second*time.Duration(100))

	var test = []struct {
		secret         string
		remainingViews int32
		ExpiresAt      time.Time
		ShouldBeNil    bool
	}{
		{"secreto", 1, availableExpiresAt, false},
		{"This is important, dont tell to anybody, please", 5, availableExpiresAt, false},
		{"You can't see this secret!", 0, availableExpiresAt, true},
	}

	s := NewSecretService(mockDB)

	for _, tc := range test {

		hash := util.CalculateSha512(tc.secret)

		rows := sqlmock.NewRows([]string{"hash", "secret_text", "remaining_views", "expires_at"}).
			AddRow(hash, tc.secret, tc.remainingViews, tc.ExpiresAt)

		mock.ExpectQuery("^SELECT (.+) FROM \"secrets\" (.+)$").WillReturnRows(rows)

		m := s.GetSecret(hash)

		if tc.ShouldBeNil && m != nil {
			t.Errorf("Expected a nil secret but get %v", m)
		}

		if !tc.ShouldBeNil {
			assert.NotNil(t, m)
			assert.Equal(t, tc.secret, m.SecretText, "The secret's should be equals")
			assert.Equal(t, hash, m.Hash, "The hash's should be equals")
		}

	}
}

func TestSecretServiceIsSecretAvailble(t *testing.T) {

	var test = []struct {
		m           *models.Secret
		IsAvailable bool
	}{
		{nil, false}, // the model is nil
		// The model should not seen cause the remainng views is zero
		{&models.Secret{RemainingViews: 0}, false},
		// The model have nil expires at
		{&models.Secret{RemainingViews: 4}, false},
		// The secret should be available
		{&models.Secret{RemainingViews: 4, ExpiresAt: time.Now().Local().Add(time.Hour*time.Duration(100) +
			time.Minute*time.Duration(100) +
			time.Second*time.Duration(100))}, true},
		// The secret should NOT  be available
		{&models.Secret{RemainingViews: 4, ExpiresAt: time.Now().Local().Add(-time.Hour - time.Duration(100) -
			time.Minute - time.Duration(100) -
			time.Second - time.Duration(100))}, false},
	}

	s := NewSecretService(nil)
	for _, tc := range test {
		isAvailable := s.IsSecretAvailble(tc.m)

		if isAvailable != tc.IsAvailable {
			t.Errorf("Expected %v, but get %v", tc.IsAvailable, isAvailable)
		}
	}
}

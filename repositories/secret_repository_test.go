package repositories

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cristianchaparroa/secret/models"
	"github.com/cristianchaparroa/secret/util"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSecretRepositoryFindByHash(t *testing.T) {
	var test = []struct {
		Secret         string
		RemainingViews int32
		ShouldBeNil    bool
	}{
		{"", 0, true},
		{"secret", 2, false},
	}

	db, mock, _ := sqlmock.New()

	mockDB, err := gorm.Open("postgres", db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	r := NewSecretRepository(mockDB)
	for _, tc := range test {

		hash := util.CalculateSha512(tc.Secret)

		var rows *sqlmock.Rows
		if !tc.ShouldBeNil {
			rows = sqlmock.NewRows([]string{"hash", "secret_text", "remaining_views"}).AddRow(hash, tc.Secret, tc.RemainingViews)
		} else {
			rows = sqlmock.NewRows([]string{"hash", "secret_text", "remaining_views"}).AddRow("", "", 0)
		}

		mock.ExpectQuery("^SELECT (.+) FROM \"secrets\" (.+)$").WillReturnRows(rows)

		m := r.FindByHash(hash)

		if !tc.ShouldBeNil {
			assert.NotNil(t, m)
			assert.Equal(t, tc.Secret, m.SecretText, "The secret's should be equals")
			assert.Equal(t, tc.RemainingViews, m.RemainingViews, "The remaining views should be equals")
			assert.Equal(t, hash, m.Hash, "The hash's should be equals")
		}

	}
}

func TestSecretRepositoryUpdate(t *testing.T) {

	var test = []struct {
		m *models.Secret
	}{
		{&models.Secret{SecretText: "je suis tombé d'amour pour toi", RemainingViews: 2}},
	}
	db, _, _ := sqlmock.New()

	mockDB, err := gorm.Open("postgres", db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	r := NewSecretRepository(mockDB)
	for _, tc := range test {
		m := r.Update(tc.m)
		assert.NotNil(t, m)

	}
}

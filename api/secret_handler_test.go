package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cristianchaparroa/secret/models"
	"github.com/cristianchaparroa/secret/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateSecret(t *testing.T) {

	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()

	mockDB, _ := gorm.Open("postgres", db)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var test = []struct {
		SecReq       *SecretPostRequest
		ExpectedCode int
	}{
		{&SecretPostRequest{Secret: "my-secret", ExpireAfterViews: 3, ExpireAfter: 20}, http.StatusOK},
		{nil, http.StatusBadRequest},
	}

	sh := NewSecretHandler(mockDB)

	for _, tc := range test {

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		if tc.SecReq != nil {
			bs, _ := json.Marshal(tc.SecReq)
			c.Request, _ = http.NewRequest("POST", "", bytes.NewBuffer(bs))
			c.Request.Header.Set("Accept", "json")
		}

		// Call the function to test
		sh.CreateSecret(c)

		if w.Code != tc.ExpectedCode {
			t.Errorf("Expected %v code but get:%v", tc.ExpectedCode, w.Code)
		}

		if w.Code == http.StatusOK {

			var sec *models.Secret
			json.NewDecoder(w.Body).Decode(&sec)
			hash := util.CalculateSha512(tc.SecReq.Secret)

			assert.NotNil(t, sec)
			assert.NotNil(t, sec.Hash)
			assert.NotNil(t, sec.CreatedAt)
			assert.NotNil(t, sec.ExpiresAt)

			assert.Equal(t, tc.SecReq.Secret, sec.SecretText, "The secret's should be equals")
			assert.Equal(t, tc.SecReq.ExpireAfterViews, sec.RemainingViews, "The remaining views should be equals")
			assert.Equal(t, hash, sec.Hash, "The hash's should be equals")
		}
	}
}

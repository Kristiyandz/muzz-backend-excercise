package apis_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kristiyandz/muzz-backend-excercise/apis"
	"github.com/Kristiyandz/muzz-backend-excercise/models/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateRandomUserHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/random-user", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apis.CreateRandomUserHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the content type is what we expect.
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	// Decode the response body to check the values.
	var gotUser user.User
	err = json.NewDecoder(rr.Body).Decode(&gotUser)
	if err != nil {
		t.Fatal(err)
	}

	expectedUser := user.User{
		Email:    "test@test.com",
		Password: "password",
		Name:     "Test",
		Gender:   "Male",
		Age:      20,
	}

	// Assert the expected user is returned.
	assert.Equal(t, expectedUser, gotUser)
}

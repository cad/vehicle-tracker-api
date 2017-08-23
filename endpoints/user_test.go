package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cad/vehicle-tracker-api/repository"
)

func TestGetAllUsersEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	user, _ := repository.CreateNewUser("test@test.com", "1234")
	_, _ = repository.CreateNewUser("test1@test.com", "1234")
	token, _ := user.RenewToken()

	// Execute
	req, _ := http.NewRequest("GET", "/user/", nil)
	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	var users []repository.User
	err := json.Unmarshal([]byte(res.Body.String()), &users)
	if err != nil {
		t.Error(errorMsg("Users", "Unmarshallable", "NotUnmarshallable"))
		return
	}

	if len(users) != 2 {
		t.Error(errorMsg("len(users)", "2", fmt.Sprintf("%d", len(users))))
		return
	}
}

func TestGetUserEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	user, _ := repository.CreateNewUser("test@test.com", "1234")
	token, _ := user.RenewToken()

	// Execute
	req, _ := http.NewRequest("GET", fmt.Sprintf("/user/%s", user.UUID), nil)
	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}

	user, err := repository.GetUserByEmail("test@test.com")
	if err != nil {
		t.Error(errorMsg("User", "ToBeFound", "NotFound"))
		return
	}

	if user.Email != "test@test.com" {
		t.Error(errorMsg("user.Email", "test@test.com", user.Email))
		return
	}

}

func TestCreateUserEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	user, _ := repository.CreateNewUser("test@test.com", "1234")
	token, _ := user.RenewToken()
	newUserEmail := "test1@test.com"
	newUserPassword := "1234"

	// Execute
	params := (struct {
		Email    string `json:"email" valid:"email"`
		Password string `json:"password"`
	}{Email: newUserEmail, Password: newUserPassword})
	paramsJSON, err := json.Marshal(&params)
	if err != nil {
		t.Error(errorMsg("UserCreatinStruct", "Marshallable", "UnMarshallable"))
	}
	body := bytes.NewBuffer(paramsJSON)
	req, _ := http.NewRequest("POST", "/user/", body)
	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 201 {
		t.Error(errorMsg("StatusCode", "201", fmt.Sprintf("%d", res.Code)))
		return
	}

	user, err = repository.GetUserByEmail(newUserEmail)
	if err != nil {
		t.Error(errorMsg("User", "ToBeFound", "NotFound"))
		return
	}

	if user.Email != newUserEmail {
		t.Error(errorMsg("user.Email", newUserEmail, user.Email))
		return
	}

}

func TestDeleteUserEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	user, _ := repository.CreateNewUser("test@test.com", "1234")
	user2, _ := repository.CreateNewUser("test1@test.com", "1234")
	token, _ := user2.RenewToken()

	// Execute
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/user/%s", user.UUID), nil)

	// Authenticate
	beforeCount := len(repository.GetAllUsers())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	afterCount := len(repository.GetAllUsers())
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}

	if beforeCount != 2 {
		t.Error(errorMsg("beforeCount", "2", fmt.Sprintf("%d", beforeCount)))
		return
	}

	if afterCount != 1 {
		t.Error(errorMsg("afterCount", "1", fmt.Sprintf("%d", afterCount)))
		return
	}
	_, err := repository.GetUserByEmail("example1@example.com")
	if err == nil {
		t.Error(errorMsg("User", "CanNotBeFound", "Found"))
	}

}

func TestDeleteUserEndpointDeleteSelf(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	user, _ := repository.CreateNewUser("test@test.com", "1234")
	token, _ := user.RenewToken()

	// Execute
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/user/%s", user.UUID), nil)

	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 403 {
		t.Error(errorMsg("StatusCode", "403", fmt.Sprintf("%d", res.Code)))
		return
	}
}

package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/cad/vehicle-tracker-api/repository"
	"github.com/gorilla/mux"
)

// swagger:parameters User
type UserParams struct {

	// Data represents user's credentials
	// in: body
	// required: true
	Data AuthorizationRequestPayload
}

// swagger:route GET /user/ Users GetAllUsers
// Get all Users.
//
//
//   Security:
//       Bearer:
//
//   Responses:
//     default: ErrorMsg
//     200: UserSuccessUsersResponse
//
func GetAllUsers(w http.ResponseWriter, req *http.Request) {
	users := repository.GetAllUsers()
	j, err := json.Marshal(users)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

// swagger:parameters GetUser
type GetUserParams struct {

	// UUID
	// in: path
	// required: true
	UUID string `json:"uuid"`
}

// swagger:route GET /user/{uuid} Users GetUser
// Get a User by UUID.
//
//   Security:
//       Bearer:
//
//   Responses:
//     default: ErrorMsg
//     200: UserSuccessUserResponse
func GetUser(w http.ResponseWriter, req *http.Request) {
	params := GetUserParams{UUID: mux.Vars(req)["uuid"]}
	//log.Println(params.UUID)
	user, _ := repository.GetUserByUUID(params.UUID)
	if user.ID == 0 {
		sendErrorMessage(w, "Not found", 404)
		return
	}
	j, err := json.Marshal(user)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

// swagger:parameters CreateNewUser
type CreateNewUserParams struct {

	// Ident represents the idetity definition of the User
	// in: body
	// required: true
	Data struct {

		// Email
		//
		// required: true
		Email string `json:"email" valid:"required"`

		// Password
		//
		// required: false
		Password string `json:"password"`
	}
}

// swagger:route POST /user/ Users CreateNewUser
// Create a new user.
//
//   Security:
//       Bearer:
//
//   Responses:
//     default: ErrorMsg
//     200: UserSuccessUserResponse
func CreateNewUser(w http.ResponseWriter, req *http.Request) {
	var params CreateNewUserParams

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&params.Data); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}
	_, err := valid.ValidateStruct(params)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := repository.CreateNewUser(
		params.Data.Email,
		params.Data.Password,
	)

	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(user)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.WriteHeader(201)
	w.Write(j)
}

// swagger:parameters DeleteUser
type DeleteUserParams struct {

	// UUID
	// in: path
	// required: true
	UUID string `json:"uuid" validate:"required"`
}

// swagger:route DELETE /user/{uuid} Users DeleteUser
// Delete a user.
//
//
//   Security:
//       Bearer:
//
//   Responses:
//     default: ErrorMsg
//     200: UserSuccessUserResponse
func DeleteUser(w http.ResponseWriter, req *http.Request) {
	params := DeleteUserParams{UUID: mux.Vars(req)["uuid"]}
	if tokenOwner, ok := UUIDFromContext(req.Context()); ok {
		if tokenOwner == params.UUID {
			sendErrorMessage(w, "User can not delete themself.", http.StatusForbidden)
			return
		}
	} else {
		sendErrorMessage(w, "tokenOwner can not be determined.", http.StatusInternalServerError)
		return
	}

	user, err := repository.DeleteUserByUUID(params.UUID)
	if err != nil {
		sendErrorMessage(w, "User to be deleted, can not be found", http.StatusNotFound)
		return
	}

	j, err := json.Marshal(user)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

type userKey int

const userUUIDKey userKey = 0

// NewUUIDContext creates a new ctx with the given UUID.
func NewUUIDContext(ctx context.Context, uuid string) context.Context {
	return context.WithValue(ctx, userUUIDKey, uuid)
}

// UUIDFromContext extracts the user UUID from ctx, if present.
func UUIDFromContext(ctx context.Context) (string, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the string type assertion returns ok=false for nil.
	uuid, ok := ctx.Value(userUUIDKey).(string)
	return uuid, ok
}

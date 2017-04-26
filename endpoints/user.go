package endpoints


import (
	"github.com/gorilla/mux"
	valid "github.com/asaskevich/govalidator"
	"net/http"
	"encoding/json"
	"github.com/cad/vehicle-tracker-api/repository"
//	"log"
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

	// Email
	// in: path
	// required: true
	Email string `json:"email"`

}

// swagger:route GET /user/{email} Users GetUser
// Get a User by email address.
//
//   Security:
//       Bearer:
//
//   Responses:
//     default: ErrorMsg
//     200: UserSuccessUserResponse
func GetUser(w http.ResponseWriter, req *http.Request) {
	params := GetUserParams{Email: mux.Vars(req)["email"]}
	user, _ := repository.GetUserByEmail(params.Email)
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

	if err != nil  {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(user)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

// swagger:parameters DeleteUser
type DeleteUserParams struct {

	// Email
	// in: path
	// required: true
	Email string `json:"email" validate:"required"`

}

// swagger:route DELETE /user/{email} Users DeleteUser
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
	params := DeleteUserParams{Email: mux.Vars(req)["email"]}

	user, err := repository.DeleteUserByEmail(params.Email)
	if err != nil {
		sendErrorMessage(w, "User to be deleted, can not be found", http.StatusNotFound)
		return
	}

	j, err := json.Marshal(user)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

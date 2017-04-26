package endpoints


import (
	"net/http"
	"strings"
	"encoding/json"
	"github.com/cad/vehicle-tracker-api/repository"
	"log"
)

func TokenAuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := strings.Split(r.Header.Get("Authorization"), " ")

		if len(s) != 2 {
			sendErrorMessage(w, "Authorization token should be in the form of Authorization: Bearer <token>", 401)
			return
		}

		if s[0] != "Bearer" {
			sendErrorMessage(w, "Authorization token should be in the form of Authorization: Bearer <token>", 401)
			return
		}

		// Check token
		token := s[1]

		if repository.CheckToken(token) != true {
			// Rejected
			sendErrorMessage(w, "Not Authorized", 401)
			return
		}
		// Permitted

		// TODO(cad): Put the tokenOwner in a Context and pass it around
		// with the http.Request object in yor API.


		h.ServeHTTP(w, r)
	}
}


type AuthorizationRequestPayload struct {
	Email string `json:"email" valid:"email"`
	Password string `json:"password"`
}

// swagger:parameters Authorize
type AuthorizationParams struct {

	// Data represents user's credentials
	// in: body
	// required: true
	Data AuthorizationRequestPayload

}


// swagger:route POST /auth/ Auth Authorize
// Get an `authorization_token`.
//
//
//   Responses:
//     default: ErrorMsg
//     200: AuthSuccessTokenResponse
func Authorize(w http.ResponseWriter, req *http.Request) {
	var params AuthorizationParams
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params.Data); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}


	user, err := repository.GetUserByEmail(params.Data.Email)
	if err != nil {
		log.Println(err.Error())
		sendErrorMessage(w, "Invalid Credentials U", 401)
		return
	}

	if user.CheckPassword(params.Data.Password) != true {
		sendErrorMessage(w, "Invalid Credentials P", 401)
		return
	}

	token, err := user.RenewToken()
	if err != nil {
		log.Println(err.Error())
		sendErrorMessage(w, "Can not create token", 500)
		return
	}
	payload := AuthorizationResponsePayload{
		AuthorizationToken: token,
	}
	j, err := json.Marshal(payload)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}


// swagger:route GET /auth/ Auth CheckAuth
// See if you are authenticated or not.
//
//   Security:
//       Bearer:
//
//   Responses:
//     default: ErrorMsg
//     200: AuthSuccessOKResponse
func CheckAuth(w http.ResponseWriter, req *http.Request) {
	payload := AuthorizationCheckResponsePayload{Authorized: true}
	j, err := json.Marshal(payload)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

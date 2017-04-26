package endpoints

import (
	"github.com/cad/vehicle-tracker-api/repository"
)


// Returns a user
// swagger:response
type UserSuccessUserResponse struct {
	// User
	// in: body
	Body repository.User
}


// Returns list of users
// swagger:response
type UserSuccessUsersResponse struct {
	// Users
	// in: body
	Body []repository.User

}

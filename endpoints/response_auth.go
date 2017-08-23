package endpoints

import "github.com/cad/vehicle-tracker-api/repository"

//"github.com/cad/vehicle-tracker-api/repository"

type AuthorizationResponsePayload struct {
	AuthorizationToken string `json:"authorization_token"`
}

type AuthorizationCheckResponsePayload struct {
	Authorized bool            `json:"authorized"`
	User       repository.User `json:"user"`
}

// Returns a `authorization_token`
// swagger:response
type AuthSuccessTokenResponse struct {
	// AuthorizationToken
	// in: body
	Body AuthorizationResponsePayload
}

// Returns ok if authenticated
// swagger:response
type AuthSuccessOKResponse struct {
	// AuthorizationCheck
	// in: body
	Body AuthorizationCheckResponsePayload
}
